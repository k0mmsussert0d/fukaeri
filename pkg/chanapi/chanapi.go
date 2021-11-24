package chanapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/k0mmsussert0d/fukaeri/pkg/models"
)

type Request struct {
	method, endpoint string
	responseChan     chan []byte
}

var api string = "https://a.4cdn.org"
var httpClientInstance *http.Client = nil
var requestsChan (chan *Request)
var exitChan (chan bool)

func httpClient() *http.Client {
	if httpClientInstance == nil {
		httpClientInstance = &http.Client{}
	}

	return httpClientInstance
}

func fetch(method, endpoint string) []byte {
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

	return bodyBytes
}

func initChans() {
	requestsChan = make(chan *Request)
	exitChan = make(chan bool)
}

func initClient() {
	var limiter = time.NewTicker(1 * time.Second)
	for {
		select {
		case req := <-requestsChan:
			<-limiter.C
			req.responseChan <- fetch(req.method, req.endpoint)
		case <-exitChan:
			return
		}
	}
}

func StartClient() {
	initChans()
	go initClient()
}

func StopClient() {
	exitChan <- true
}

func Threads(board string) models.Threads {
	request := &Request{"GET", fmt.Sprintf("%s/%s/threads.json", api, board), make(chan []byte)}

	requestsChan <- request
	response := <-request.responseChan

	var returnObj models.Threads
	json.Unmarshal(response, &returnObj)
	return returnObj
}

func Thread(board string, id string) models.Thread {
	request := &Request{"GET", fmt.Sprintf("%s/%s/thread/%s.json", api, board, id), make(chan []byte)}

	requestsChan <- request
	response := <-request.responseChan

	var returnObj models.Thread
	json.Unmarshal(response, &returnObj)
	return returnObj
}
