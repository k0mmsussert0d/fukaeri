package db

import (
	"context"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/internal/conf"
	"github.com/k0mmsussert0d/fukaeri/internal/log"
	"github.com/k0mmsussert0d/fukaeri/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitCollections(ctx context.Context) {
	boards := conf.GetConfig().Archive.Boards
	log.Info().Printf("Initializing database collections for boards: %v", boards)

	for _, board := range boards {
		log.Info().Printf("Initializing %v collection", board)
		CreateCollectionIfNotExists(board, ctx)
	}
}

func CreateCollectionIfNotExists(board string, ctx context.Context) {
	exists, _ := db.DB(db.MongoClient()).ListCollectionNames(
		ctx,
		bson.D{{"name", board}},
	)

	if len(exists) == 0 {
		err := db.DB(db.MongoClient()).CreateCollection(
			ctx,
			board,
			options.CreateCollection(),
		)

		internal.HandleError(err)
	}
}
