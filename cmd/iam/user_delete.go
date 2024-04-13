package iam

import (
	"github.com/alpacanetworks/alpacon-cli/api/iam"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var userDeleteCmd = &cobra.Command{
	Use:     "delete [USER NAME]",
	Aliases: []string{"rm"},
	Short:   "Delete a specified user",
	Long: `
	This command is used to permanently delete a specified user account from the Alpacon. 
	The command requires the exact username as an argument.
	`,
	Example: ` 
	alpacon user delete [USER NAME]	
	alpacon user rm [USER NAME]
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userName := args[0]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		if alpaconClient.Privileges == "general" {
			utils.CliError("You do not have the permission to delete users.")
		}

		err = iam.DeleteUser(alpaconClient, userName)
		if err != nil {
			utils.CliError("Failed to delete the user: %s.", err)
		}

		utils.CliInfo("User successfully deleted: %s.", userName)
	},
}
