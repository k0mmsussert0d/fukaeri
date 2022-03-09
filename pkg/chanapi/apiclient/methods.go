package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/pkg/models"
)

func (client ApiClient) Threads(ctx context.Context, board string) (*models.Threads, error) {
	res, err := client.fetch(ctx, "GET", fmt.Sprintf("%s/%s/threads.json", client.Endpoint, board))
	if err != nil {
		return nil, internal.WrapError(err, "Failed to fetch %s threadslist", board)
	}
	var returnObj models.Threads
	json.Unmarshal(res, &returnObj)
	return &returnObj, nil
}

func (client ApiClient) Thread(ctx context.Context, board string, id string) (*models.Thread, error) {
	res, err := client.fetch(ctx, "GET", fmt.Sprintf("%s/%s/thread/%s.json", client.Endpoint, board, id))
	if err != nil {
		return nil, internal.WrapError(err, "Failed to fetch thread %s/%s", board, id)
	}
	var returnObj models.Thread
	json.Unmarshal(res, &returnObj)
	return &returnObj, nil
}

func (client ApiClient) ThreadSince(ctx context.Context, board, id string, since time.Time) (*models.Thread, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/thread/%s.json", client.Endpoint, board, id), nil)
	if err != nil {
		return nil, internal.WrapError(err, "Failed to fetch thread %s/%s since %v", board, id, since)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("If-Modified-Since", since.Local().Format(http.TimeFormat))

	resp, body, err := client.HttpClient.Do(ctx, req)
	if err != nil {
		return nil, internal.WrapError(err, "Failed to fetch thread %s/%s since %v", board, id, since)
	}

	if resp.StatusCode == 304 {
		return nil, nil
	}

	var returnObj models.Thread
	json.Unmarshal(body, &returnObj)
	return &returnObj, nil
}

func (client ApiClient) File(ctx context.Context, board string, timestamp int64, ext string) ([]byte, error) {
	res, err := client.fetch(ctx, "GET", fmt.Sprintf("%s/%s/%d%s", client.MediaEndpoint, board, timestamp, ext))
	if err != nil {
		return nil, internal.WrapError(err, "Failed to fetch file %s/%v.%s", board, timestamp, ext)
	}
	return res, nil
}
