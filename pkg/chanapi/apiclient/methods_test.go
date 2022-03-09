package apiclient_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"github.com/k0mmsussert0d/fukaeri/pkg/chanapi/apiclient"
	mockedhttpclient "github.com/k0mmsussert0d/fukaeri/pkg/chanapi/mocked_http_client"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

var apiClient *apiclient.ApiClient = nil
var mockedHttpClient *mockedhttpclient.MockedHttpClient = nil

var corsHeaders = map[string][]string{
	"access-control-allow-headers": {"If-Modified-Since"},
	"access-control-allow-methods": {"GET", "OPTIONS"},
	"access-control-allow-origin":  {"https://example.com"},
	"access-control-max-age":       {"1728000"},
}

func TestMain(m *testing.M) {
	mockedHttpClient = mockedhttpclient.New()
	apiClient = &apiclient.ApiClient{
		HttpClient:    mockedHttpClient,
		Endpoint:      "https://a.example.com",
		MediaEndpoint: "https://i.example.com",
	}
	code := m.Run()
	os.Exit(code)
}

func TestThreads(t *testing.T) {
	responseBody, err := os.ReadFile("./test_bodies/threads.json")
	if err != nil {
		internal.HandleError(err)
	}
	mockedHttpClient.On("GET", "https://a.example.com/po/threads.json").Return(200, responseBody, headers(map[string][]string{}))

	threads, err := apiClient.Threads(context.TODO(), "po")

	assert.NilError(t, err)
	assert.Equal(t, len(*threads), 10)
	assert.Equal(t, (*threads)[0].Threads[0].No, 570368)
	assert.Equal(t, (*threads)[0].Threads[0].LastModified, 1546294897)
	assert.Equal(t, (*threads)[0].Threads[0].Replies, 2)

	for idx, page := range *threads {
		assert.Equal(t, page.Page, idx+1)
		assert.Equal(t, len(page.Threads), 15)
	}
}

func TestThread(t *testing.T) {
	responseBody, err := os.ReadFile("./test_bodies/570368.json")
	if err != nil {
		internal.HandleError(err)
	}
	mockedHttpClient.On("GET", "https://a.example.com/po/thread/570368.json").Return(200, responseBody, headers(map[string][]string{}))

	thread, err := apiClient.Thread(context.TODO(), "po", "570368")

	assert.NilError(t, err)
	assert.Equal(t, len(thread.Posts), 3)
	assert.Equal(t, thread.Posts[0].No, 570368)
	assert.Equal(t, thread.Posts[0].Sub, "Welcome to /po/!")
	assert.Equal(t, thread.Posts[0].Filename, "yotsuba_folding")
	assert.Equal(t, thread.Posts[0].Md5, "uZUeZeB14FVR+Mc2ScHvVA==")
}

func TestThreadSince(t *testing.T) {
	responseBody, err := os.ReadFile("./test_bodies/570368.json")
	internal.HandleError(err)
	lastModified := time.Now().Add(time.Duration(5) * time.Minute).Format(http.TimeFormat)

	t.Run("No new posts since", func(t *testing.T) {
		mockedHttpClient.On("GET", "https://a.example.com/po/thread/570368.json").Return(304, responseBody, headers(map[string][]string{"last-modified": {lastModified}}))

		thread, err := apiClient.ThreadSince(context.TODO(), "po", "570368", time.Now())

		assert.NilError(t, err)
		assert.Assert(t, is.Nil(thread))
	})

	t.Run("New posts since", func(t *testing.T) {
		mockedHttpClient.On("GET", "https://a.example.com/po/thread/570368.json").Return(200, responseBody, headers(map[string][]string{"last-modified": {lastModified}}))

		thread, err := apiClient.ThreadSince(context.TODO(), "po", "570368", time.Now())

		assert.NilError(t, err)
		assert.Equal(t, len(thread.Posts), 3)
		assert.Equal(t, thread.Posts[0].No, 570368)
		assert.Equal(t, thread.Posts[0].Sub, "Welcome to /po/!")
		assert.Equal(t, thread.Posts[0].Filename, "yotsuba_folding")
		assert.Equal(t, thread.Posts[0].Md5, "uZUeZeB14FVR+Mc2ScHvVA==")
	})
}

func headers(headers map[string][]string) map[string][]string {
	for k, v := range corsHeaders {
		headers[k] = v
	}
	now := time.Now().Format(http.TimeFormat)
	headers["date"] = []string{now}
	_, lmPresent := headers["last-modified"]
	if !lmPresent {
		headers["last-modified"] = []string{now}
	}
	return headers
}
