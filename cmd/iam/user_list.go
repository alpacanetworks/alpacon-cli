package iam

import (
	"github.com/alpacanetworks/alpacon-cli/api/iam"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var userListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "Display a list of all users",
	Long: `
	Display a detailed list of all users registered in the Alpacon.
	This command provides information such as name, email, status, and other relevant details.
	`,
	Example: `
	alpacon user ls
	alpacon user list
	alpacon user all
	`,
	Run: func(cmd *cobra.Command, args []string) {
		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Failed to create a connection to the Alpacon API %s", err)
		}

		userList, err := iam.GetUserList(alpaconClient)
		if err != nil {
			utils.CliError("Failed to retrieve the user list %s", err)
		}

		utils.PrintTable(userList)
	},
}
