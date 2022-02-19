package workers

import (
	"context"
	"sync"

	"github.com/k0mmsussert0d/fukaeri/internal/conf"
	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/apiclient"
)

func StartArchiving(ctx context.Context) {
	boards := conf.GetConfig().Archive.Boards
	chanapi := apiclient.New()

	boardsWg := &sync.WaitGroup{}

	for _, board := range boards {
		boardsWg.Add(1)
		go ArchiveBoard(board, *chanapi, boardsWg, ctx)
	}

	boardsWg.Wait()
}
