package headscale

import (
	"net/http"
	"time"
)

type Client struct {
	ApiURL string
	ApiKey string
	HTTP	 *http.Client
}

func newClient() *Client {
	return &Client{
		ApiURL: "",
		ApiKey: "",
		HTTP: &http.Client{
			Timeout: time.Minute,
		},
	}
}

