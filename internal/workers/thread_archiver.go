package workers

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/internal/log"
	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/apiclient"
	"github.com/k0mmsussert0d/fukaeri/pkg/db"
	"github.com/k0mmsussert0d/fukaeri/pkg/models"
)

func ArchiveThread(board string, id int, chanapi apiclient.ApiClient, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	log.Info().Printf("Started archiver for thread %v/%v", board, id)

	var lastUpdateTime time.Time

	refreshRates := []int{10, 20, 30, 45, 60, 90, 180}

	refreshRateIdx := 0
	initRefreshRate := refreshRates[0]

	refreshThreadTicker := time.NewTicker(time.Duration(initRefreshRate) * time.Second)

	doWork := func() {
		defer func() {
			r := recover()
			if r != nil {
				log.Error().Printf("An error occured while archiving thread %v/%v: %v", board, id, r)
			}
		}()

		writeToDB := func(thread models.Thread) {
			db.SaveOrUpdateThread(board, strconv.Itoa(id), thread, ctx)
			db.SaveFilesFromThread(board, thread, chanapi, ctx)
		}

		refreshThreadTicker.Stop()
		if lastUpdateTime.IsZero() {
			log.Info().Printf("Fetching full %v/%v thread for the first time", board, id)
			thread, err := chanapi.Thread(ctx, board, strconv.Itoa(id))
			internal.HandleError(err)
			refreshThreadTicker.Reset(time.Duration(initRefreshRate) * time.Second)
			writeToDB(*thread)
		} else {
			log.Info().Printf("Refreshing thread %v/%v for new posts since %v", board, id, lastUpdateTime)
			thread, err := chanapi.ThreadSince(ctx, board, strconv.Itoa(id), lastUpdateTime)
			internal.HandleError(err)
			if thread != nil {
				log.Info().Printf("Updating thread %v/%v as new posts appeared", board, id)
				log.Debug().Printf("Next refresh for %v/%v in %v seconds", board, id, initRefreshRate)
				refreshRateIdx = 0
				refreshThreadTicker.Reset(time.Duration(initRefreshRate) * time.Second)
				writeToDB(*thread)
			} else {
				log.Info().Printf("No new posts for thread %v/%v", board, id)
				if refreshRateIdx < len(refreshRates)-1 {
					refreshRateIdx += 1
				}
				log.Debug().Printf("Next refresh for %v/%v in %v seconds", board, id, refreshRates[refreshRateIdx])
				refreshThreadTicker.Reset(time.Duration(refreshRates[refreshRateIdx]) * time.Second)
			}
		}
		lastUpdateTime = time.Now()
	}

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
