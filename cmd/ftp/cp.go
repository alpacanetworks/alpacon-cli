package ftp

import (
	"github.com/alpacanetworks/alpacon-cli/api/ftp"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
	"strings"
)

var CpCmd = &cobra.Command{
	Use:   "cp [SOURCE] [DESTINATION]",
	Short: "Copy files between local and remote locations",
	Long: `The cp command allows you to copy files between your local machine and a remote server.
	Example usages:
	- Upload: alpacon cp /path/to/local/file.txt servername:/remote/path/
	- Download: alpacon cp servername:/remote/file.txt /local/path/`,
	Example: `
	Upload: alpacon cp /path/to/local/file.txt servername:/remote/path/
	Download: alpacon cp servername:/remote/file.txt /local/path/
	`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]
		dest := args[1]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		if isRemotePath(src) && isLocalPath(dest) {
			downloadFile(alpaconClient, src, dest)
		} else if isLocalPath(src) && isRemotePath(dest) {
			uploadFile(alpaconClient, src, dest)
		} else {
			utils.CliError("Invalid combination of source and destination paths. Please use the format 'cp [SOURCE] [DESTINATION]' for copying files.")
		}
	},
}

// isRemotePath determines if the given path is a remote server path.
func isRemotePath(path string) bool {
	// Implement your logic to determine if the path is a remote path
	// For example, check if it contains ':'
	return strings.Contains(path, ":")
}

// isLocalPath determines if the given path is a local file system path.
func isLocalPath(path string) bool {
	return !isRemotePath(path)
}

func downloadFile(client *client.AlpaconClient, src, dest string) {
	err := ftp.DownloadFile(client, src, dest)
	if err != nil {
		utils.CliError("Failed to download the file from server: %s", err)
		return
	}
	utils.CliInfo("`%s` successfully downloaded from server `%s`", src, dest)
}

func uploadFile(client *client.AlpaconClient, src, dest string) {
	err := ftp.UploadFile(client, src, dest)
	if err != nil {
		utils.CliError("Failed to upload the file to server %s", err)
		return
	}
	utils.CliInfo("`%s` successfully uploaded to `%s` ", src, dest)
}
