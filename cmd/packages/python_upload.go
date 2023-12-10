package packages

import (
	"github.com/alpacanetworks/alpacon-cli/api/packages"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var pythonPackageUploadCmd = &cobra.Command{
	Use:     "upload [FILE PATH]",
	Aliases: []string{"cp"},
	Short:   "Upload a python package to alpacon",
	Long: `
	The 'upload' command allows users to upload a Python package to the alpacon. 
	This command is designed to facilitate the transfer of your locally developed Python packages to a remote server environment for further usage or distribution.
	`,
	Example: `
	alpacon package python upload alpamon-1.1.0-py3-none-any.whl
	alpacon package python cp /home/alpacon/alpamon-1.1.0-py3-none-any.whl
	`,
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		err = packages.UploadPackage(alpaconClient, file, "python")
		if err != nil {
			utils.CliError("Failed to upload the python packages to alpacon %s", err)
		}

		utils.CliInfo("`%s` successfully uploaded to alpacon ", file)
	},
}
