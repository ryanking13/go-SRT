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

// SearchTrain is used to search trains from SRT server
func (c *Client) SearchTrain(params *SearchParams) ([]*Train, error) {
	if !c.isLogin {
		return nil, errors.New("Not Logged In")
	}

	if params.Date == "" {
		params.Date = today()
	}

	if params.Time == "" {
		params.Time = "000000"
	}

	depCode, ok := stationCode[params.Dep]
	if !ok {
		return nil, fmt.Errorf("Station `%s` Not Exists", params.Dep)
	}
	arrCode, ok := stationCode[params.Arr]
	if !ok {
		return nil, fmt.Errorf("Station `%s` Not Exists", params.Arr)
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
		"dptDt": params.Date,
		// departure time
		"dptTm": params.Time,
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

	if !params.IncludeSoldOut {
		trainsAvailable := make([]*Train, 0)
		for _, train := range trains {
			if train.SeatsAvailable() {
				trainsAvailable = append(trainsAvailable, train)
			}
		}
		return trainsAvailable, nil
	}

	return trains, nil
}

// Reserve is used to reserve SRT train
func (c *Client) Reserve(r *ReserveParams) (*Reservation, error) {
	if !c.isLogin {
		return nil, errors.New("Not Logged In")
	}

	if r.Train.trainName != "SRT" {
		return nil, fmt.Errorf("SRT is expected for a train name, %s given", r.Train.trainName)
	}

	if r.Passengers == nil {
		r.Passengers = []*Passenger{Adult(1)}
	}

	formData := map[string]string{
		"reserveType":    "11",
		"jobId":          "1101", // 개인 예약
		"jrnyCnt":        "1",
		"jrnyTpCd":       "11",
		"jrnySqno1":      "001",
		"stndFlg":        "N",
		"trnGpCd1":       "300", // 열차그룹코드 (좌석선택은 SRT만 가능하기때문에 무조건 300을 셋팅한다)
		"stlbTrnClsfCd1": r.Train.trainCode,
		"dptDt1":         r.Train.depDate,
		"dptTm1":         r.Train.depTime,
		"runDt1":         r.Train.depDate,
		"trnNo1":         fmt.Sprintf("%05s", r.Train.trainNumber),
		"dptRsStnCd1":    r.Train.depStationCode,
		"dptRsStnCdNm1":  r.Train.depStationName,
		"arvRsStnCd1":    r.Train.arrStationCode,
		"arvRsStnCdNm1":  r.Train.arrStationName,
	}

	for k, v := range passengers2Params(r.Passengers, r.SpecialSeatOnly) {
		formData[k] = v
	}

	resp, err := c.httpClient.R().
		SetFormData(formData).
		Post(srtReserveURL)

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

	dupMsg := "요청하신 승차권과 동일한 시간대에 예약 또는 발권하신 승차권이 존재합니다."
	if strings.Contains(parser.String(), dupMsg) {
		c.log("WARNING: 이미 같은 시간대의 예약이 존재합니다. 중복 예약되었습니다.")
	}

	c.debug(parser.String())

	reservationNo := parser.Data()["reservListMap"].([]interface{})[0].(map[string]interface{})["pnrNo"].(string)
	reservations, err := c.Reservations()

	if err != nil {
		return nil, err
	}

	for _, reservation := range reservations {
		if reservation.reservationNumber == reservationNo {
			return reservation, nil
		}
	}

	return nil, errors.New("Ticket not found: check reservation status")
}

// Reservations is used to retrieve your SRT reservations
func (c *Client) Reservations() ([]*Reservation, error) {
	if !c.isLogin {
		return nil, errors.New("Not Logged In")
	}

	formData := map[string]string{
		"pageNo": "0",
	}

	resp, err := c.httpClient.R().
		SetFormData(formData).
		Post(srtReservationsURL)

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

	trainData := parser.Data()["trainListMap"].([]interface{})
	payData := parser.Data()["payListMap"].([]interface{})
	reservations := make([]*Reservation, 0)

	for i, _train := range trainData {
		train := _train.(map[string]interface{})
		pay := payData[i].(map[string]interface{})
		tickets, err := c.TicketsByNumber(train["pnrNo"].(string))
		if err != nil {
			return nil, err
		}

		reservations = append(reservations, &Reservation{
			reservationNumber: train["pnrNo"].(string),
			totalCost:         fmt.Sprintf("%.0f", train["rcvdAmt"].(float64)),
			seatCount:         fmt.Sprintf("%.0f", train["tkSpecNum"].(float64)),
			trainCode:         pay["stlbTrnClsfCd"].(string),
			trainName:         trainName[pay["stlbTrnClsfCd"].(string)],
			trainNumber:       pay["trnNo"].(string),
			depDate:           pay["dptDt"].(string),
			depTime:           pay["dptTm"].(string),
			depStationCode:    pay["dptRsStnCd"].(string),
			depStationName:    stationName[pay["dptRsStnCd"].(string)],
			arrTime:           pay["arvTm"].(string),
			arrStationCode:    pay["arvRsStnCd"].(string),
			arrStationName:    stationName[pay["arvRsStnCd"].(string)],
			paymentDate:       pay["iseLmtDt"].(string),
			paymentTime:       pay["iseLmtTm"].(string),
			tickets:           tickets,
		})
	}
	return reservations, nil
}

// Tickets is used to see information of the reservation (reservation consists of multiple tickets)
func (c *Client) Tickets(reservation *Reservation) ([]*Ticket, error) {
	return c.TicketsByNumber(reservation.reservationNumber)
}

// TicketsByNumber is used to see information of the reservation by it's number directly
func (c *Client) TicketsByNumber(reservationNumber string) ([]*Ticket, error) {
	if !c.isLogin {
		return nil, errors.New("Not Logged In")
	}

	formData := map[string]string{
		"pnrNo":    reservationNumber,
		"jrnySqno": "1",
	}

	resp, err := c.httpClient.R().
		SetFormData(formData).
		Post(srtTicketInfoURL)

	if err != nil {
		return nil, err
	}

	parser := &responseParser{}
	err = parser.Parse(resp.Body())

	if err != nil {
		return nil, err
	}

	if !parser.Success() {
		return nil, errors.New("Response Parsing Failed")
	}

	tickets := make([]*Ticket, 0)
	_ticketList := parser.Data()["trainListMap"].([]interface{})
	for _, ticket := range _ticketList {
		t := ticket.(map[string]interface{})
		tickets = append(tickets, &Ticket{
			car:               fmt.Sprintf("%.0f", t["scarNo"].(float64)),
			seat:              t["seatNo"].(string),
			seatTypeCode:      t["psrmClCd"].(string),
			seatType:          seatType[t["psrmClCd"].(string)],
			passengerTypeCode: t["psgTpCd"].(string),
			passengerType:     passengerType[t["psgTpCd"].(string)],
			price:             strings.TrimLeft(t["rcvdAmt"].(string), "0"),
			originalPrice:     strings.TrimLeft(t["stdrPrc"].(string), "0"),
			discount:          strings.TrimLeft(t["dcntPrc"].(string), "0"),
		})
	}

	return tickets, nil
}

// Cancel is used to cancel reservation
func (c *Client) Cancel(reservation *Reservation) error {
	return c.CancelByNumber(reservation.reservationNumber)
}

// CancelByNumber is used to cancel reservation by it's number directly
func (c *Client) CancelByNumber(reservationNumber string) error {
	if !c.isLogin {
		return errors.New("Not Logged In")
	}

	formData := map[string]string{
		"pnrNo":     reservationNumber,
		"jrnyCnt":   "1",
		"rsvChgTno": "0",
	}

	resp, err := c.httpClient.R().
		SetFormData(formData).
		Post(srtCancelURL)

	if err != nil {
		return err
	}

	parser := &responseParser{}
	err = parser.Parse(resp.Body())

	if err != nil {
		return err
	}
	if !parser.Success() {
		c.debug(string(resp.Body()))
		return errors.New("Response Parsing Failed")
	}

	c.debug(parser.String())

	return nil
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
