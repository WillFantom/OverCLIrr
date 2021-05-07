package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/willfantom/goverseerr"
	"github.com/willfantom/overclirr/ui"
	"github.com/willfantom/overclirr/utility"
)

// Flags

var (
	searchPage int
)

// Commands

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for and request new media",
	PreRun: func(cmd *cobra.Command, args []string) {
		setOverseer(overseerrProfileName)
	},
	Run: func(cmd *cobra.Command, args []string) {
		ui.ColorPrintBold("Search Media\n", ui.White)
		searchTerm := ui.GetInput("Enter the search term", utility.NonEmptyValidator)
		ui.StartSpinner()
		results, err := overseerr.Search(searchTerm, searchPage)
		if err != nil {
			ui.Fatal("Could not fetch search results", err)
		}
		ui.StopSpinner()
		ui.ColorPrint(fmt.Sprintf("Found %d results on page %d of %d\n", len(results.Results), results.Page, results.TotalPages), ui.White)
		ui.SelectorTemplated("Select a request:", results, resultSelectionTemplate)
	},
}

func init() {
	searchCmd.Flags().IntVarP(&searchPage, "page", "p", 1, "Which page of results to view")
	RootCmd.AddCommand(searchCmd)
}

// Result Selection

var resultSelectionTemplate = &promptui.SelectTemplates{
	Label:    "{{ . }}:",
	Active:   "🍿  {{ if .Title }} {{ .Title | magenta }}  {{else}} {{ .Name | magenta }} {{end}} {{if .ReleaseDate }}{{ (.ReleaseDate) | white }}{{end}}{{if .FirstAiredDate }}{{ (.FirstAiredDate) | white }}{{end}} {{ .MediaType.ToEmoji }}",
	Inactive: "   {{ if .Title }} {{ .Title | cyan }}  {{else}} {{ .Name | magenta }} {{end}} {{if .ReleaseDate }}{{ (.ReleaseDate) | red }}{{end}}{{if .FirstAiredDate }}{{ (.FirstAiredDate) | red }}{{end}} {{ .MediaType.ToEmoji }}",
	Selected: "🍿  {{ if .Title }} {{ .Title | cyan }}  {{else}} {{ .Name | magenta }} {{end}} {{ .MediaTypeEmoji }}",
	Details: `
|- Search Result Details -|
{{ "Title:" | faint }}	{{ if .Title }} {{ .Title }}  {{else}} {{ .Name }} {{end}}
{{ "Media Type:" | faint }}	{{ .MediaType.ToEmoji }} {{ .MediaType }}
{{ "Media Status:" | faint }}	{{ .MediaInfo.Status.ToEmoji }} {{ .MediaInfo.Status.ToString }}
{{ if .ReleaseDate }}{{"Release Date:" | faint }}	{{ .ReleaseDate }}{{end}}{{ if .FirstAiredDate }}{{ "Release Date:" | faint }}	{{ .FirstAiredDate }}{{end}}`,
}

// Result Actions

var resultActionsPromptTemplate = &promptui.SelectTemplates{
	Label:    "{{ . }}:",
	Active:   "➡️  {{ .Name | magenta }}",
	Inactive: "   {{ .Name | cyan }}",
	Selected: "➡️  {{ .Name | magenta }}",
	Details: `
|- Action Details
{{ "Name:" | faint }}	{{ .Name }}
{{ "Description:" | faint }}	{{ .Description }}`,
}

type ResultAction struct {
	Name        string
	Description string
	Handler     func(r *goverseerr.GenericSearchResult)
	Validator   func(r *goverseerr.GenericSearchResult) bool
}

var ResultActions = []ResultAction{
	{
		Name:        "Request",
		Description: "Request this media!",
		Validator: func(r *goverseerr.GenericSearchResult) bool {
			if r.MediaType == goverseerr.MediaTypePerson {
				ui.Error("Can not request a person!")
				return false
			}
			if r.MediaInfo.Status != goverseerr.MediaStatusUnknown {
				ui.Error("Media is already requested or available!")
				return false
			}
			return true
		},
		Handler: func(r *goverseerr.GenericSearchResult) {

		},
	},
	{
		Name:        "Details",
		Description: "See more details about this item",
		Validator: func(r *goverseerr.GenericSearchResult) bool {
			return true
		},
		Handler: func(r *goverseerr.GenericSearchResult) {
		},
	},
}
