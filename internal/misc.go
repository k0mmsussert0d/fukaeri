package internal

import (
	"go.mongodb.org/mongo-driver/bson"
)

func ToBSONDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
