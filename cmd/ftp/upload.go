package ftp

import (
	"github.com/alpacanetworks/alpacon-cli/api/ftp"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
	"strings"
)

var UploadCmd = &cobra.Command{
	Use:     "upload [FILE PATH] [SERVER NAME]:[UPLOAD PATH]",
	Aliases: []string{"cp"},
	Short:   "Transfer a file to a remote server",
	Example: `
	alpacon upload alpacon.txt myserver:/home/alpacon/
	alpacon cp /Users/alpacon.txt myserver:/home/alpacon/
	`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		upload := strings.SplitN(args[1], ":", 2)

		if len(upload) != 2 {
			utils.CliError("invalid argument format for [SERVER NAME]:[UPLOAD PATH]. Please provide the server name and path separated by a colon")
		}

		serverName := upload[0]
		path := upload[1]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		err = ftp.UploadFile(alpaconClient, file, serverName, path)
		if err != nil {
			utils.CliError("Failed to upload the file to server %s", err)
		}

		utils.CliInfo("`%s` successfully uploaded to server `%s` at `%s` ", file, serverName, path)
	},
}
