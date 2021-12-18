package workers

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/apiclient"
)

func ArchiveBoard(board string, chanapi apiclient.ApiClient, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	log.Printf("Started archiver for %v board", board)
	// threads currently processed by archiver and their workers cancelation functions
	currentThreads := make(map[int]func())
	threadsWg := &sync.WaitGroup{}

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
				log.Printf("Thread %v/%v has been archived or deleted – aborting archivization", board, currentThreadId)
				cancel()
				delete(currentThreads, currentThreadId)
			}
		}

		// start archiving new threads
		for activeThreadId := range activeThreads {
			_, alreadyKnown := currentThreads[activeThreadId]
			if !alreadyKnown {
				log.Printf("New thread %v/%v appeared. Spinning up new archiver", board, activeThreadId)
				threadCtx, cancelThread := context.WithCancel(ctx)
				threadsWg.Add(1)
				go ArchiveThread(board, activeThreadId, chanapi, threadsWg, threadCtx)
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
			threadsWg.Wait()
			return
		}
	}
}
