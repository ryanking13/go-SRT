package srt

import (
	"fmt"
	"strings"
)

// Train struct is used to represent a SRT Train
type Train struct {
	trainCode        string
	trainName        string
	trainNumber      string
	depDate          string
	depTime          string
	depStationCode   string
	depStationName   string
	arrDate          string
	arrTime          string
	arrStationCode   string
	arrStationName   string
	generalSeatState string
	specialSeatState string
}

// GeneralSeatsAvailable is used to check there are empty general seats
func (t *Train) GeneralSeatsAvailable() bool {
	return strings.Contains(t.generalSeatState, "예약가능")
}

// SpecialSeatsAvailable is used to check there are empty special seats
func (t *Train) SpecialSeatsAvailable() bool {
	return strings.Contains(t.specialSeatState, "예약가능")
}

// SeatsAvailable is used to check there are empty seats
func (t *Train) SeatsAvailable() bool {
	return t.GeneralSeatsAvailable() || t.SpecialSeatsAvailable()
}

func (t *Train) String() string {
	return fmt.Sprintf(
		"[%s] %s월 %s일, %s~%s(%s:%s~%s:%s) 특실: %s, 일반실: %s",
		t.trainName,
		t.depDate[4:6],
		t.depDate[6:8],
		t.depStationName,
		t.arrStationName,
		t.depTime[0:2],
		t.depTime[2:4],
		t.arrTime[0:2],
		t.arrTime[2:4],
		t.specialSeatState,
		t.generalSeatState,
	)
}
