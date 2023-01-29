package mockedhttpclient

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/k0mmsussert0d/fukaeri/internal/log"
)

type Request struct {
	method   string
	endpoint string
}

type Response struct {
	code    int
	body    []byte
	headers map[string][]string
}

type MockedHttpClient struct {
	requests map[Request]Response
}

type MockedHttpClientBuilderRequestPart struct {
	client  *MockedHttpClient
	request *Request
}

func New() *MockedHttpClient {
	return &MockedHttpClient{
		requests: make(map[Request]Response),
	}
}

func (client *MockedHttpClient) On(method, endpoint string) *MockedHttpClientBuilderRequestPart {
	return &MockedHttpClientBuilderRequestPart{
		client: client,
		request: &Request{
			method:   method,
			endpoint: endpoint,
		},
	}
}

func (builder *MockedHttpClientBuilderRequestPart) Return(code int, body []byte, headers map[string][]string) *MockedHttpClient {
	builder.client.requests[*builder.request] = Response{
		code:    code,
		body:    body,
		headers: headers,
	}
	return builder.client
}

var httpStatuses = map[int]string{
	200: "200 OK",
	201: "201 Created",
	202: "202 Accepted",
	204: "204 No Content",
	301: "301 Moved Permanently",
	302: "302 Found",
	304: "304 Not Modified",
	307: "307 Temporary Redirect",
	308: "308 Permanent Redirect",
	400: "400 Bad Request",
	401: "401 Unauthrized",
	403: "403 Forbidden",
	404: "404 Not Found",
	405: "405 Method Not Allowed",
	408: "408 Request Timeout",
	411: "411 Length Required",
	413: "413 Payload Too Large",
	414: "414 Request-URI Too Long",
	415: "415 Unsupported Media Type",
	429: "429 Too Many Requests",
	500: "500 Internal Server Error",
	502: "502 Bad Gateway",
	503: "503 Service Unavailable",
	504: "504 Gateway Timeout",
}

func (client *MockedHttpClient) Do(ctx context.Context, req *http.Request) (*http.Response, []byte, error) {
	method, host := req.Method, req.URL.String()
	response, exists := client.requests[Request{method, host}]
	if !exists {
		log.Logger().Panicw("Response to the request has not been defined",
			"method", method,
			"host", host,
		)
	}

	responseStatusString, exists := httpStatuses[response.code]
	if !exists {
		responseStatusString = fmt.Sprintf("%v", response.code)
	}

	httpResponse := &http.Response{
		Status:        responseStatusString,
		StatusCode:    response.code,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        response.headers,
		Body:          ioutil.NopCloser(bytes.NewReader(response.body)),
		ContentLength: int64(len(response.body)),
		Close:         false,
		Uncompressed:  true,
		Request:       req,
	}

	return httpResponse, response.body, nil
}
