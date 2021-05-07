package cmd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/willfantom/goverseerr"
	"github.com/willfantom/overclirr/ui"
)

// Flags

var (
	reqPage     int
	reqPageSize int
	allRequests bool
)

// Commands

var requestsCmd = &cobra.Command{
	Use:     "requests",
	Aliases: []string{"my-requests"},
	Short:   "View the requests made by the current profile user",
	PreRun: func(cmd *cobra.Command, args []string) {
		setOverseer(overseerrProfileName)
	},
	Run: func(cmd *cobra.Command, args []string) {
		ui.ColorPrintBold("Overseerr Requests\n", ui.White)
		ui.StartSpinner()
		me, err := overseerr.GetLoggedInUser()
		if err != nil {
			ui.Fatal("Could not get the deatails for the current profile user", err)
		}
		var requests []*goverseerr.MediaRequest
		var pageInfo *goverseerr.Page
		if !allRequests {
			r, p, err := overseerr.GetUserRequests(me.ID, reqPage-1, reqPageSize)
			if err != nil {
				ui.Fatal("Could not get the request list for the current user", err)
			}
			requests = r
			pageInfo = p
		} else {
			r, p, err := overseerr.GetRequests(reqPage-1, reqPageSize, goverseerr.RequestFileterAll, goverseerr.RequestSortAdded)
			if err != nil {
				ui.Fatal("Could not get the full request list", err)
			}
			requests = r
			pageInfo = p
		}
		ui.StopSpinner()
		ui.ColorPrint(fmt.Sprintf("Found %d results on page %d of %d\n", len(requests), pageInfo.Page, pageInfo.Pages), ui.White)
		request := selectRequest(requests)
		actIdx, _ := ui.SelectorTemplated("Select an Action", RequestActions, requestActionsPromptTemplate)
		if !RequestActions[actIdx].Validator(request) {
			ui.Fatal("Action not compatible with this request!", errors.New("action not compatible with request"))
		}
		RequestActions[actIdx].Handler(request)
	},
}

func init() {
	requestsCmd.Flags().IntVarP(&reqPage, "page", "p", 1, "Which page of requests to view")
	requestsCmd.Flags().IntVar(&reqPageSize, "page-size", 20, "How many requests to show per page")
	requestsCmd.Flags().BoolVarP(&allRequests, "all", "a", false, "Attempt to get all requests, not just cirrent user's")
	RootCmd.AddCommand(requestsCmd)
}

// Request Selection

var requestSelectionTemplate = &promptui.SelectTemplates{
	Label:    "{{ . }}:",
	Active:   "🍿  {{ .ContentTitle | magenta }} ({{ .ContentDate | white }}) {{ .MediaTypeEmoji }}",
	Inactive: "   {{ .ContentTitle | cyan }} ({{ .ContentDate | red }}) {{ .MediaTypeEmoji }}",
	Selected: "🍿  {{ .ContentTitle | magenta | cyan }} {{ .MediaTypeEmoji }}",
	Details: `
|- Media Request Details -|
{{ "Title:" | faint }}	{{ .ContentTitle }}
{{ "Request ID:" | faint }}	{{ .ID }}
{{ "Media Type:" | faint }}	{{ .MediaType }}
{{ "Request Status:" | faint }}	{{ .StatusEmoji }} {{ .Status }}
{{ "Media Status:" | faint }}	{{ .MediaStatusEmoji }} {{ .MediaStatus }}
{{ "Creator:" | faint }}	{{ .CreatorEmail }} [{{ .CreatedDate }}]`,
}

func selectRequest(requests []*goverseerr.MediaRequest) *goverseerr.MediaRequest {
	ui.StartSpinner()
	friendly := make([]*goverseerr.FriendlyMediaRequest, len(requests))
	var wg sync.WaitGroup
	for idx, req := range requests {
		wg.Add(1)
		go func(i int, r *goverseerr.MediaRequest) {
			defer wg.Done()
			f, err := r.ToFriendly(overseerr)
			if err != nil {
				if i == 0 {
					fmt.Printf("%+v\n", f)
				}
				logrus.New().WithFields(logrus.Fields{
					"requestId": r.ID,
					"extended":  err.Error(),
				}).Errorln("failed to convert request to friendly request")
			}
			friendly[i] = f
		}(idx, req)
	}
	wg.Wait()
	ui.StopSpinner()
	idx, _ := ui.SelectorTemplated("Select a request:", friendly, requestSelectionTemplate)
	return requests[idx]
}

// Request Actions

var requestActionsPromptTemplate = &promptui.SelectTemplates{
	Label:    "{{ . }}:",
	Active:   "➡️  {{ .Name | magenta }}",
	Inactive: "   {{ .Name | cyan }}",
	Selected: "➡️  {{ .Name | magenta }}",
	Details: `
|- Action Details
{{ "Name:" | faint }}	{{ .Name }}
{{ "Description:" | faint }}	{{ .Description }}`,
}

type RequestAction struct {
	Name        string
	Description string
	Handler     func(r *goverseerr.MediaRequest)
	Validator   func(r *goverseerr.MediaRequest) bool
}

var RequestActions = []RequestAction{
	{
		Name:        "Delete Request",
		Description: "Remove a request from the Overseerr server",
		Validator: func(r *goverseerr.MediaRequest) bool {
			return true
		},
		Handler: func(r *goverseerr.MediaRequest) {
			ui.DestructiveConfirmation()
			if err := overseerr.DeleteRequest(r.ID); err != nil {
				ui.Error("Failed to delete that request")
				ui.ColorPrint("You can only remove pending requests if you are not an admin user\n", ui.White)
				ui.FatalQuiet("Failed to delete request from server", err)
			} else {
				ui.Success("Request Deleted")
			}
		},
	},
	{
		Name:        "Retry Request",
		Description: "Resend a request to the relevant management service",
		Validator: func(r *goverseerr.MediaRequest) bool {
			return r.Status != goverseerr.RequestStatusAvailable
		},
		Handler: func(r *goverseerr.MediaRequest) {
			if _, err := overseerr.RetryRequest(r.ID); err != nil {
				ui.Error("Failed to retry that request")
				ui.ColorPrint("You can only retry requests if you can manage requests (admin)\n", ui.White)
				ui.ColorPrint("Availabe requests can not be retried\n", ui.White)
				ui.FatalQuiet("Failed to retry request", err)
			} else {
				ui.Success("Request Resent")
			}
		},
	},
	{
		Name:        "Approve Request",
		Description: "Set a request's status to Approved",
		Validator: func(r *goverseerr.MediaRequest) bool {
			return (r.Status != goverseerr.RequestStatusApproved && r.Status != goverseerr.RequestStatusAvailable)
		},
		Handler: func(r *goverseerr.MediaRequest) {
			ui.DestructiveConfirmation()
			if _, err := overseerr.ApproveRequest(r.ID); err != nil {
				ui.Error("Failed to approve that request")
				ui.FatalQuiet("Failed to approve request", err)
			} else {
				ui.Success("Request Approved")
			}
		},
	},
	{
		Name:        "Decline Request",
		Description: "Set a request's status to Declined",
		Validator: func(r *goverseerr.MediaRequest) bool {
			return (r.Status != goverseerr.RequestStatusDeclined && r.Status != goverseerr.RequestStatusAvailable)
		},
		Handler: func(r *goverseerr.MediaRequest) {
			ui.DestructiveConfirmation()
			if _, err := overseerr.DeclineRequest(r.ID); err != nil {
				ui.Error("Failed to decline that request")
				ui.FatalQuiet("Failed to decline request", err)
			} else {
				ui.Success("Request Declined")
			}
		},
	},
	{
		Name:        "Details",
		Description: "Print more details about this request",
		Validator: func(r *goverseerr.MediaRequest) bool {
			return true
		},
		Handler: func(r *goverseerr.MediaRequest) {
			if friendly, err := r.ToFriendly(overseerr); err != nil {
				ui.Error("Failed to get more info about this request")
				ui.FatalQuiet("Failed to get request details", err)
			} else {
				fmt.Printf("%+v\n", friendly)
			}
		},
	},
}
