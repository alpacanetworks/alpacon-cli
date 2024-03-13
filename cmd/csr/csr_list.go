package csr

import (
	"github.com/alpacanetworks/alpacon-cli/api/cert"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var csrListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "Display a list of all certificate signing requests",
	Long: `
	Display CSRs, optionally filtered by state ('requested', 'processing', 'signed', 
    'issued', 'canceled', 'denied'). Use the --state flag to specify the status.
	`,
	Example: `
	alpacon csr ls
	alpacon csr list
	alpacon csr all
	`,
	Run: func(cmd *cobra.Command, args []string) {
		state, _ := cmd.Flags().GetString("state")

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		csrList, err := cert.GetCSRList(alpaconClient, state)
		if err != nil {
			utils.CliError("Failed to retrieve the csr list: %s.", err)
		}

		utils.PrintTable(csrList)
	},
}

func init() {
	var state string

	csrListCmd.Flags().StringVarP(&state, "state", "s", "", "Specify the status of the CSR (e.g., 'denied', 'signed')")
}
