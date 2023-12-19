package agent

import (
	"github.com/alpacanetworks/alpacon-cli/api/agent"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var restartAgentCmd = &cobra.Command{
	Use:     "restart [SERVER NAME]",
	Short:   "Restart server's agent(alpamon)",
	Example: `alpacon agent restart myserver`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		err = agent.RequestAgentAction(alpaconClient, serverName, "restart")
		if err != nil {
			utils.CliError("Failed to restart the agent %s", err)
		}

		utils.CliInfo("Agent restart request successful. Verify in events.(alpacon events)")
	},
}
