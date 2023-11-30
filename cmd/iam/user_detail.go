package iam

import (
	"github.com/alpacanetworks/alpacon-cli/api/iam"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var userDetailCmd = &cobra.Command{
	Use:   "describe [USER NAME]",
	Short: "Display detailed information about a specific user",
	Long: `
	The describe command fetches and displays detailed information about a specific user, 
	including its description, shell and other relevant attributes. 
	`,
	Example: ` 
	# Display details of a user named 'admin'
  	alpacon user describe admin
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userName := args[0]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Failed to create a connection to the Alpacon API %s", err)
		}

		userDetail, err := iam.GetUserDetail(alpaconClient, userName)
		if err != nil {
			utils.CliError("Failed to retrieve the user details %s", err)
		}

		utils.PrintJson(userDetail)
	},
}
