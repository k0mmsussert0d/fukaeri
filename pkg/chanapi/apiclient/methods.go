package apiclient

import (
	"encoding/json"
	"fmt"

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
