package srt // import "github.com/ryanking13/go-SRT"

import (
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// New method creates a new SRT client.
func New() *Client {
	client := &Client{
		httpClient: resty.New(),
		logger: &logrus.Logger{
			Out: os.Stdout,
			Formatter: &logrus.TextFormatter{
				DisableTimestamp: true,
			},
			Level: logrus.InfoLevel,
		},
	}

	client.httpClient.SetHeaders(defaultHeaders)
	return client
}
