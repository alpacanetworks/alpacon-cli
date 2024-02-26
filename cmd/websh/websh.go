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
	alpacon websh [SERVER_NAME]

	// Execute a command directly on a server and retrieve the output
	alpacon websh [SERVER_NAME] [COMMAND]

	// Open a websh terminal as a root user
	alpacon websh -r [SERVER_NAME]
	alapcon websh -u root [SERVER_NAME]
	
	// Open a websh terminal specifying username and groupname
	alpacon websh -u [USER_NAME] -g [GROUP_NAME] [SERVER_NAME]

	// Run a command as [USER_NAME]/[GROUP_NAME]
	alpacon websh -u [USER_NAME] -g [GROUP_NAME] [SERVER_NAME] [COMMAND]

	Flags:
	-r          					   Run the websh terminal as the root user.
	-u / --username [USER_NAME]        Specify the username under which the command should be executed.
	-g / --groupname [GROUP_NAME]      Specify the group name under which the command should be executed.

	Note:
	- All flags (-r, -u, -g) must be placed before the [SERVER_NAME].
	- The -u (or --username) and -g (or --groupname) flags require an argument specifying the user or group name, respectively.
	`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			username, groupname, serverName string
			commandArgs                     []string
		)

		for i := 0; i < len(args); i++ {
			switch {
			case args[i] == "-r" || args[i] == "--root":
				username = "root"
			case args[i] == "-h" || args[i] == "--help":
				cmd.Help()
				return
			case strings.HasPrefix(args[i], "-u") || strings.HasPrefix(args[i], "--username"):
				username, i = extractValue(args, i)
			case strings.HasPrefix(args[i], "-g") || strings.HasPrefix(args[i], "--groupname"):
				groupname, i = extractValue(args, i)
			default:
				if serverName == "" {
					serverName = args[i]
				} else {
					commandArgs = append(commandArgs, args[i:]...)
					i = len(args)
				}
			}
		}

		if serverName == "" {
			utils.CliError("Server name is required.")
		}

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		if len(commandArgs) > 0 {
			command := strings.Join(commandArgs, " ")
			result, err := event.RunCommand(alpaconClient, serverName, command, username, groupname)
			if err != nil {
				utils.CliError("Failed to run the '%s' command on the '%s' server: %s", command, serverName, err)
			}
			fmt.Println(result)
		} else {
			session, err := websh.CreateWebshConnection(alpaconClient, serverName, username, groupname)
			if err != nil {
				utils.CliError("Failed to create the websh connection: %s", err)
			}
			websh.OpenNewTerminal(alpaconClient, session)
		}
	},
}

func extractValue(args []string, i int) (string, int) {
	if strings.Contains(args[i], "=") {
		parts := strings.SplitN(args[i], "=", 2)
		return parts[1], i
	}
	if i+1 < len(args) {
		return args[i+1], i + 1
	}
	return "", i
}
