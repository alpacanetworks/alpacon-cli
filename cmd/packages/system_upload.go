package packages

import (
	"github.com/alpacanetworks/alpacon-cli/api/packages"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var systemPackageUploadCmd = &cobra.Command{
	Use:     "upload [FILE PATH]",
	Aliases: []string{"cp"},
	Short:   "Upload a system package to alpacon",
	Long: `
	The 'upload' command allows users to upload a System package to the alpacon. 
	This command is designed to facilitate the transfer of System packages to a remote server environment for further usage or distribution.
	`,
	Example: `
	alpacon package system upload osquery-5.10.2-1.linux.x86_64.rpm
	alpacon package system cp /home/alpacon/osquery_5.8.2-1.linux_amd64.deb
	`,
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		err = packages.UploadPackage(alpaconClient, file, "system")
		if err != nil {
			utils.CliError("Failed to upload the system packages to alpacon %s", err)
		}

		utils.CliInfo("`%s` successfully uploaded to alpacon ", file)
	},
}
