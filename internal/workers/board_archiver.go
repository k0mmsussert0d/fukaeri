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

	log.Info().Printf("Started archiver for %v board", board)
	// threads currently processed by archiver and their workers cancelation functions
	currentThreads := make(map[int]func())
	threadsWg := &sync.WaitGroup{}

	doWork := func() {
		defer func() {
			r := recover()
			if r != nil {
				log.Error().Printf("An error occured while archiving board %v: %v", board, r)
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
				log.Info().Printf("Thread %v/%v has been archived or deleted – aborting archivization...", board, currentThreadId)
				cancel()
				delete(currentThreads, currentThreadId)
			}
		}

		// start archiving new threads
		for activeThreadId := range activeThreads {
			_, alreadyKnown := currentThreads[activeThreadId]
			if !alreadyKnown {
				log.Info().Printf("New thread %v/%v appeared. Spinning up new archiver...", board, activeThreadId)
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
			log.Info().Printf("Board %v archiver received exit signal. Waiting for thread archivers to exit...", board)
			threadsWg.Wait()
			log.Info().Printf("All thread archivers for board %v terminated. Shutting down board archiver", board)
			return
		}
	}
}
