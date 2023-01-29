package db

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/internal/conf"
	"github.com/k0mmsussert0d/fukaeri/internal/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileDetails struct {
	Posts  []int64
	Width  int
	Height int
	Size   int
}

func FileExists(md5 []byte) bool {
	mongoDB := DB(MongoClient())
	fs := conf.Get().DB.Files

	log.Logger().Debugw("Checkinf if file already exists in the bucket",
		"md5", md5,
		"bucket", fs,
	)

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	res := mongoDB.Collection(fmt.Sprintf("%s.%s", fs, "files")).FindOne(
		ctx,
		bson.D{{"metadata", bson.D{{"md5", md5}}}},
	)
	if res.Err() != mongo.ErrNoDocuments {
		log.Logger().Debugw("File already exists",
			"md5", md5,
		)
		return true
	} else {
		log.Logger().Debugw("File has not been archived yet",
			"md5", md5,
		)
		return false
	}

}

func SaveFile(file []byte, md5 []byte, details FileDetails) {
	mongoDB := DB(MongoClient())
	fs := conf.Get().DB.Files

	log.Logger().Debugw("Saving file",
		"md5", md5,
		"bucket", fs,
	)

	bucket, err := gridfs.NewBucket(
		mongoDB,
		options.GridFSBucket().SetName(fs),
	)
	if err != nil {
		log.Logger().Errorw("Failed to initialize GridFS bucket",
			"bucket", fs,
		)
		internal.HandleError(err)
	}

	md5sum := hex.EncodeToString(md5)

	uploadStream, err := bucket.OpenUploadStream(
		md5sum,
		options.GridFSUpload().SetMetadata(bson.D{{"md5", md5}}),
	)
	if err != nil {
		log.Logger().Errorw("Failed to open UploadStream",
			"file", md5sum,
		)
		internal.HandleError(err)
	}

	err = uploadStream.SetWriteDeadline(time.Now().Add(5 * time.Second))
	internal.HandleError(err)

	if _, err = uploadStream.Write(file); err != nil {
		log.Logger().Errorw("Failed to save file",
			"md5", md5sum,
		)
		internal.HandleError(err)
	}

	err = uploadStream.Close()
	if err != nil {
		log.Logger().Errorw("Failed to commit file metadata",
			"md5", md5sum,
		)
	}
	internal.HandleError(err)

	log.Logger().Debugw("File saved successfully",
		"md5", md5,
	)
}

func AddPostToFile(md5 []byte, postNo int64) {
	mongoDB := DB(MongoClient())
	fs := conf.Get().DB.Files

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	mongoDB.Collection(fmt.Sprintf("%s.%s", fs, "files")).UpdateOne(
		ctx,
		bson.D{{"metadata", bson.D{{"md5", md5}}}},
		bson.D{{"$push", bson.D{{"posts", postNo}}}},
	)
}
