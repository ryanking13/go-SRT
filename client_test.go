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

func tomorrow() string {
	t := time.Now().AddDate(0, 0, 1)
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func TestSRT(t *testing.T) {

	client := srt.New()
	// client.SetDebug()

	username := os.Getenv("SRT_USERNAME")
	password := os.Getenv("SRT_PASSWORD")

	t.Run("Username/Password Check", func(t *testing.T) {
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

	t.Run("Search/Reserve/Cancel Test", func(t *testing.T) {
		defer func() {
			t.Log("Clean Up")
			reservations, _ := client.Reservations()
			for _, r := range reservations {
				client.Cancel(r)
			}
		}()

		searchParams := &srt.SearchParams{
			Dep:  "수서",
			Arr:  "부산",
			Date: tomorrow(),
		}
		trains, err := client.SearchTrain(searchParams)
		if err != nil {
			t.Errorf("SRT SearchTrain Failed: %s", err.Error())
			return
		}
		t.Log(trains)

		reserveParams := &srt.ReserveParams{
			Train:      trains[len(trains)-1],
			Passengers: []*srt.Passenger{srt.Adult(2), srt.Child(1)},
		}

		reservation, err := client.Reserve(reserveParams)
		if err != nil {
			t.Errorf("SRT Reserve Failed: %s", err.Error())
			return
		}
		t.Log(reservation)

		reservations, err := client.Reservations()
		if err != nil {
			t.Errorf("SRT Reservations Failed: %s", err.Error())
			return
		}
		t.Log(reservations)

		tickets, err := client.Tickets(reservations[0])
		if err != nil {
			t.Errorf("SRT Tickets Failed: %s", err.Error())
			return
		}
		t.Log(tickets)

		err = client.Cancel(reservations[0])
		if err != nil {
			t.Errorf("SRT Cancel Failed: %s", err.Error())
			return
		}
	})

	t.Run("Logout Test", func(t *testing.T) {
		err := client.Logout()
		if err != nil {
			t.Errorf("SRT Logout Failed: %s", err.Error())
		}
	})
}
