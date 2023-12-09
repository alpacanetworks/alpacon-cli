package packages

import (
	"errors"
	"github.com/spf13/cobra"
)

var PackagesCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"packages"},
	Short:   "Commands to manage and interact with packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("subcommand error")
	},
}

func init() {
	PackagesCmd.AddCommand(systemCmd)
	PackagesCmd.AddCommand(pythonCmd)
}
