package ftp

import (
	"github.com/alpacanetworks/alpacon-cli/api/ftp"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
	"strings"
)

var CpCmd = &cobra.Command{
	Use:   "cp [SOURCE...] [DESTINATION]",
	Short: "Copy files between local and remote locations",
	Long: `The cp command allows you to copy files between your local machine and a remote server.
	Example usages:
	- Upload: alpacon cp /path/to/local/file1.txt /path/to/local/file2.txt servername:/remote/path/
	- Download: alpacon cp servername:"/remote/path1 /remote/path2" /local/destination/path`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			utils.CliError("You must specify at least two arguments.")
			return
		}

		dest := args[len(args)-1]
		sources := args[:len(args)-1]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
			return
		}

		if isLocalPaths(sources) && isRemotePath(dest) {
			uploadFile(alpaconClient, sources, dest)
		} else if isRemotePath(sources[0]) && isLocalPath(dest) {
			downloadFile(alpaconClient, sources[0], dest)
		} else {
			utils.CliError("Invalid combination of source and destination paths.")
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

func isLocalPaths(paths []string) bool {
	for _, path := range paths {
		if isRemotePath(path) {
			return false
		}
	}
	return true
}

func downloadFile(client *client.AlpaconClient, src string, dest string) {
	err := ftp.DownloadFile(client, src, dest)
	if err != nil {
		utils.CliError("Failed to download the file from server: %s", err)
		return
	}
	utils.CliInfo("`%s` successfully downloaded from server `%s`", src, dest)
}

func uploadFile(client *client.AlpaconClient, src []string, dest string) {
	err := ftp.UploadFile(client, src, dest)
	if err != nil {
		utils.CliError("Failed to upload the file to server %s", err)
	}
	utils.CliInfo("`%s` successfully uploaded to `%s` ", src, dest)
}
