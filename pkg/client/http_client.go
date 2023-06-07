package client

import (
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type httpClient struct {
	client *http.Client
}

func (dc *httpClient) Do(r *http.Request) (*http.Response, error) {
	return dc.client.Do(r)
}

func NewHTTPClient(timeoutInSec int) HTTPClient {
	return &httpClient{
		client: &http.Client{
			Timeout: time.Second * time.Duration(timeoutInSec),
		},
	}
}
