package srt

import (
	"fmt"
	"strconv"
	"time"
)

func reverseMap(m map[string]string) map[string]string {
	_m := map[string]string{}
	for k, v := range m {
		_m[v] = k
	}

	return _m
}

func today() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

// tick one second (HHMMSS --> HHMMSS + 1 sec)
func tickSecond(s string) string {
	hour, _ := strconv.Atoi(s[0:2])
	minute, _ := strconv.Atoi(s[2:4])
	second, _ := strconv.Atoi(s[4:6])

	second++
	if second == 60 {
		second = 0
		minute++
	}
	if minute == 60 {
		minute = 0
		hour++
	}
	if hour == 24 {
		hour = 0
	}

	return fmt.Sprintf("%.2d%.2d%.2d", hour, minute, second)
}
