package apiclient

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/limitedhttpclient"
)

type HttpClient interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

type ApiClient struct {
	httpClient    HttpClient
	endpoint      string
	mediaEndpoint string
}

func New() *ApiClient {
	return &ApiClient{
		httpClient:    limitedhttpclient.New(),
		endpoint:      "https://a.4cdn.org",
		mediaEndpoint: "https://i.4cdn.org",
	}
}

func (client ApiClient) fetch(ctx context.Context, method, endpoint string, dst *[]byte) (err error) {
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

	err = client.fetchRequest(ctx, req, dst)
	if err != nil {
		return err
	}

	return nil
}

func (client ApiClient) fetchRequest(ctx context.Context, request *http.Request, dst *[]byte) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = internal.ConvertPanicToError(r)
		}
	}()

	resp, err := client.httpClient.Do(ctx, request)
	internal.HandleError(err)

	defer resp.Body.Close()
	*dst, err = ioutil.ReadAll(resp.Body)
	internal.HandleError(err)

	return nil
}
