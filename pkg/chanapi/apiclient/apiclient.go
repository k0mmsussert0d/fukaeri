package apiclient

import (
	"time"

	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/restclient"
	"github.com/k0mmsussert0d/fukaeri/pkg/models"
)

type ApiClient struct {
	restClient restclient.RestClient
	limiter    *time.Ticker
}

func (client *ApiClient) Start() {
	client.limiter = time.NewTicker(1 * time.Second)
	client.restClient = *restclient.New()
}

func (client *ApiClient) Threads(board string) models.Threads {
	<-client.limiter.C

	return client.restClient.Threads(board)
}

func (client *ApiClient) Thread(board string, id string) models.Thread {
	<-client.limiter.C

	return client.restClient.Thread(board, id)
}
