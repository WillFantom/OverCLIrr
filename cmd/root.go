package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/willfantom/goverseerr"
	"github.com/willfantom/overclirr/ui"
	"github.com/willfantom/overclirr/utility"
)

const (
	defaultProfile string = "default"
)

// persistent flags
var (
	logLevel             string
	overseerrProfileName string
	noTitle              bool
)

// instance to use
var overseerr *goverseerr.Overseerr

func setOverseer(profileName string) {
	ui.StartSpinner()
	profile, err := utility.GetOverseerrProfile(profileName)
	if err != nil {
		ui.Fatal("Overseerr profile does not exist", err)
	}
	instance, err := profile.Connect()
	if err != nil {
		ui.Fatal("Could not connect using overseerr profile: "+profileName, err)
	}
	overseerr = instance
	ui.StopSpinner()
}

var RootCmd = &cobra.Command{
	Use:     "overclirr",
	Aliases: []string{"ocrr", "overseerr", "overseerr-cli"},
	Short:   "Manage media servers from the command line",
	Long:    `A simple command line tool for managing media server(s) with Overseerr!`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupLogger()
		if !noTitle {
			ui.PrintTitleBox("OverCLIrr", "An Overseerr Management Tool")
		}
		logrus.WithFields(logrus.Fields{
			"command": cmd.Name(),
			"args":    args,
		}).Debugln("running command")
	},
	Run: func(cmd *cobra.Command, args []string) {
		allProfiles := utility.GetAllOverseerrProfiles()
		ui.ColorPrintBold("Profiles Found in Configuration: ", ui.Blue)
		ui.ColorPrint(fmt.Sprintf("%d\n\n", len(allProfiles)), ui.White)
		cmd.Help()
	},
}

func setupLogger() {
	if level, err := logrus.ParseLevel(logLevel); err != nil {
		ui.Error("Invalid log level given: " + logLevel)
		logrus.SetLevel(logrus.PanicLevel)
	} else {
		logrus.SetLevel(level)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVar(&noTitle, "no-title", false, "stop printing the big fu*king title")
	RootCmd.PersistentFlags().StringVar(&logLevel, "log", "panic", "set the log level (fatal, error, info, debug, trace)")
	RootCmd.PersistentFlags().StringVar(&overseerrProfileName, "profile", defaultProfile, "use a specific overseerr login profile name")
	viper.BindPFlag("log", RootCmd.PersistentFlags().Lookup("log"))
}
