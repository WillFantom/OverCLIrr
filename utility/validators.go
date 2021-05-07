package utility

import (
	"errors"
	"net/url"
	"regexp"

	"golang.org/x/text/language"
)

func URLValidator(u string) error {
	if len(u) == 0 {
		return errors.New("url must not be empty")
	}
	if parsedURL, err := url.Parse(u); err != nil {
		return errors.New("not parseable as a url")
	} else {
		if !parsedURL.IsAbs() {
			return errors.New("url must be absolute (have a scheme such as https://)")
		}
	}
	return nil
}

func LocaleValidator(l string) error {
	if _, err := language.Parse(l); err != nil {
		return errors.New("locale not known")
	}
	return nil
}

func NonEmptyValidator(in string) error {
	if len(in) == 0 {
		return errors.New("cannot be empty")
	}
	return nil
}

func EmailValidator(in string) error {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(in) {
		return errors.New("must be a valid email address")
	}
	return nil
}
