package srt

import (
	"fmt"
	"regexp"
)

const (
	srtScheme = "https"
	srtHost   = "app.srail.co.kr"
	srtPort   = 443
)

var srtMobile string = fmt.Sprintf("%s://%s:%d", srtScheme, srtHost, srtPort)

var (
	srtMainURL           = srtMobile + "/main/main.do"
	srtLoginURL          = srtMobile + "/apb/selectListApb01080_n.do"
	srtLogoutURL         = srtMobile + "/login/loginOut.do"
	srtSearchScheduleURL = srtMobile + "/ara/selectListAra10007_n.do"
	srtReserveURL        = srtMobile + "/arc/selectListArc05013_n.do"
	srtReservationsURL   = srtMobile + "/atc/selectListAtc14016_n.do"
	srtTicketInfoURL     = srtMobile + "/ard/selectListArd02017_n.do"
	srtCancelURL         = srtMobile + "/ard/selectListArd02045_n.do"
)

var (
	defaultHeaders map[string]string = map[string]string{
		"User-Agent": "Mozilla/5.0 (Linux; Android 5.1.1; LGM-V300K Build/N2G47H) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/39.0.0.0 Mobile Safari/537.36SRT-APP-Android V.1.0.6",
		"Accept":     "application/json",
	}
)

var (
	regexEmail, _ = regexp.Compile("[^@]+@[^@]+\\.[^@]+")
	regexPhone, _ = regexp.Compile("(\\d{3})-(\\d{3,4})-(\\d{4})")
)

const (
	loginTypeID    = "1"
	loginTypeEmail = "2"
	loginTypePhone = "3"
)

const (
	passengerAdult          = "1"
	passengerDisability1To3 = "2"
	passengerDisability4To6 = "3"
	passengerSenior         = "4"
	passengerChild          = "5"
)

var (
	stationCode map[string]string = map[string]string{
		"수서":      "0551",
		"동탄":      "0552",
		"지제":      "0553",
		"천안아산":    "0502",
		"오송":      "0297",
		"대전":      "0010",
		"공주":      "0514",
		"익산":      "0030",
		"정읍":      "0033",
		"광주송정":    "0036",
		"나주":      "0037",
		"목포":      "0041",
		"김천구미":    "0507",
		"동대구":     "0015",
		"신경주":     "0508",
		"울산(통도사)": "0509",
		"울산":      "0509",
		"통도사":     "0509",
		"부산":      "0020",
	}
	stationName map[string]string = reverseMap(stationCode)
	trainName   map[string]string = map[string]string{
		"00": "KTX",
		"02": "무궁화",
		"03": "통근열차",
		"04": "누리로",
		"05": "전체",
		"07": "KTX-산천",
		"08": "ITX-새마을",
		"09": "ITX-청춘",
		"17": "SRT",
	}

	seatType map[string]string = map[string]string{
		"1": "일반실",
		"2": "특실",
	}

	passengerType map[string]string = map[string]string{
		passengerAdult:          "어른/청소년",
		passengerDisability1To3: "장애 1~3급",
		passengerDisability4To6: "장애 4~6급",
		passengerSenior:         "경로",
		passengerChild:          "어린이",
	}
)
