package srt

import (
	"fmt"
	"strconv"
)

type Passenger struct {
	count    int
	name     string
	typeCode string
}

func (p *Passenger) String() string {
	return fmt.Sprintf("%s %d명", p.name, p.count)
}

func Adult(count int) *Passenger {
	return &Passenger{
		count:    count,
		name:     passengerType[passengerAdult],
		typeCode: passengerAdult,
	}
}

func Disability1To3(count int) *Passenger {
	return &Passenger{
		count:    count,
		name:     passengerType[passengerDisability1To3],
		typeCode: passengerDisability1To3,
	}
}

func Disability4To6(count int) *Passenger {
	return &Passenger{
		count:    count,
		name:     passengerType[passengerDisability4To6],
		typeCode: passengerDisability4To6,
	}
}

func Senior(count int) *Passenger {
	return &Passenger{
		count:    count,
		name:     passengerType[passengerSenior],
		typeCode: passengerSenior,
	}
}

func Child(count int) *Passenger {
	return &Passenger{
		count:    count,
		name:     passengerType[passengerChild],
		typeCode: passengerChild,
	}
}

func passengers2Params(passengers []*Passenger, specialSeatOnly bool) map[string]string {
	totalCnt := 0
	for _, p := range passengers {
		totalCnt += p.count
	}

	params := map[string]string{
		"totPrnb":    strconv.Itoa(totalCnt),
		"psgGridcnt": strconv.Itoa(len(passengers)),
	}

	for i, p := range passengers {
		params[fmt.Sprintf("psgTpCd%d", i+1)] = p.typeCode
		params[fmt.Sprintf("psgInfoPerPrnb%d", i+1)] = strconv.Itoa(p.count)
		// seat location ('000': 기본, '012': 창측, '013': 복도측)
		// TODO: 선택 가능하게
		params[fmt.Sprintf("locSeatAttCd%d", i+1)] = "000"
		// seat requirement ('015': 일반, '021': 휠체어)
		// TODO: 선택 가능하게
		params[fmt.Sprintf("rqSeatAttCd%d", i+1)] = "015"
		// seat direction ('009': 정방향)
		params[fmt.Sprintf("dirSeatAttCd%d", i+1)] = "009"
		params[fmt.Sprintf("smkSeatAttCd%d", i+1)] = "000"
		params[fmt.Sprintf("etcSeatAttCd%d", i+1)] = "000"
		// seat type: ('1': 일반실, '2': 특실)
		if specialSeatOnly {
			params[fmt.Sprintf("psrmClCd%d", i+1)] = "2"
		} else {
			params[fmt.Sprintf("psrmClCd%d", i+1)] = "1"
		}
	}
	return params
}
