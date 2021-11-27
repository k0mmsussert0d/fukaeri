package db

import (
	"context"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/internal/conf"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoClient() *mongo.Client {
	config := conf.GetConfig()
	uri := config.DB.Connstring
	mongoClientInstance, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	internal.HandleError(err)
	return mongoClientInstance
}

func DB(client *mongo.Client) *mongo.Database {
	return client.Database(conf.GetConfig().DB.Name)
}
