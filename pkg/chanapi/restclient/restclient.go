package restclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/k0mmsussert0d/fukaeri/pkg/models"
)

type RestClient struct {
	httpClient *http.Client
	endpoint   string
}

func New() *RestClient {
	client := &RestClient{}
	client.Init()
	return client
}

func (client *RestClient) Init() {
	client.httpClient = &http.Client{}
	client.endpoint = "https://a.4cdn.org"
}

func (client RestClient) fetch(method, endpoint string) []byte {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	return bodyBytes
}

func (client RestClient) Threads(board string) models.Threads {
	res := client.fetch("GET", fmt.Sprintf("%s/%s/threads.json", client.endpoint, board))
	var returnObj models.Threads
	json.Unmarshal(res, &returnObj)
	return returnObj
}

func (client RestClient) Thread(board string, id string) models.Thread {
	res := client.fetch("GET", fmt.Sprintf("%s/%s/thread/%s.json", client.endpoint, board, id))
	var returnObj models.Thread
	json.Unmarshal(res, &returnObj)
	return returnObj
}
