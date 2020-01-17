package srt

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

// Client struct is used to create a SRT client.
type Client struct {
	httpClient *resty.Client
}

// Login is used to login to SRT server
func (c *Client) Login() error {
	return errors.New("Not Implemented")
}

// Logout is used to logout from SRT server
func (c *Client) Logout() error {
	return errors.New("Not Implemented")
}

// SearchTrain is used to search trains from SRT server
func (c *Client) SearchTrain() error {
	return errors.New("Not Implemented")
}

// Reserve is used to reserve SRT train
func (c *Client) Reserve() error {
	return errors.New("Not Implemented")
}

// Reservations is used to retrieve your SRT reservations
func (c *Client) Reservations() error {
	return errors.New("Not Implemented")
}

// TicketInfo is used to see information of the ticket
func (c *Client) TicketInfo() error {
	return errors.New("Not Implemented")
}

// Cancel is used to cancel reservation
func (c *Client) Cancel() error {
	return errors.New("Not Implemented")
}
