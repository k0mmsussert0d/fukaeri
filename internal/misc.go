package internal

import (
	"errors"

	"github.com/k0mmsussert0d/fukaeri/internal/log"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleError(err error) {
	if err != nil {
		log.Error().Println(err)
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

func ConvertPanicToError(panicVal interface{}) error {
	switch x := panicVal.(type) {
	case string:
		return errors.New(x)
	case error:
		return x
	default:
		return errors.New("unknown panic type")
	}
}
