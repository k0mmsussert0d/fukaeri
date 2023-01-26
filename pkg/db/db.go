package db

import (
	"context"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/internal/conf"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoClient() *mongo.Client {
	uri := conf.Get().DB.Connstring
	mongoClientInstance, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	internal.HandleError(err)
	return mongoClientInstance
}

func DB(client *mongo.Client) *mongo.Database {
	return client.Database(conf.Get().DB.Collection)
}

func Bucket(client *mongo.Client) *gridfs.Bucket {
	bucket_name := conf.Get().DB.Files
	bucket, err := gridfs.NewBucket(
		DB(client),
		&options.BucketOptions{
			Name: &bucket_name,
		},
	)
	internal.HandleError(err)
	return bucket
}
