package cmd

import (
	"github.com/spf13/cobra"
	"github.com/willfantom/overclirr/ui"
	"github.com/willfantom/overclirr/utility"
)

var addProfileCmd = &cobra.Command{
	Use:     "add-profile",
	Aliases: []string{"add"},
	Short:   "Add a new Overseerr login profile",
	Long:    "Create a new Overseerr login profile and write it to the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		ui.ColorPrintBold("Create New Profile\n", ui.Blue)
		ui.ColorPrint(" | profile: "+overseerrProfileName+"\n", ui.White)
		var profile utility.OverseerrProfile
		profile.URL = ui.GetInput("URL of the Overseerr instance:", utility.URLValidator)
		profile.Locale = ui.GetInput("Local of the Overseerr instance (e.g. en):", utility.LocaleValidator)
		profile.Auth.Type = utility.SelectProfileAuthType()
		switch profile.Auth.Type {
		case utility.OverseerrAuthTypeKey:
			profile.Auth.Key = ui.GetMaskedInput("API key for the Overseerr instance:", utility.NonEmptyValidator)
		case utility.OverseerrAuthTypeLocal:
			profile.Auth.Email = ui.GetInput("Email address for the user account:", utility.EmailValidator)
			profile.Auth.Password = ui.GetMaskedInput("Password for the user account:", nil)
		case utility.OverseerrAuthTypePlex:
			profile.Auth.PlexToken = ui.GetMaskedInput("Plex Token for the user account:", utility.NonEmptyValidator)
		}
		if _, err := profile.Connect(); err != nil {
			ui.Fatal("Could not connect to Overseer with given profile information", err)
		}
		if err := utility.WriteOverseerrProfile(overseerrProfileName, profile, true); err != nil {
			ui.Fatal("Valid profile, but can not be written to config", err)
		}
		ui.Success("Profile Added")
	},
}

var delProfileCmd = &cobra.Command{
	Use:     "delete-profile",
	Aliases: []string{"del-profile", "del", "delete", "remove"},
	Short:   "Remove an Overseerr login profile",
	Long:    "Remove an Overseerr login profile from OverCLIrr's configuration",
	Run: func(cmd *cobra.Command, args []string) {
		ui.ColorPrintBold("Attempting to remove profile...\n", ui.Blue)
		ui.ColorPrint(" | profile: "+overseerrProfileName+"\n", ui.White)
		ui.DestructiveConfirmation()
		if err := utility.DeleteOverseerrProfile(overseerrProfileName); err != nil {
			ui.Fatal("Failed to remove profile", err)
		}
		ui.Success("Profile Added")
	},
}

var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "View and test installed login profiles",
	Run: func(cmd *cobra.Command, args []string) {
		allProfiles := utility.GetAllOverseerrProfiles()
		for name, profile := range allProfiles {
			if _, err := profile.Connect(); err != nil {
				ui.Error("Could not connect using profile: " + name)
			} else {
				ui.Success("Connection made using profile: " + name)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(addProfileCmd)
	RootCmd.AddCommand(delProfileCmd)
	RootCmd.AddCommand(profilesCmd)
}
