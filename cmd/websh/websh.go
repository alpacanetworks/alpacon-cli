package websh

import (
	"github.com/alpacanetworks/alpacon-cli/api/websh"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var root bool

var WebshCmd = &cobra.Command{
	Use:   "websh [SERVER NAME]",
	Short: "Open a websh terminal for a server",
	Long: ` 
	This command opens a websh terminal for interacting with the specified server. 
	It provides a terminal interface to manage and control the server remotely.
	`,
	Example: `
	alpacon websh [SERVER NAME]
	alpacon websh -r [SERVER_NAME]
	alpacon websh [SERVER NAME] --root
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		session, err := websh.CreateWebshConnection(alpaconClient, serverName, root)
		if err != nil {
			utils.CliError("Failed to create the websh connection: %s", err)
		}

		websh.OpenNewTerminal(alpaconClient, session)
	},
}

func init() {
	WebshCmd.Flags().BoolVarP(&root, "root", "r", false, "Run as root user")
}
