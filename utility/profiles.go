package utility

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/willfantom/goverseerr"
	"github.com/willfantom/overclirr/ui"
)

type OverseerrAuthType string

const (
	OverseerrAuthTypeKey   OverseerrAuthType = "key"
	OverseerrAuthTypeLocal OverseerrAuthType = "local"
	OverseerrAuthTypePlex  OverseerrAuthType = "plex"
)

// OverseerrAuthConfig contains the authorization details for an account
type OverseerrAuthConfig struct {
	Type      OverseerrAuthType `json:"type"`
	Key       string            `json:"key,omitempty"`
	Email     string            `json:"email,omitempty"`
	Password  string            `json:"password,omitempty"`
	PlexToken string            `json:"plexToken,omitempty"`
}

// OverseerrProfile contains all the details to login to an account on Overseerr
type OverseerrProfile struct {
	URL           string              `json:"url"`
	CustomHeaders map[string]string   `json:"customHeaders"`
	Locale        string              `json:"locale" default:"en"`
	Auth          OverseerrAuthConfig `json:"auth"`
}

// QuickValidate ensures both the URL and Locale are valid for a profile
func (c OverseerrProfile) QuickValidate() error {
	if err := URLValidator(c.URL); err != nil {
		return err
	}
	if err := LocaleValidator(c.Locale); err != nil {
		return err
	}
	return nil
}

// Connect creates a logged in Overseerr instance from a profile
func (p OverseerrProfile) Connect() (*goverseerr.Overseerr, error) {
	switch p.Auth.Type {
	case OverseerrAuthTypeKey:
		logrus.Debugln("connecting to overseerr using a key based profile")
		if o, err := goverseerr.NewKeyAuth(p.URL, p.CustomHeaders, p.Locale, p.Auth.Key); err == nil {
			return o, nil
		} else {
			return nil, err
		}
	case OverseerrAuthTypeLocal:
		logrus.Debugln("connecting to overseerr using a email/password based profile")
		if o, err := goverseerr.NewLocalAuth(p.URL, p.CustomHeaders, p.Locale, p.Auth.Email, p.Auth.Password); err == nil {
			return o, nil
		} else {
			return nil, err
		}
	case OverseerrAuthTypePlex:
		logrus.Debugln("connecting to overseerr using a plex token based profile")
		if o, err := goverseerr.NewPlexAuth(p.URL, p.CustomHeaders, p.Locale, p.Auth.PlexToken); err == nil {
			return o, nil
		} else {
			return nil, err
		}
	default:
		return nil, errors.New("profile has an invalid authorization type: " + string(p.Auth.Type))
	}
}

// GetAllOverseerrProfiles returns all configured overseerr profiles
// and panics if fails
func GetAllOverseerrProfiles() map[string]OverseerrProfile {
	profiles := make(map[string]OverseerrProfile)
	if err := viper.UnmarshalKey("profiles", &profiles); err != nil {
		ui.Fatal("Failed to parse Overseerr profiles from the file", err)
	}
	return profiles
}

// GetOverseerrProfile returns a configured overseerr profile if it exists.
// If not, an error is returned
func GetOverseerrProfile(profileName string) (OverseerrProfile, error) {
	profiles := GetAllOverseerrProfiles()
	if profile, ok := profiles[profileName]; ok {
		return profile, nil
	}
	return OverseerrProfile{}, errors.New("overseerr profile not found")
}

// WriteOverseerrProfile writes a profile to the configuration
func WriteOverseerrProfile(profileName string, profile OverseerrProfile, overwrite bool) error {
	if err := profile.QuickValidate(); err != nil {
		return err
	}
	profiles := GetAllOverseerrProfiles()
	if _, ok := profiles[profileName]; ok && !overwrite {
		return errors.New("overseer profile with that name already exists")
	}
	profiles[profileName] = profile
	viper.Set("profiles", profiles)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

// WriteOverseerrProfile writes a profile to the configuration
func DeleteOverseerrProfile(profileName string) error {
	profiles := GetAllOverseerrProfiles()
	if _, ok := profiles[profileName]; !ok {
		return errors.New("profile with name '" + profileName + "' does not exist")
	}
	delete(profiles, profileName)
	viper.Set("profiles", profiles)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

func SelectProfileAuthType() OverseerrAuthType {
	options := []string{"API Key [admins]", "Plex Token", "Email/Password"}
	idx, _ := ui.Selector("What auth type should be used?", options)
	switch idx {
	case 0:
		return OverseerrAuthTypeKey
	case 1:
		return OverseerrAuthTypePlex
	case 2:
		return OverseerrAuthTypeLocal
	}
	return OverseerrAuthTypeKey
}
