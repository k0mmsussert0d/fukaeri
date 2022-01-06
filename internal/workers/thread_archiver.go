package workers

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/k0mmsussert0d/fukaeri/internal/log"
	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/apiclient"
	"github.com/k0mmsussert0d/fukaeri/pkg/db"
)

func ArchiveThread(board string, id int, chanapi apiclient.ApiClient, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	log.Info().Printf("Started archiver for thread %v/%v", board, id)

	var lastUpdateTime time.Time

	doWork := func() {
		if lastUpdateTime.IsZero() {
			log.Info().Printf("Fetching full %v/%v thread for the first time", board, id)
			thread := chanapi.Thread(board, strconv.Itoa(id))
			db.SaveOrUpdateThread(board, strconv.Itoa(id), thread, ctx)
		} else {
			log.Info().Printf("Refreshing thread %v/%v for new posts since %v", board, id, lastUpdateTime)
			thread, updated := chanapi.ThreadSince(board, strconv.Itoa(id), lastUpdateTime)
			if updated {
				log.Info().Printf("Updating thread %v/%v as new posts appeared", board, id)
				db.SaveOrUpdateThread(board, strconv.Itoa(id), thread, ctx)
			} else {
				log.Info().Printf("No new posts for thread %v/%v", board, id)
			}
		}
		lastUpdateTime = time.Now()
	}

	refreshThreadTicker := time.NewTicker(30 * time.Second)

	for {
		doWork()

		select {
		case <-refreshThreadTicker.C:
			continue
		case <-ctx.Done():
			log.Info().Printf("Thread %v/%v archiver received exit signal. Shutting down.", board, id)
			return
		}
	}
}
