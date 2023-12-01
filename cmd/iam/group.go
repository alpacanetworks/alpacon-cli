package iam

import (
	"errors"
	"github.com/spf13/cobra"
)

var GroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage Group (Identity and Access Management) resources",
	Long: `
	The 'group' command is used to manage IAM group resources,
	including groups, and permissions within the Alpacon.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("subcommand error")
	},
}

func init() {
	GroupCmd.AddCommand(groupListCmd)
	GroupCmd.AddCommand(groupDetailCmd)
}