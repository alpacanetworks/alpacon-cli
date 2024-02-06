package token

import (
	"errors"
	"github.com/spf13/cobra"
)

var TokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Commands to manage api tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("subcommand error")
	},
}

func init() {
	TokenCmd.AddCommand(tokenCreateCmd)
	TokenCmd.AddCommand(tokenListCmd)
	TokenCmd.AddCommand(tokenDeleteCmd)
}
