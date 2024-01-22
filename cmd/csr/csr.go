package csr

import (
	"errors"
	"github.com/spf13/cobra"
)

var CsrCmd = &cobra.Command{
	Use:   "csr",
	Short: "Generate and manage Certificate Signing Request (CSR) operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("subcommand error")
	},
}

func init() {
	CsrCmd.AddCommand(csrCreateCmd)
	CsrCmd.AddCommand(csrListCmd)
	CsrCmd.AddCommand(csrApproveCmd)
	CsrCmd.AddCommand(csrDenyCmd)
	CsrCmd.AddCommand(csrDeleteCmd)
	CsrCmd.AddCommand(csrDetailCmd)
}
