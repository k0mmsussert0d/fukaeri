package workers

import (
	"context"
	"log"
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

	log.Printf("Started archiver for thread %v/%v", board, id)

	lastUpdateTime := time.Now()

	doWork := func() {
		log.Printf("Refreshing thread %v/%v for new posts since %v", board, id, lastUpdateTime)
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

	refreshThreadTicker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-refreshThreadTicker.C:
			doWork()
		case <-ctx.Done():
			log.Printf("Thread %v/%v archiver received exit signal. Shutting down.", board, id)
			return
		}
	}
}
