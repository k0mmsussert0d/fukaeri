package db

import (
	"context"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/internal/conf"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _mongoClientInstance *mongo.Client = nil

func mongoClient() *mongo.Client {
	if _mongoClientInstance == nil {
		config := conf.GetConfig()
		uri := config.DB.Connstring
		var err error
		_mongoClientInstance, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		internal.HandleError(err)
	}
	return _mongoClientInstance
}

func database() *mongo.Database {
	return mongoClient().Database(conf.GetConfig().DB.Name)
}
