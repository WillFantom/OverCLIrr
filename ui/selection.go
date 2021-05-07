package ui

import (
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
)

func Selector(title string, options []string) (int, string) {
	prompt := promptui.Select{
		Label: title,
		Items: options,
		Size:  5,
	}
	logrus.WithField("title", title).Traceln("running selector menu")
	index, option, err := prompt.Run()
	if err != nil {
		logrus.WithField("title", title).Errorln("selector menu failed")
		FatalQuiet("Selector Failed", err)
	}
	logrus.WithField("selection", option).Traceln("ending selector menu")
	return index, option
}

func SelectorTemplated(title string, options interface{}, template *promptui.SelectTemplates) (int, string) {
	prompt := promptui.Select{
		Label:     title,
		Items:     options,
		Size:      5,
		Templates: template,
	}
	logrus.WithField("title", title).Traceln("running selector menu")
	index, option, err := prompt.Run()
	if err != nil {
		logrus.WithField("title", title).Errorln("selector menu failed")
		FatalQuiet("Selector Failed", err)
	}
	logrus.WithField("selection", option).Traceln("ending selector menu")
	return index, option
}
