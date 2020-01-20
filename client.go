package srt

import (
	"errors"
	"fmt"
	"net/http"
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
	if !c.isLogin {
		return nil
	}

	resp, err := c.httpClient.R().
		Post(srtLogoutURL)

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("Logout Failed")
	}

	return nil
}

// SearchTrain is used to search trains from SRT server,
// SearchTrain returns *only* available seats,
// if you want to retrieve trains with no seats as well, Use SearchTrainAll() instead
func (c *Client) SearchTrain(dep, arr, date, depTime string) ([]*Train, error) {
	// TODO filtering
	trains, err := c.SearchTrainAll(dep, arr, date, depTime)

	if err != nil {
		return nil, err
	}

	trainsAvailable := make([]*Train, 0)
	for _, train := range trains {
		if train.SeatsAvailable() {
			trainsAvailable = append(trainsAvailable, train)
		}
	}
	return trainsAvailable, nil
}

// SearchTrainAll is used to search *all* trains including trains with no seats
func (c *Client) SearchTrainAll(dep, arr, date, depTime string) ([]*Train, error) {
	if !c.isLogin {
		return nil, errors.New("Not Loggin In")
	}

	depCode, ok := stationCode[dep]
	if !ok {
		return nil, fmt.Errorf("Station `%s` Not Exists", dep)
	}
	arrCode, ok := stationCode[arr]
	if !ok {
		return nil, fmt.Errorf("Station `%s` Not Exists", arr)
	}

	formData := map[string]string{
		// course (1: 직통, 2: 환승, 3: 왕복)
		// TODO: support 환승, 왕복
		"chtnDvCd":   "1",
		"arriveTime": "N",
		"seatAttCd":  "015",
		// 검색 시에는 1명 기준으로 검색
		"psgNum":  "1",
		"trnGpCd": "109",
		// train type (05: 전체, 17: SRT)
		"stlbTrnClsfCd": "05",
		// departure date
		"dptDt": date,
		// departure time
		"dptTm": depTime,
		// arrival station code
		"arvRsStnCd": arrCode,
		// departure station code
		"dptRsStnCd": depCode,
	}

	resp, err := c.httpClient.R().
		SetFormData(formData).
		Post(srtSearchScheduleURL)

	if err != nil {
		return nil, err
	}

	parser := &responseParser{}
	err = parser.Parse(resp.Body())

	if err != nil {
		return nil, err
	}
	if !parser.Success() {
		c.debug(string(resp.Body()))
		return nil, errors.New("Response Parsing Failed")
	}

	trainsData := parser.
		Data()["outDataSets"].(map[string]interface{})["dsOutput1"].([]interface{})

	toTrain := func(t map[string]interface{}) *Train {
		return &Train{
			trainCode:        t["stlbTrnClsfCd"].(string),
			trainName:        trainName[t["stlbTrnClsfCd"].(string)],
			trainNumber:      t["trnNo"].(string),
			depDate:          t["dptDt"].(string),
			depTime:          t["dptTm"].(string),
			depStationCode:   t["dptRsStnCd"].(string),
			depStationName:   stationName[t["dptRsStnCd"].(string)],
			arrDate:          t["arvDt"].(string),
			arrTime:          t["arvTm"].(string),
			arrStationCode:   t["arvRsStnCd"].(string),
			arrStationName:   stationName[t["arvRsStnCd"].(string)],
			generalSeatState: t["gnrmRsvPsbStr"].(string),
			specialSeatState: t["sprmRsvPsbStr"].(string),
		}
	}

	trains := make([]*Train, 0)
	for _, t := range trainsData {
		trains = append(trains, toTrain(t.(map[string]interface{})))
	}

	// Note: updated api uses pagination,
	//      therefore, to retreive all trains, retry searching unless there are no more trains
	for len(trains) > 0 {
		nextDepTime := tickSecond(trains[len(trains)-1].depTime)
		formData["dptTm"] = nextDepTime

		resp, err := c.httpClient.R().
			SetFormData(formData).
			Post(srtSearchScheduleURL)

		if err != nil {
			return nil, err
		}

		parser := &responseParser{}
		err = parser.Parse(resp.Body())

		if err != nil {
			return nil, err
		}

		if !parser.Success() {
			break
		}

		trainsData := parser.
			Data()["outDataSets"].(map[string]interface{})["dsOutput1"].([]interface{})

		for _, t := range trainsData {
			trains = append(trains, toTrain(t.(map[string]interface{})))
		}
	}

	return trains, nil
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
