package server

import (
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var serverListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "Display a list of all servers",
	Long: `
	Display a detailed list of all servers registered in the Alpacon.
	This command provides information such as server ID, name, status, and other relevant details.
	`,
	Example: `
	alpacon server ls
	alpacon server list
	alpacon server all
	`,
	Run: func(cmd *cobra.Command, args []string) {
		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Failed to create a connection to the Alpacon API %s", err)
		}

		serverList, err := server.GetServerList(alpaconClient)
		if err != nil {
			utils.CliError("Failed to retrieve the servers %s", err)
		}

		utils.PrintTable(serverList)
	},
}
