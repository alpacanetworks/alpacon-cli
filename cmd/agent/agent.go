package agent

import (
	"errors"
	"github.com/spf13/cobra"
)

var AgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Commands to manage server's agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("subcommand error")
	},
}

func init() {
	AgentCmd.AddCommand(upgradeAgentCmd)
	AgentCmd.AddCommand(restartAgentCmd)
	AgentCmd.AddCommand(shutdownAgentCmd)
}
