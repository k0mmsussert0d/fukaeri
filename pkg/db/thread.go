package db

import (
	"context"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Threads(client *mongo.Client) *mongo.Collection {
	return DB(client).Collection("threads")
}

func GetThread(client *mongo.Client, no int) models.Thread {
	var res models.Thread
	err := Threads(client).FindOne(
		context.TODO(),
		bson.D{{"no", no}},
	).Decode(&res)

	internal.HandleError(err)

	return res
}

func SaveThread(client *mongo.Client, thread models.Thread) {
	Threads(client).InsertOne(
		context.TODO(),
		thread,
	)
}
