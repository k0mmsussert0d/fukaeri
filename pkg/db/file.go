package db

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/internal/conf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FileExists(md5 string) bool {
	mongoDB := DB(MongoClient())
	fs := conf.GetConfig().DB.Files
	if fs == nil {
		*fs = "fs"
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	return mongoDB.Collection(fmt.Sprintf("%s.%s", *fs, "files")).FindOne(
		ctx,
		bson.D{{"filename", md5}},
	).Err() != mongo.ErrNoDocuments
}

func SaveFile(file []byte, filetype string) {
	mongoDB := DB(MongoClient())
	fs := conf.GetConfig().DB.Files
	if fs == nil {
		*fs = "fs"
	}

	bucket, err := gridfs.NewBucket(
		mongoDB,
		options.GridFSBucket().SetName(*fs),
	)
	internal.HandleError(err)

	fileMd5 := fmt.Sprintf("%x", md5.Sum(file))

	uploadStream, err := bucket.OpenUploadStream(
		fileMd5,
		options.GridFSUpload(),
	)
	internal.HandleError(err)

	err = uploadStream.SetWriteDeadline(time.Now().Add(5 * time.Second))
	internal.HandleError(err)

	if _, err = uploadStream.Write(file); err != nil {
		internal.HandleError(err)
	}
}
