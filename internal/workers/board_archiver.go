package workers

import (
	"context"
	"sync"
	"time"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/internal/log"
	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/apiclient"
)

func ArchiveBoard(board string, chanapi apiclient.ApiClient, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	log.Logger().Infof("Started archiver for %v board", board)
	// threads currently processed by archiver and their workers cancelation functions
	currentThreads := make(map[int]func())
	threadsWg := &sync.WaitGroup{}

	doWork := func() {
		defer func() {
			r := recover()
			if r != nil {
				log.Logger().Errorw("An error occurred while archiving board",
					"board", board,
					"error", r,
				)
			}
		}()

		// threads currently active on a board
		activeThreads := make(map[int]interface{})

		// fetch all active threads from a board
		threads, err := chanapi.Threads(ctx, board)
		internal.HandleError(err)
		for _, threadsOnPage := range *threads {
			for _, threadInfo := range threadsOnPage.Threads {
				activeThreads[threadInfo.No] = nil
			}
		}

		// cancel all archived/deleted threads archiving
		for currentThreadId, cancel := range currentThreads {
			_, stillActive := activeThreads[currentThreadId]
			if !stillActive {
				log.Logger().Infow("Thread has been archived or deleted, aborting the task",
					"board", board,
					"thread", currentThreadId,
				)
				cancel()
				delete(currentThreads, currentThreadId)
			}
		}

		// start archiving new threads
		for activeThreadId := range activeThreads {
			_, alreadyKnown := currentThreads[activeThreadId]
			if !alreadyKnown {
				log.Logger().Infow("New thread appreared, starting new task",
					"board", board,
					"thread", activeThreadId,
				)
				threadCtx, cancelThread := context.WithCancel(ctx)
				threadsWg.Add(1)
				go ArchiveThread(board, activeThreadId, chanapi, threadsWg, threadCtx)
				currentThreads[activeThreadId] = cancelThread
			}
		}

	}

	refreshBoardTicker := time.NewTicker(60 * time.Second)
	for {
		doWork()

		select {
		case <-refreshBoardTicker.C:
			continue
		case <-ctx.Done():
			log.Logger().Infof("Board %v archiver received exit signal. Waiting for thread archivers to shutdown.", board)
			threadsWg.Wait()
			log.Logger().Infof("All thread archivers for board %v terminated. Shutting down board archiver.", board)
			return
		}
	}
}
