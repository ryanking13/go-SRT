package main

import (
	"fmt"
	"time"
)

// Time struct is used to represent time parameter of SRT
type Time struct {
	time time.Time
}

// New is used to instantiate new Time instance
func today() *Time {
	return &Time{time: time.Now()}
}

// NextDay is used to retrieve next day of Time
func (t *Time) NextDay(d int) *Time {
	return &Time{time: t.time.AddDate(0, 0, d)}
}

func (t *Time) String() string {
	return fmt.Sprintf("%d%02d%02d", t.time.Year(), t.time.Month(), t.time.Day())
}
