package ui

import (
	"strings"

	"github.com/fatih/color"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
)

func stopUIElements() {
	StopSpinner()
}

const (
	Red     color.Attribute = color.FgRed
	Blue    color.Attribute = color.FgBlue
	Magenta color.Attribute = color.FgMagenta
	White   color.Attribute = color.FgWhite
)

func ColorPrint(message string, fgColor color.Attribute) {
	color.New(fgColor).Printf("%s", message)
}

func ColorPrintBold(message string, fgColor color.Attribute) {
	color.New(fgColor, color.Bold).Printf("%s", message)
}

func Table(values [][]string) {
	pterm.DefaultTable.WithHasHeader().WithData(values).Render()
}

func Success(message string) {
	stopUIElements()
	color.New(color.FgBlack, color.BgGreen).Printf(" SUCCESS ")
	color.New(color.FgGreen).Printf(" %s\n", message)
}

func Error(message string) {
	stopUIElements()
	color.New(color.FgBlack, color.BgRed, color.Bold).Printf(" ERROR ")
	color.New(color.FgRed).Printf(" %s\n", message)
}

func ErrorSub(message, sub string) {
	stopUIElements()
	color.New(color.FgBlack, color.BgRed, color.Bold).Printf(" ERROR ")
	color.New(color.FgRed).Printf(" %s\n", message)
	color.New(color.FgWhite, color.Bold).Printf(" > %s", sub)
}

func FatalQuiet(message string, err error) {
	stopUIElements()
	logrus.WithField("extended", err.Error()).Fatalln(strings.ToLower(message))
}

func Fatal(message string, err error) {
	stopUIElements()
	color.New(color.FgBlack, color.BgRed, color.Bold).Printf(" FATAL ")
	color.New(color.FgRed).Printf(" %s\n", message)
	logrus.WithField("extended", err.Error()).Fatalln(strings.ToLower(message))
}
