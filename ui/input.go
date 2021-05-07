package ui

import (
	"strings"

	"github.com/manifoldco/promptui"
)

func GetInput(title string, validator func(string) error) string {
	prompt := promptui.Prompt{
		Label:    title,
		Validate: validator,
	}
	result, err := prompt.Run()
	if err != nil {
		Fatal("User input failed", err)
	}
	return result
}

func GetMaskedInput(title string, validator func(string) error) string {
	prompt := promptui.Prompt{
		Label:    title,
		Validate: validator,
		Mask:     '*',
	}
	result, err := prompt.Run()
	if err != nil {
		Fatal("User input failed", err)
	}
	return result
}

func DestructiveConfirmation() {
	prompt := promptui.Prompt{
		Label:     "What you are about to do could be destructive, continue?",
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil {
		Fatal("destructive confirmation failed", err)
	}
	if strings.ToLower(result) != "y" {
		Error("Aborted!")
		FatalQuiet("destructive confirmation rejected", err)
	}
}
