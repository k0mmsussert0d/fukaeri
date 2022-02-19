package limitedhttpclient

import (
	"context"
	"net/http"

	"github.com/k0mmsussert0d/fukaeri/internal/log"
	"golang.org/x/time/rate"
)

type LimitedHttpClient struct {
	httpClient *http.Client
	limiter    *rate.Limiter
	queue      chan http.Request
}

func New() *LimitedHttpClient {
	return &LimitedHttpClient{
		httpClient: &http.Client{},
		queue:      make(chan http.Request, 1),
		limiter:    rate.NewLimiter(rate.Limit(1), 1),
	}
}

func (client *LimitedHttpClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	log.Debug().Printf("Request %v %v queued. Waiting...", req.Method, req.URL)
	select {
	case <-client.rateLimiterChan(ctx):
		log.Debug().Printf("Sending request %v %v", req.Method, req.URL)
		return client.httpClient.Do(req)
	case <-ctx.Done():
		log.Debug().Printf("Request %v %v aborted", req.Method, req.URL)
		return nil, ctx.Err()
	}
}

func (client *LimitedHttpClient) rateLimiterChan(ctx context.Context) <-chan interface{} {
	readyStream := make(chan interface{}, 1)
	go func() {
		defer close(readyStream)
		for {
			if err := client.limiter.Wait(ctx); err != nil {
				return
			} else {
				readyStream <- 0
			}
		}
	}()
	return readyStream
}
