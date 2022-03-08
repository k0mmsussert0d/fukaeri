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
