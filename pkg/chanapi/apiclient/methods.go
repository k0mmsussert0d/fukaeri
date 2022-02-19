package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/k0mmsussert0d/fukaeri/internal/log"
	"github.com/k0mmsussert0d/fukaeri/pkg/models"
)

func (client ApiClient) Threads(ctx context.Context, board string) (*models.Threads, error) {
	var res []byte
	err := client.fetch(ctx, "GET", fmt.Sprintf("%s/%s/threads.json", client.endpoint, board), &res)
	if err != nil {
		return nil, err
	}
	var returnObj models.Threads
	json.Unmarshal(res, &returnObj)
	return &returnObj, nil
}

func (client ApiClient) Thread(ctx context.Context, board string, id string) (*models.Thread, error) {
	var res []byte
	err := client.fetch(ctx, "GET", fmt.Sprintf("%s/%s/thread/%s.json", client.endpoint, board, id), &res)
	if err != nil {
		return nil, err
	}
	var returnObj models.Thread
	json.Unmarshal(res, &returnObj)
	return &returnObj, nil
}

func (client ApiClient) ThreadSince(ctx context.Context, board, id string, since time.Time) (*models.Thread, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/thread/%s.json", client.endpoint, board, id), nil)
	if err != nil {
		log.Error().Println(err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("If-Modified-Since", since.Local().Format(http.TimeFormat))

	resp, err := client.httpClient.Do(ctx, req)
	if err != nil {
		log.Error().Println(err)
		return nil, err
	}

	if resp.StatusCode == 304 {
		return nil, nil
	}

	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Println(err)
	}
	var returnObj models.Thread
	json.Unmarshal(res, &returnObj)
	return &returnObj, nil
}

func (client ApiClient) File(ctx context.Context, board string, timestamp int64, ext string) (*[]byte, error) {
	var res []byte
	err := client.fetch(ctx, "GET", fmt.Sprintf("%s/%s/%d%s", client.mediaEndpoint, board, timestamp, ext), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
