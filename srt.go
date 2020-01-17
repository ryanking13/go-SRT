package srt // import "github.com/ryanking13/go-SRT"

import "github.com/go-resty/resty/v2"

// Client is a SRT client
type Client struct {
	client *resty.Client
}
