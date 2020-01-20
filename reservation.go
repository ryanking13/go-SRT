package srt

import "fmt"

// Ticket represents a SRT Ticket
type Ticket struct {
	car               string
	seat              string
	seatTypeCode      string
	seatType          string
	passengerTypeCode string
	passengerType     string
	price             string
	originalPrice     string
	discount          string
}

func (t *Ticket) String() string {
	return fmt.Sprintf(
		"%s호차 %s (%s) %s [%s원(%s원 할인)]",
		t.car,
		t.seat,
		t.seatType,
		t.passengerType,
		t.price,
		t.discount,
	)
}

// Reservation represents a reservation of SRT train.
// Reservation consists of train information and multiple Tickets
type Reservation struct {
	reservationNumber string
	totalCost         string
	seatCount         string
	trainCode         string
	trainName         string
	trainNumber       string
	depDate           string
	depTime           string
	depStationCode    string
	depStationName    string
	arrTime           string
	arrStationCode    string
	arrStationName    string
	paymentDate       string
	paymentTime       string
	tickets           []*Ticket
}

func (r *Reservation) String() string {
	return fmt.Sprintf(
		"[%s] %s월 %s일, %s~%s(%s:%s~%s:%s) %s원(%s석), 구입기한 %s월 %s일 %s:%s",
		r.trainName,
		r.depDate[4:6],
		r.depDate[6:8],
		r.depStationName,
		r.arrStationName,
		r.depTime[0:2],
		r.depTime[2:4],
		r.arrTime[0:2],
		r.arrTime[2:4],
		r.totalCost,
		r.seatCount,
		r.paymentDate[4:6],
		r.paymentDate[6:8],
		r.paymentTime[0:2],
		r.paymentTime[2:4],
	)
}
