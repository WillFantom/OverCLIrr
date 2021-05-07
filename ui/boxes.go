package ui

import (
	box "github.com/Delta456/box-cli-maker/v2"
	"github.com/pterm/pterm"
)

var (
	titleBox = box.New(box.Config{
		Type:         "Double",
		Color:        "HiMagenta",
		ContentAlign: "Center",
		TitlePos:     "Inside",
		Px:           3,
		Py:           0,
	})

	contentsBox = box.New(box.Config{
		Type:         "Round",
		Color:        "White",
		ContentAlign: "Center",
		TitlePos:     "Top",
		Px:           2,
		Py:           1,
	})

	compactBox = box.New(box.Config{
		Type:         "Round",
		Color:        "White",
		ContentAlign: "Center",
		TitlePos:     "Inside",
		Px:           2,
		Py:           0,
	})
)

const maxLineLength int = 60

func stringToParagraph(in string) string {
	return pterm.DefaultParagraph.WithMaxWidth(maxLineLength).Sprint(in)
}

func PrintTitleBox(title, sub string) {
	stopUIElements()
	titleBox.Println(title, stringToParagraph(sub))
}

func PrintBox(title, sub string) {
	stopUIElements()
	contentsBox.Println(title, stringToParagraph(sub))
}

func PrintCompactBox(title, sub string) {
	stopUIElements()
	compactBox.Println(title, stringToParagraph(sub))
}
