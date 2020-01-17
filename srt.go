package srt // import "github.com/ryanking13/go-SRT"

import "github.com/go-resty/resty/v2"

// New method creates a new SRT client.
func New() *Client {
	client := &Client{
		httpClient: resty.New(),
	}

	client.httpClient.SetHeaders(defaultHeaders)
	return client
}
