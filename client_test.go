package srt_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	srt "github.com/ryanking13/go-SRT"
)

func today() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func TestSRT(t *testing.T) {

	client := srt.New()
	// client.SetDebug()

	username := os.Getenv("SRT_USERNAME")
	password := os.Getenv("SRT_PASSWORD")

	t.Run("Username, Password Check", func(t *testing.T) {
		if username == "" || password == "" {
			t.Error("SRT_USERNAME or SRT_PASSWORD not set")
		}
	})

	t.Run("Login Success Test", func(t *testing.T) {
		err := client.Login(username, password)
		if err != nil {
			t.Errorf("SRT Login Failed: %s", err.Error())
		}
	})

	t.Run("Login Fail Test (ID)", func(t *testing.T) {
		err := client.Login("deadbeef", password)
		if err == nil {
			t.Error("Invalid Login Credential Bypassed (ID)")
		}
	})

	t.Run("Login Fail Test (PW)", func(t *testing.T) {
		err := client.Login(username, "deadbeef")
		if err == nil {
			t.Error("Invalid Login Credential Bypassed (PW)")
		}
	})

	t.Run("SearchTrain Test", func(t *testing.T) {
		trains, err := client.SearchTrainAll("수서", "부산", today(), "000000")
		if err != nil {
			t.Errorf("SRT SearchTrain Failed: %s", err.Error())
		}
		t.Log(trains)
	})

	t.Run("Reserve Test", func(t *testing.T) {
		err := client.Reserve()
		if err != nil {
			t.Errorf("SRT Reserve Failed: %s", err.Error())
		}
	})

	t.Run("Reservations Test", func(t *testing.T) {
		err := client.Reservations()
		if err != nil {
			t.Errorf("SRT Reservations Failed: %s", err.Error())
		}
	})

	t.Run("TicketInfo Test", func(t *testing.T) {
		err := client.TicketInfo()
		if err != nil {
			t.Errorf("SRT TicketInfo Failed: %s", err.Error())
		}
	})

	t.Run("Cancel Test", func(t *testing.T) {
		err := client.Cancel()
		if err != nil {
			t.Errorf("SRT Cancel Failed: %s", err.Error())
		}
	})

	t.Run("Logout Test", func(t *testing.T) {
		err := client.Logout()
		if err != nil {
			t.Errorf("SRT Logout Failed: %s", err.Error())
		}
	})
}
