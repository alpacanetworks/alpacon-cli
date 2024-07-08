package ftp

import (
	"fmt"
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
	Copy files between your local machine and a remote server using the cp command.
	This command supports uploading, downloading, and specifying authentication details
	such as username and groupname.
	
	Example usages:
	- To upload multiple files to a remote server:
	  alpacon cp /local/path/file1.txt /local/path/file2.txt [SERVER_NAME]:/remote/path/
	
	- To download files from a remote server to a local destination:
	  alpacon cp [SERVER_NAME]:/remote/path1 /remote/path2 /local/destination/path
	
	- To specify username and groupname:
	  alpacon cp -u [USER_NAME] -g [GROUP_NAME] /local/path/file.txt [SERVER_NAME]:/remote/path/
	`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		groupname, _ := cmd.Flags().GetString("groupname")

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
			uploadFile(alpaconClient, sources, dest, username, groupname)
		} else if isRemotePath(sources[0]) && isLocalPath(dest) {
			downloadFile(alpaconClient, sources[0], dest, username, groupname)
		} else {
			utils.CliError("Invalid combination of source and destination paths.")
		}
	},
}

func init() {
	var username, groupname string

	CpCmd.Flags().StringVarP(&username, "username", "u", "", "Specify username")
	CpCmd.Flags().StringVarP(&groupname, "groupname", "g", "", "Specify groupname")
}

// isRemotePath determines if the given path is a remote server path.
func isRemotePath(path string) bool {
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

func downloadFile(client *client.AlpaconClient, src, dest, username, groupname string) {
	err := ftp.DownloadFile(client, src, dest, username, groupname)
	if err != nil {
		utils.CliError("Failed to download the file from server: %s.", err)
		return
	}
	utils.CliInfo("Download request for %s to server %s successful.", src, dest)
}

func uploadFile(client *client.AlpaconClient, src []string, dest, username, groupname string) {
	result, err := ftp.UploadFile(client, src, dest, username, groupname)
	if err != nil {
		utils.CliError("Failed to upload the file to server: %s.", err)
	}
	utils.CliInfo("Upload request for %s to %s successful.", src, dest)
	fmt.Printf("Result: %s.\n", result)
}
