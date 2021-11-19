package chanapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/k0mmsussert0d/fukaeri/pkg/models"
)

var _api string = "https://a.4cdn.org"
var _httpClientInstance *http.Client = nil

func httpClient() *http.Client {
	if _httpClientInstance == nil {
		_httpClientInstance = &http.Client{}
	}

	return _httpClientInstance
}

func fetch(method, endpoint string, responseObject interface{}) {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient().Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	json.Unmarshal(bodyBytes, &responseObject)
}

func Threads(board string) models.Threads {
	var responseObject models.Threads
	fetch("GET", fmt.Sprintf("%s/%s/threads.json", _api, board), &responseObject)
	return responseObject
}

func Thread(board string, id string) models.Thread {
	var responseObject models.Thread
	fetch("GET", fmt.Sprintf("%s/%s/thread/%s.json", _api, board, id), &responseObject)
	return responseObject
}
