package authority

import (
	"github.com/alpacanetworks/alpacon-cli/api/cert"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var authorityListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "Display a list of all certificate authorities",
	Long: `
 	Displays a comprehensive list of all certificate authorities that have been initialized within the system, 
	including their status and configuration details
	`,
	Example: `
	alpacon authority ls
	alpacon authority list
	alpacon authority all
	`,
	Run: func(cmd *cobra.Command, args []string) {

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		authorityList, err := cert.GetAuthorityList(alpaconClient)
		if err != nil {
			utils.CliError("Failed to retrieve the authority list: %s.", err)
		}

		utils.PrintTable(authorityList)
	},
}
