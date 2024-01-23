package csr

import (
	"github.com/alpacanetworks/alpacon-cli/api/cert"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var csrDetailCmd = &cobra.Command{
	Use:     "describe [CSR ID]",
	Aliases: []string{"desc"},
	Short:   "Display detailed information about a specific Certificate Signing Request",
	Long: `
	The describe command fetches and displays detailed information about a specific certificate signing request, 
	including its csr text, certificate authority and other relevant attributes. 
	`,
	Example: ` 
	alpacon csr describe [CSR ID]
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csrId := args[0]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		csrDetail, err := cert.GetCSRDetail(alpaconClient, csrId)
		if err != nil {
			utils.CliError("Failed to retrieve the csr details %s", err)
		}

		utils.PrintJson(csrDetail)
	},
}
