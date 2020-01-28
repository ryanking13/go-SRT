package main

import (
	"errors"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
	srt "github.com/ryanking13/go-SRT"
)

func hi() {
	// log.Println(aurora.Red("***************************"))
	// log.Println("************SRT************")
	// log.Println("***************************")
}

func setColors() {
	// TODO: linux compatible
	log.SetFlags(0)
	log.SetOutput(colorable.NewColorableStdout())
}

func required(msg string) func(val interface{}) error {
	return func(val interface{}) error {
		err := survey.Required(val)
		if err != nil {
			return errors.New(msg)
		}
		return nil
	}
}

func login() (*srt.Client, error) {
	log.Println(aurora.Cyan("***** SRT 서버에 로그인 *****"))

	username := ""
	password := ""

	survey.AskOne(
		&survey.Input{
			Message: "ID:",
		},
		&username,
		survey.WithValidator(required("ID를 입력하세요")),
	)

	survey.AskOne(
		&survey.Password{
			Message: "Password:",
		},
		&password,
		survey.WithValidator(required("비밀번호를 입력하세요")),
	)

	return srt.New(), nil
}

func main() {
	setColors()
	hi()

	_, err := login()

	if err != nil {
		log.Println(aurora.Red("Login Error"))
	}
}
