package cmd

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/willfantom/goverseerr"
	"github.com/willfantom/overclirr/ui"
)

var retryAllCmd = &cobra.Command{
	Use:   "retry-requests",
	Short: "Retry all the non-available requests",
	Long:  "Resend all media requests to the manager service provided the content is not already available",
	PreRun: func(cmd *cobra.Command, args []string) {
		setOverseer(overseerrProfileName)
	},
	Run: func(cmd *cobra.Command, args []string) {
		ui.ColorPrintBold("Fetching Requests...\n", ui.White)
		ui.StartSpinner()
		pgNumber := 0
		anyErrors := false
		var allRequests []*goverseerr.MediaRequest
		for {
			r, p, err := overseerr.GetRequests(pgNumber, 50, goverseerr.RequestFileterUnavailable, goverseerr.RequestSortAdded)
			if err != nil {
				ui.Fatal("Could not get the unavailable request list", err)
			}
			allRequests = append(allRequests, r...)
			pgNumber++
			if pgNumber >= p.Pages {
				break
			}
		}
		ui.StopSpinner()
		ui.ColorPrint("Retrying Requests...\n", ui.White)
		ui.StartSpinner()
		var wg sync.WaitGroup
		for _, req := range allRequests {
			wg.Add(1)
			go func(r *goverseerr.MediaRequest) {
				defer wg.Done()
				if _, err := overseerr.RetryRequest(r.ID); err != nil {
					ui.Error("Failed to retry request: " + fmt.Sprintf("%d", r.ID))
					anyErrors = true
				}
			}(req)
		}
		wg.Wait()
		ui.StopSpinner()
		if anyErrors {
			ui.ColorPrintBold("Completed with errors\n", ui.Red)
		} else {
			ui.Success("Retried all requests")
		}
	},
}

func init() {
	RootCmd.AddCommand(retryAllCmd)
}
