package srt_test

import (
	"testing"

	srt "github.com/ryanking13/go-SRT"
)

func TestSRT(t *testing.T) {

	client := srt.New()
	
	t.Run("Login Test", func(t *testing.T) {
		err := client.Login()
		if err != nil {
			t.Errorf("SRT Login Failed: %s", err.Error())
		}
	})

	t.Run("SearchTrain Test", func(t *testing.T) {
		err := client.SearchTrain()
		if err != nil {
			t.Errorf("SRT SearchTrain Failed: %s", err.Error())
		}
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
