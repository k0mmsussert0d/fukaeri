package apiclient

import (
	"errors"
	"net/http"
	"time"
)

type ApiClient struct {
	httpClient *http.Client
	limiter    *time.Ticker
	queue      chan http.Request
}

type httpClientRequest struct {
	request  *http.Request
	response <-chan http.Response
}

func (client *ApiClient) Init() {
	client.httpClient = &http.Client{}
	client.queue = make(chan http.Request, 1)
	client.limiter = time.NewTicker(1 * time.Second)
}

func (client *ApiClient) Do(req *http.Request) (*http.Response, error) {
	request := httpClientRequest{req, make(chan http.Response)}

	for response := range request.response {
		return &response, nil
	}

	return nil, errors.New("whatever")
}
