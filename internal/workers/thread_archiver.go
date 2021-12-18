package workers

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/apiclient"
	"github.com/k0mmsussert0d/fukaeri/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ArchiveThread(board string, id int, chanapi apiclient.ApiClient, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	lastUpdateTime := time.Now()

	doWork := func() {
		thread := chanapi.ThreadSince(board, strconv.Itoa(id), lastUpdateTime)

		mongo := db.DB(db.MongoClient())

		doc, err := internal.ToBSONDoc(thread)
		internal.HandleError(err)

		// replace existing thread document with updated version or create a new one
		if _, err := mongo.Collection(board).ReplaceOne(
			ctx,
			bson.D{{"_id", id}},
			doc,
			options.Replace().SetUpsert(true),
		); err != nil {
			internal.HandleError(err)
		}
	}

	refreshThreadTicker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-refreshThreadTicker.C:
			doWork()
		case <-ctx.Done():
			return
		}
	}
}
