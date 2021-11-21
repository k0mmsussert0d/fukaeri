package db

import (
	"context"
	"log"

	"github.com/k0mmsussert0d/fukaeri/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func coll() *mongo.Collection {
	return database().Collection("threads")
}

func GetThread(no int) models.Thread {
	var res models.Thread
	err := coll().FindOne(
		context.TODO(),
		bson.D{{"no", no}},
	).Decode(&res)

	if err != nil {
		log.Fatal("thread not found")
	}

	return res
}

func SaveThread(thread models.Thread) {
	coll().InsertOne(
		context.TODO(),
		thread,
	)
}
