package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/willfantom/overclirr/cmd"
	"github.com/willfantom/overclirr/ui"
)

const (
	initialLogLevel string = "panic"
)

func main() {

	if err := cmd.RootCmd.Execute(); err != nil {
		if logrus.GetLevel() < logrus.WarnLevel {
			ui.Error("Try setting the log level higher (e.g. info) to see what is going on!")
		}
		logrus.WithField("extended", err.Error()).
			Fatalln("an error occurred executing the command")
	}
}

// init adds the config information to the global viper
func init() {
	//set initial log level
	lvl, _ := logrus.ParseLevel(initialLogLevel)
	logrus.SetLevel(lvl)

	// define configuration file info
	viper.SetConfigName("overclirr")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath(".")

	// set default configuration values
	viper.SetDefault("showLoadingSpinner", true)
	viper.SetDefault("log", initialLogLevel)
	viper.SetDefault("profiles", nil)

	// create config file
	if err := viper.SafeWriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
			ui.Fatal("Configuration file could not be created!", err)
		}
	}

	// read in existing config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			ui.ColorPrint("Try removing the file and readding the configuration values\nOr check the wiki on GitHub...", ui.Blue)
			ui.Fatal("Existing configuration file could not be read in", err)
		} else {
			ui.ColorPrint("Config file should have been created automatically...", ui.White)
			ui.Fatal("No configuration file could be found", err)
		}
	}

	// config cascade
	if viper.GetString("log") != initialLogLevel {
		viper.Set("showLoadingSpinner", false)
	}

	logrus.Traceln("configuration init success")
}
