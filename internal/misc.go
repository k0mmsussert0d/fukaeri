package internal

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func ToBSONDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
