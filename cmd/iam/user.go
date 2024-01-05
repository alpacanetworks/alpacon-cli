package iam

import (
	"errors"
	"github.com/spf13/cobra"
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage User resources",
	Long: `
	The 'user' command is used to manage IAM user resources,
	including users and permissions within the Alpacon.
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
	UserCmd.AddCommand(userListCmd)
	UserCmd.AddCommand(userDetailCmd)
	UserCmd.AddCommand(userDeleteCmd)
	UserCmd.AddCommand(userCreateCmd)
}
