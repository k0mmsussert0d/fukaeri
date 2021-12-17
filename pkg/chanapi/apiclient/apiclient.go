package apiclient

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/limitedhttpclient"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ApiClient struct {
	httpClient HttpClient
	endpoint   string
}

func New(ctx context.Context) *ApiClient {
	return &ApiClient{
		httpClient: limitedhttpclient.New(ctx),
		endpoint:   "https://a.4cdn.org",
	}
}

func (client ApiClient) fetch(method, endpoint string) []byte {
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
