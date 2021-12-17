package limitedhttpclient

import (
	"context"
	"net/http"

	"golang.org/x/time/rate"
)

type LimitedHttpClient struct {
	httpClient *http.Client
	limiter    *rate.Limiter
	queue      chan http.Request
	ctx        context.Context
}

func New(ctx context.Context) *LimitedHttpClient {
	return &LimitedHttpClient{
		httpClient: &http.Client{},
		queue:      make(chan http.Request, 1),
		limiter:    rate.NewLimiter(rate.Limit(1), 1),
		ctx:        ctx,
	}
}

func (client *LimitedHttpClient) Do(req *http.Request) (*http.Response, error) {
	if err := client.limiter.Wait(client.ctx); err != nil {
		return nil, err
	}

	return client.httpClient.Do(req)
}
