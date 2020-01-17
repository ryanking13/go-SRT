package srt

import (
	"errors"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Client struct is used to create a SRT client.
type Client struct {
	httpClient *resty.Client
	isLogin    bool
	logger     *logrus.Logger
}

// Login is used to login to SRT server
func (c *Client) Login(id, pw string) error {

	var loginType string
	if regexEmail.MatchString(id) {
		c.debug("Setting login type to email")
		loginType = loginTypeEmail
	} else if regexPhone.MatchString(id) {
		c.debug("Setting login type to phone number")
		id = strings.ReplaceAll(id, "-", "")
		loginType = loginTypePhone
	} else {
		c.debug("Setting login type to membership id")
		loginType = loginTypeID
	}

	resp, err := c.httpClient.R().
		SetFormData(map[string]string{
			"auto":          "Y",
			"check":         "Y",
			"page":          "menu",
			"deviceKey":     "-",
			"customerYn":    "",
			"login_referer": srtMainURL,
			"srchDvCd":      loginType,
			"srchDvNm":      id,
			"hmpgPwdCphd":   pw,
		}).
		Post(srtLoginURL)

	if err != nil {
		return err
	}

	body := resp.String()
	c.debug(body)
	if strings.Contains(body, "존재하지않는 회원입니다") {
		return errors.New("존재하지 않는 회원입니다. ID를 확인하세요")
	}
	if strings.Contains(body, "비밀번호 오류입니다.") {
		return errors.New("비밀번호 오류. 비밀번호를 확인하세요")
	}

	c.isLogin = true
	return nil
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

// SetDebug is used to change SRT Client logger setting to print all debug logs
func (c *Client) SetDebug() {
	c.logger.Level = logrus.DebugLevel
}

func (c *Client) debug(msg string) {
	c.logger.Debug(msg)
}

func (c *Client) log(msg string) {
	c.logger.Info(msg)
}
