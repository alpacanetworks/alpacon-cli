package packages

import (
	"github.com/alpacanetworks/alpacon-cli/api/packages"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var systemPackageListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "Display a list of all system packages",
	Long: `
	Display a detailed list of all python packages registered in the Alpacon.
	This command provides information such as name, version, platform and other relevant details.
	`,
	Example: `
	alpacon package system ls
	alpacon package system list
	alpacon package system all
	`,
	Run: func(cmd *cobra.Command, args []string) {
		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		packageList, err := packages.GetSystemPackageEntry(alpaconClient)
		if err != nil {
			utils.CliError("Failed to retrieve the system packages %s", err)
		}

		utils.PrintTable(packageList)
	},
}
