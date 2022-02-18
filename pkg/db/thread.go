package db

import (
	"context"
	"encoding/base64"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/apiclient"
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

func SaveFilesFromThread(board string, thread models.Thread, chanapi apiclient.ApiClient, ctx context.Context) {
	savePost := func(thread models.Thread, post_idx int) {
		post := thread.Posts[post_idx]
		if post.Filename != "" {
			md5, err := base64.StdEncoding.DecodeString(post.Md5)
			internal.HandleError(err)

			if !FileExists(md5) {
				details := FileDetails{
					[]int64{post.Tim},
					post.W,
					post.H,
					post.Fsize,
				}
				SaveFile(chanapi.File(board, post.Tim, post.Ext), md5, details)
			} else {
				AddPostToFile(md5, post.Tim)
			}
		}
	}

	for idx := range thread.Posts {
		select {
		case <-ctx.Done():
			return
		default:
			savePost(thread, idx)
		}
	}
}
