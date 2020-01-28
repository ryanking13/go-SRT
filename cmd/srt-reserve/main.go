package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/logrusorgru/aurora"
	srt "github.com/ryanking13/go-SRT"
)

func required(msg string) func(val interface{}) error {
	return func(val interface{}) error {
		err := survey.Required(val)
		if err != nil {
			return errors.New(msg)
		}
		return nil
	}
}

func isPositive(val interface{}) error {
	if str, ok := val.(string); ok {
		num, err := strconv.Atoi(str)
		if err != nil {
			return err
		}

		if num < 1 {
			return errors.New("승객 수는 1 이상이여야 합니다.")
		}
	}
	return nil
}

func login() (*srt.Client, error) {
	username := ""
	password := ""

	client := srt.New()

	for {
		err := survey.AskOne(
			&survey.Input{
				Message: "SRT ID:",
			},
			&username,
			survey.WithValidator(required("ID를 입력하세요")),
		)

		if err != nil {
			return nil, err
		}

		err = survey.AskOne(
			&survey.Password{
				Message: "SRT Password:",
			},
			&password,
			survey.WithValidator(required("비밀번호를 입력하세요")),
		)

		if err != nil {
			return nil, err
		}

		err = client.Login(username, password)

		if err != nil {
			log.Println(err.Error())
			continue
		}

		break
	}

	return client, nil
}

func selectTrain(client *srt.Client) (*srt.SearchParams, error) {

	stationDep := ""
	stationArr := ""
	dateDep := ""

	stations := []string{
		"수서",
		"동탄",
		"지제",
		"천안아산",
		"오송",
		"대전",
		"공주",
		"익산",
		"정읍",
		"광주송정",
		"나주",
		"목포",
		"김천구미",
		"동대구",
		"신경주",
		"울산(통도사)",
		"부산",
	}

	// Step 1) Select Departure

	err := survey.AskOne(
		&survey.Select{
			Message: "출발역:",
			Options: stations,
			Default: stations[0],
		},
		&stationDep,
	)

	if err != nil {
		return nil, err
	}

	// Step 2) Select Arrival

	err = survey.AskOne(
		&survey.Select{
			Message: "도착역:",
			Options: stations,
			Default: stations[0],
		},
		&stationArr,
	)

	if err != nil {
		return nil, err
	}

	// Step 3) Select Date

	nextCnt := 0
	numDays := 10
	toPrev := "이전 날짜로"
	toNext := "다음 날짜로"
	for {
		isNext := nextCnt > 0
		days := make([]string, 0)

		if isNext {
			days = append(days, toPrev)
		}

		date := today().NextDay(nextCnt * numDays)
		for i := 0; i < numDays; i++ {
			days = append(days, date.String())
			date = date.NextDay(1)
		}
		days = append(days, toNext)

		err = survey.AskOne(
			&survey.Select{
				Message: "날짜:",
				Options: days,
				Default: days[0],
			},
			&dateDep,
		)

		if err != nil {
			return nil, err
		}

		if dateDep == toPrev {
			nextCnt--
			continue
		} else if dateDep == toNext {
			nextCnt++
			continue
		}
		break
	}

	if err != nil {
		return nil, err
	}

	return &srt.SearchParams{
		Dep:  stationDep,
		Arr:  stationArr,
		Date: dateDep,
	}, nil
}

func searchTrain(client *srt.Client, params *srt.SearchParams) (*srt.ReserveParams, error) {

	passengers := make([]*srt.Passenger, 0)
	passengersCnt := 0

	// Select Passengers
	err := survey.AskOne(
		&survey.Input{
			Message: "인원 수:",
		},
		&passengersCnt,
		survey.WithValidator(
			survey.ComposeValidators(
				required("인원 수를 입력하세요"),
				isPositive,
			),
		),
	)

	passengers = append(passengers, srt.Adult(passengersCnt))

	// Search Train
	trains, err := client.SearchTrain(params)
	if err != nil {
		return nil, err
	}

	trainsStr := make([]string, len(trains))
	for i, train := range trains {
		trainsStr[i] = train.String()
	}

	selectedIndex := 0
	err = survey.AskOne(
		&survey.Select{
			Message: "열차 선택:",
			Options: trainsStr,
			Default: "부산",
		},
		&selectedIndex,
	)

	if err != nil {
		return nil, err
	}

	selected := trains[selectedIndex]

	return &srt.ReserveParams{
		Train:      selected,
		Passengers: passengers,
	}, nil
}

func reserve(client *srt.Client, params *srt.ReserveParams) (*srt.Reservation, error) {
	result, err := client.Reserve(params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func main() {
	setColors()
	hi()

	client, err := login()

	if dieOnInterrupt(err) != nil {
		log.Println(aurora.Red(fmt.Sprintf("Login Error: %v", err)))
		os.Exit(1)
	}

	searchParams, err := selectTrain(client)

	if dieOnInterrupt(err) != nil {
		log.Println(aurora.Red(fmt.Sprintf("Select Error: %v", err)))
	}

	reserveParams, err := searchTrain(client, searchParams)

	if dieOnInterrupt(err) != nil {
		log.Println(aurora.Red(fmt.Sprintf("Search Error: %v", err)))
	}

	reserveResult, err := reserve(client, reserveParams)

	if dieOnInterrupt(err) != nil {
		log.Println(aurora.Red(fmt.Sprintf("Reserve Error: %v", err)))
	}

	log.Println(aurora.Green("예약 완료! 홈페이지에서 결제를 완료하세요."))
	log.Println(reserveResult)
}
