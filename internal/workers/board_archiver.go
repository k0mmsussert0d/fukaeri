package workers

import (
	"context"
	"time"

	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/apiclient"
)

func ArchiveBoard(board string, chanapi apiclient.ApiClient, ctx context.Context) {
	// threads currently processed by archiver and their workers cancelation functions
	currentThreads := make(map[int]func())

	doWork := func() {
		// threads currently active on a board
		activeThreads := make(map[int]interface{})

		// fetch all active threads from a board
		for _, threadsOnPage := range chanapi.Threads(board) {
			for _, threadInfo := range threadsOnPage.Threads {
				activeThreads[threadInfo.No] = nil
			}
		}

		// cancel all archived/deleted threads archiving
		for currentThreadId, cancel := range currentThreads {
			_, stillActive := activeThreads[currentThreadId]
			if !stillActive {
				cancel()
				delete(currentThreads, currentThreadId)
			}
		}

		// start archiving new threads
		for activeThreadId := range activeThreads {
			_, alreadyKnown := currentThreads[activeThreadId]
			if !alreadyKnown {
				threadCtx, cancelThread := context.WithCancel(ctx)
				go ArchiveThread(board, activeThreadId, chanapi, threadCtx)
				currentThreads[activeThreadId] = cancelThread
			}
		}

	}

	refreshBoardTicker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-refreshBoardTicker.C:
			doWork()
		case <-ctx.Done():
			return
		}
	}
}
