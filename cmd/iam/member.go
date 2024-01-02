package iam

import (
	"errors"
	"github.com/spf13/cobra"
)

var MemberCmd = &cobra.Command{
	Use:   "member",
	Short: "Manage group members",
	Long: `
	Member command provides tools for managing group members. 
	It includes functionalities to add, delete, and modify member roles within groups. 
	Use this command to oversee group membership and control access to group resources.
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
	MemberCmd.AddCommand(memberAddCmd)
	MemberCmd.AddCommand(memberDeleteCmd)
}
