package srt

import "fmt"

const (
	srtScheme = "https"
	srtHost   = "app.srail.co.kr"
	srtPort   = 443
)

var srtMobile string = fmt.Sprintf("%s://%s:%d", srtScheme, srtHost, srtPort)

var (
	srtMainURL           = srtMobile + "/neo/main/main.do"
	srtLoginURL          = srtMobile + "/neo/apb/selectListApb01080_n.do"
	srtLogoutURL         = srtMobile + "/neo/login/loginOut.do"
	srtSearchScheduleURL = srtMobile + "/neo/ara/selectListAra10007_n.do"
	srtReserveURL        = srtMobile + "/neo/arc/selectListArc05013_n.do"
	srtTicketsURL        = srtMobile + "/neo/atc/selectListAtc14016_n.do"
	srtTicketInfoURL     = srtMobile + "/neo/ard/selectListArd02017_n.do"
	srtCancelURL         = srtMobile + "/neo/ard/selectListArd02045_n.do"
)

var (
	defaultHeaders map[string]string = map[string]string{
		"User-Agent": "Mozilla/5.0 (Linux; Android 5.1.1; LGM-V300K Build/N2G47H) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/39.0.0.0 Mobile Safari/537.36SRT-APP-Android V.1.0.6",
		"Accept":     "application/json",
	}
)
