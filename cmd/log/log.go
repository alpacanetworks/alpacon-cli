package log

import (
	"github.com/alpacanetworks/alpacon-cli/api/log"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var pageSize int

var LogCmd = &cobra.Command{
	Use:     "log [SERVER NAME]",
	Aliases: []string{"logs"},
	Short:   "Retrieve and display server logs",
	Long: `Retrieve and display logs for a specified server. This command allows you 
	to view logs of different levels and types associated with a server. Use the '--tail' flag 
	to limit the output to the last N log entries. Suitable for debugging and monitoring 
	server activities.`,
	Example: `
	alpacon log [SERVER NAME]
	alpacon logs [SERVER_NAME]
	alpacon log [SERVER NAME] --tail=10
	alpacon logs [SERVER NAME] --tail=10
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		logList, err := log.GetSystemLogList(alpaconClient, serverName, pageSize)
		if err != nil {
			utils.CliError("Failed to get logs: %s.", err)
		}

		utils.PrintTable(logList)
	},
}

func init() {
	LogCmd.Flags().IntVarP(&pageSize, "tail", "t", 25, "Number of log entries to show from the end")
}
