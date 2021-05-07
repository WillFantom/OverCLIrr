package cmd

import (
	"context"

	"github.com/google/go-github/v35/github"
	semver "github.com/hashicorp/go-version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/willfantom/overclirr/ui"
)

const (
	defaultVersion string = "no-version"
	ghUser         string = "willfantom"
	ghRepo         string = "goverseerr"
)

var version string = defaultVersion

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of OverCLIrr",
	Run: func(cmd *cobra.Command, args []string) {
		ui.ColorPrint("OverCLIrr Version: ", ui.White)
		ui.ColorPrintBold(version+"\n", ui.Magenta)
	},
}

var overseerrVersion = &cobra.Command{
	Use:   "overseerr-version",
	Short: "Print the version of the connected Overseerr instance",
	PreRun: func(cmd *cobra.Command, args []string) {
		setOverseer(overseerrProfileName)
	},
	Run: func(cmd *cobra.Command, args []string) {
		status, err := overseerr.Status()
		if err != nil {
			ui.Fatal("Could not determine Overseerr version", err)
		}
		ui.ColorPrint("Overseerr Version: ", ui.White)
		ui.ColorPrintBold(status.Version+"\n", ui.Magenta)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check if OverCLIrr is up to date",
	Run: func(cmd *cobra.Command, args []string) {
		available, err := checkForUpdate()
		if available {
			ui.ColorPrint("An update is available\nCheck the releases on GitHub\n", ui.Blue)
		} else if !available && err == nil {
			ui.ColorPrint("You are running the latest version according to GitHub\n", ui.Blue)
		} else {
			ui.Error("Could not perform version check")
			ui.Fatal("Perhaps this is unversioned? Or you can't connect to GitHub?", err)
		}
	},
}

func checkForUpdate() (bool, error) {
	//semver local version
	semverVersion, err := semver.NewSemver(version)
	if err != nil {
		logrus.WithField("extended", err.Error()).
			Errorln("could not determine the version of overclirr")
		return false, err
	}

	//get latest github release tag
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), ghUser, ghRepo)
	if err != nil {
		logrus.WithField("extended", err.Error()).
			Errorln("could get repository release info from github")
		return false, err
	}
	ghVer, err := semver.NewSemver(*release.TagName)
	if err != nil {
		logrus.WithField("extended", err.Error()).
			Errorln("latest github release tag not semver compliant")
		return false, err
	}

	//compare
	if ghVer.GreaterThan(semverVersion) {
		logrus.Infoln("found a more recent release on github")
		return true, nil
	}
	logrus.Infoln("overclirr found to be latest version")
	return false, nil
}

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(updateCmd)
	RootCmd.AddCommand(overseerrVersion)
}
