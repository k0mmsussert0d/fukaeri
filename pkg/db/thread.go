package db

import (
	"context"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveOrUpdateThread(board string, id string, thread models.Thread, ctx context.Context) {
	doc, err := internal.ToBSONDoc(thread)
	internal.HandleError(err)

	mongoDB := DB(MongoClient())

	if _, err := mongoDB.Collection(board).ReplaceOne(
		ctx,
		bson.D{{"_id", id}},
		doc,
		options.Replace().SetUpsert(true),
	); err != nil {
		internal.HandleError(err)
	}
}
