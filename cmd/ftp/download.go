package ftp

import (
	"github.com/alpacanetworks/alpacon-cli/api/ftp"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
	"strings"
)

var DownloadCmd = &cobra.Command{
	Use:     "download [SERVER NAME]:[FILE PATH]",
	Aliases: []string{"cp"},
	Short:   "Transfer a file from a remote server",
	Example: `
	alpacon download myserver:/home/alpacon/alpacon.txt
	alpacon cp myserver:/home/alpacon/alpacon.txt
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		upload := strings.SplitN(args[0], ":", 2)

		if len(upload) != 2 {
			utils.CliError("invalid argument format for [SERVER NAME]:[FILE PATH]. Please provide the server name and path separated by a colon")
		}

		serverName := upload[0]
		path := upload[1]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		downloadURL, err := ftp.DownloadFile(alpaconClient, serverName, path)
		if err != nil {
			utils.CliError("Failed to download the file from server: %s", err)
		}

		utils.CliInfo("`%s` successfully downloaded from server `%s`. Download URL: %s", path, serverName, downloadURL)
	},
}
