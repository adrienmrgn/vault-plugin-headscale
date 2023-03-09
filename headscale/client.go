package headscale

import (
	"net/http"
	"time"
)

// Client : struct that define a Headscale Client
type Client struct {
	APIURL string
	APIKey string
	HTTP   *http.Client
}

// NewClient returns a Headscale client and instantiatethe http client
func NewClient() *Client {
	return &Client{
		APIURL: "",
		APIKey: "",
		HTTP: &http.Client{
			Timeout: time.Minute,
		},
	}
}
