package server

import (
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var serverDetailCmd = &cobra.Command{
	Use:     "describe [SERVER NAME]",
	Aliases: []string{"desc"},
	Short:   "Display detailed information about a specific server",
	Long: `
	The describe command fetches and displays detailed information about a specific server, 
	including its status, and other relevant attributes. 
	This command is useful for getting an in-depth understanding of a server's current state and configuration.
	`,
	Example: ` 
	# Display details of a server named 'myserver'
  	alpacon server describe myserver
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		serverDetail, err := server.GetServerDetail(alpaconClient, serverName)
		if err != nil {
			utils.CliError("Failed to retrieve the server details: %s.", err)
		}

		utils.PrintJson(serverDetail)
	},
}
