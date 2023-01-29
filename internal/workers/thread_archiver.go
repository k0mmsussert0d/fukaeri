package workers

import (
	"context"
	"errors"
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

	log.Logger().Infow("Started new thread archived",
		"board", board,
		"id", id,
	)

	var lastUpdateTime time.Time

	refreshRates := []int{10, 20, 30, 45, 60, 90, 180}

	refreshRateIdx := 0
	initRefreshRate := refreshRates[0]

	refreshThreadTicker := time.NewTicker(time.Duration(initRefreshRate) * time.Second)

	doWork := func() {
		defer func() {
			r := recover()
			if r != nil {
				switch x := r.(type) {
				case internal.FukaeriError:
				case error:
					if !errors.Is(x, context.Canceled) { // context cancelation errors are expected
						internal.HandleError(internal.WrapError(x, "An error occured while archiving thread %v/%v", board, id))
					}
				default:
					err := internal.WrapError(errors.New("Unknown error occured while archiving thread"), "Unknown error occured while archiving thread %v/%v", board, id)
					err.Misc["org_value"] = x
					internal.HandleError(err)
				}
			}
		}()

		writeToDB := func(thread models.Thread) {
			if ctx.Err() != nil {
				return
			}
			db.SaveOrUpdateThread(board, strconv.Itoa(id), thread, ctx)
			db.SaveFilesFromThread(board, thread, chanapi, ctx)
		}

		refreshThreadTicker.Stop()
		if lastUpdateTime.IsZero() {
			log.Logger().Infow("Fetching full thread for the first time",
				"board", board,
				"id", id,
			)
			thread, err := chanapi.Thread(ctx, board, strconv.Itoa(id))
			internal.HandleError(err)
			refreshThreadTicker.Reset(time.Duration(initRefreshRate) * time.Second)
			writeToDB(*thread)
		} else {
			log.Logger().Infow("Refreshing thread",
				"board", board,
				"id", id,
				"modifiedSince", lastUpdateTime,
			)
			thread, err := chanapi.ThreadSince(ctx, board, strconv.Itoa(id), lastUpdateTime)
			internal.HandleError(err)
			if thread != nil {
				log.Logger().Infow("Updating thread as new posts appeared",
					"board", board,
					"id", id,
				)
				writeToDB(*thread)
				refreshRateIdx = 0
				refreshThreadTicker.Reset(time.Duration(initRefreshRate) * time.Second)
				log.Logger().Debugw("Next refresh scheduled",
					"board", board,
					"id", id,
					"nextRefreshIn", initRefreshRate,
				)
			} else {
				log.Logger().Infow("No new posts appeared",
					"board", board,
					"id", id,
				)
				if refreshRateIdx < len(refreshRates)-1 {
					refreshRateIdx += 1
				}
				log.Logger().Debugw("Next refresh scheduled",
					"board", board,
					"id", id,
					"nextRefreshIn", refreshRates[refreshRateIdx],
				)
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
			log.Logger().Infow("Thread archived received exit signal, shutting down",
				"board", board,
				"id", id,
			)
			return
		}
	}
}
