package cert

import (
	"errors"
	"github.com/spf13/cobra"
)

var CertCmd = &cobra.Command{
	Use:     "cert",
	Aliases: []string{"certificate"},
	Short:   "Manage and interact with SSL/TLS certificates",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("subcommand error")
	},
}

func init() {
	CertCmd.AddCommand(certListCmd)
}
