package event

import (
	"github.com/alpacanetworks/alpacon-cli/api/event"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var pageSize int
var serverName string
var userName string

var EventCmd = &cobra.Command{
	Use:     "event",
	Aliases: []string{"events"},
	Short:   "Retrieve and display recent Alpacon events.",
	Long: `
	Retrieve and display a list of recent events from the Alpacon, with options to filter by server, user, and the number of events. 
	Use the '--tail' flag to limit the output to the last N event entries. 
	Specify a server with '--server' or filter events by user with '--user' to narrow down the results.
	`,
	Example: `
	alpacon event
	alpacon events
	alpacon event -tail 10 -s myserver -u admin
	alpacon event --tail=10 --server=myserver --user=admin
	`,
	Run: func(cmd *cobra.Command, args []string) {
		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		eventList, err := event.GetEventList(alpaconClient, pageSize, serverName, userName)
		if err != nil {
			utils.CliError("Failed to get event %s", err)
		}

		utils.PrintTable(eventList)
	},
}

func init() {
	EventCmd.Flags().IntVarP(&pageSize, "tail", "t", 25, "Number of event entries to show from the end")
	EventCmd.Flags().StringVarP(&serverName, "server", "s", "", "Specify server for events")
	EventCmd.Flags().StringVarP(&userName, "user", "u", "", "Specify request user for events")
}
