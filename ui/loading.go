package ui

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/viper"
	"github.com/willfantom/reticulating-go"
)

var (
	loadingSpinner     *spinner.Spinner
	loadingMessageFunc func() string
)

func init() {
	loadingSpinner = spinner.New(spinner.CharSets[2], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	loadingSpinner.HideCursor = true
	loadingSpinner.Color("white")
	loadingMessageFunc = reticulating.GetLoadingMessage
}

func StartSpinner() {
	if viper.GetBool("showLoadingSpinner") {
		loadingSpinner.Start()
		go updateLoadingMessage()
	}
}

func updateLoadingMessage() {
	for loadingSpinner.Active() {
		loadingSpinner.Suffix = "  " + loadingMessageFunc()
		time.Sleep(500 * time.Millisecond)
	}
}

func StopSpinner() {
	loadingSpinner.Stop()
}
