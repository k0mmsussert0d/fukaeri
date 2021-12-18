package db

import (
	"context"
	"log"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/internal/conf"
	"github.com/k0mmsussert0d/fukaeri/pkg/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitCollections(ctx context.Context) {
	boards := conf.GetConfig().Archive.Boards
	log.Printf("Initializing database collections for boards: %v", boards)

	for _, board := range boards {
		log.Printf("Initializing %v collection", board)
		CreateCollectionIfNotExists(board, ctx)
	}
}

func CreateCollectionIfNotExists(board string, ctx context.Context) {
	err := db.DB(db.MongoClient()).CreateCollection(
		ctx,
		board,
		options.CreateCollection(),
	)

	if err != nil {
		_, ok := err.(mongo.CommandError)
		if !ok {
			internal.HandleError(err)
		} else {
			log.Printf("Collection for %v board already exists", board)
		}
	} else {
		log.Printf("Created collection for %v board", board)
	}
}
