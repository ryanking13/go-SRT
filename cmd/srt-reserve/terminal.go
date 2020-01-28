package main

import (
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
)

func hi() {
	log.Println(aurora.Cyan(`
░██████╗██████╗░████████╗
██╔════╝██╔══██╗╚══██╔══╝
╚█████╗░██████╔╝░░░██║░░░
░╚═══██╗██╔══██╗░░░██║░░░
██████╔╝██║░░██║░░░██║░░░
╚═════╝░╚═╝░░╚═╝░░░╚═╝░░░
`))
}

func setColors() {
	// TODO: linux compatible
	log.SetFlags(0)
	log.SetOutput(colorable.NewColorableStdout())
}

func dieOnInterrupt(err error) error {
	if err == terminal.InterruptErr {
		os.Exit(0)
	}
	return err
}
