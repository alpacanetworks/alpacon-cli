package cert

import (
	"github.com/alpacanetworks/alpacon-cli/api/cert"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var certListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "Display a list of all certificates",
	Long: `
	Retrieves and shows a detailed list of all the SSL/TLS certificates currently managed by the system, 
	including their issuance status and validity.
	`,
	Example: `
	alpacon cert ls
	alpacon cert list
	alpacon cert all
	`,
	Run: func(cmd *cobra.Command, args []string) {
		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		certList, err := cert.GetCertificateList(alpaconClient)
		if err != nil {
			utils.CliError("Failed to retrieve the certificate list: %s.", err)
		}

		utils.PrintTable(certList)
	},
}
