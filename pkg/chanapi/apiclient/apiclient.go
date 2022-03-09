package apiclient

import (
	"context"
	"net/http"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/limitedhttpclient"
)

type HttpClient interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, []byte, error)
}

type ApiClient struct {
	HttpClient    HttpClient
	Endpoint      string
	MediaEndpoint string
}

func New() *ApiClient {
	return &ApiClient{
		HttpClient:    limitedhttpclient.New(),
		Endpoint:      "https://a.4cdn.org",
		MediaEndpoint: "https://i.4cdn.org",
	}
}

func (client ApiClient) fetch(ctx context.Context, method, endpoint string) (body []byte, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = internal.ConvertPanicToError(r)
		}
	}()

	req, err := http.NewRequest(method, endpoint, nil)
	internal.HandleError(err)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	body, err = client.fetchRequest(ctx, req)
	return
}

func (client ApiClient) fetchRequest(ctx context.Context, request *http.Request) (body []byte, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = internal.ConvertPanicToError(r)
		}
	}()

	_, body, err = client.HttpClient.Do(ctx, request)
	internal.HandleError(err)
	return
}
