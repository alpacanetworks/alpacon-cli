package csr

import (
	"github.com/alpacanetworks/alpacon-cli/api/cert"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var csrDenyCmd = &cobra.Command{
	Use:   "deny",
	Short: "Deny a CSR",
	Long: `
	Rejects a Certificate Signing Request, marking it as denied and stopping any further processing 
	or signing activities for that request
	`,
	Example: `alpacon csr deny [CSR ID] `,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csrId := args[0]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		_, err = cert.DenyCSR(alpaconClient, csrId)
		if err != nil {
			utils.CliError("Failed to deny the csr: %s.", err)
		}

		utils.CliInfo("CSR denial request successful. Please verify the CSR status.")
	},
}
