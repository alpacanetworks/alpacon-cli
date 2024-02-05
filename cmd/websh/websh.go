package websh

import (
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/event"
	"github.com/alpacanetworks/alpacon-cli/api/websh"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
	"strings"
)

var WebshCmd = &cobra.Command{
	Use:   "websh [SERVER NAME] [COMMAND]",
	Short: "Open a websh terminal or execute a command on a server",
	Long: ` 
	This command either opens a websh terminal for interacting with the specified server or executes a specified command directly on the server.
	It provides a terminal interface for managing and controlling the server remotely or for executing commands and retrieving their output directly.
	`,
	Example: `
	// Open a websh terminal for a server
	alpacon websh [SERVER NAME]

	// Execute a command directly on a server and retrieve the output
	alpacon websh [SERVER NAME] [COMMAND]

	// Additional examples with flags
	alpacon websh -r [SERVER_NAME]
	alpacon websh [SERVER NAME] --root
	`,
	Run: func(cmd *cobra.Command, args []string) {
		root, _ := cmd.Flags().GetBool("root")

		if len(args) < 1 {
			cmd.Usage()
			return
		}

		serverName := args[0]
		commandArgs := args[1:]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		if len(commandArgs) > 0 {
			command := strings.Join(commandArgs, " ")
			result, err := event.RunCommand(alpaconClient, serverName, command)
			if err != nil {
				utils.CliError("Failed to run the '%s' command on the '%s' server: %s", command, serverName, err)
			}
			fmt.Println(result)
		} else {
			session, err := websh.CreateWebshConnection(alpaconClient, serverName, root)
			if err != nil {
				utils.CliError("Failed to create the websh connection: %s", err)
			}

			websh.OpenNewTerminal(alpaconClient, session)
		}
	},
}

func init() {
	var root bool
	WebshCmd.Flags().BoolVarP(&root, "root", "r", false, "Run as root user")
}
