package internal

import "log"

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
