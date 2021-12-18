package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/pkg/models"
)

func (client ApiClient) Threads(board string) models.Threads {
	res := client.fetch("GET", fmt.Sprintf("%s/%s/threads.json", client.endpoint, board))
	var returnObj models.Threads
	json.Unmarshal(res, &returnObj)
	return returnObj
}

func (client ApiClient) Thread(board string, id string) models.Thread {
	res := client.fetch("GET", fmt.Sprintf("%s/%s/thread/%s.json", client.endpoint, board, id))
	var returnObj models.Thread
	json.Unmarshal(res, &returnObj)
	return returnObj
}

func (client ApiClient) ThreadSince(board, id string, since time.Time) models.Thread {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/thread/%s.json", client.endpoint, board, id), nil)
	internal.HandleError(err)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("If-Modified-Since", since.Local().Format(http.TimeFormat))

	var returnObj models.Thread
	res := client.fetchRequest(req)
	json.Unmarshal(res, &returnObj)
	return returnObj
}
