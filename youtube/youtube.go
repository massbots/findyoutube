package youtube

import (
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func NewClient(apiKey string) (*Client, error) {
	service, err := youtube.New(&http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	})
	return &Client{
		apiKey:  apiKey,
		Service: service,
	}, err
}

type Client struct {
	apiKey string
	*youtube.Service
}
