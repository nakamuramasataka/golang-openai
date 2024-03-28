package openai

import (
	"net/http"
)

type Client struct {
	HTTPClient *http.Client
	APIKey     string
}

func NewClient(apiKey string) *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		APIKey:     apiKey,
	}
}
