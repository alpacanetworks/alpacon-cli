package token

import (
	"github.com/alpacanetworks/alpacon-cli/api/auth"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var tokenDeleteCmd = &cobra.Command{
	Use:   "delete [tok NAME]",
	Short: "Delete a specified api token",
	Long: `
	Removes an existing API token from the system. 
	This command requires the token name to identify the token to be deleted.
	`,
	Example: ` 
	alpacon token delete [TOKEN NAME]	
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tokenName := args[0]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		err = auth.DeleteAPIToken(alpaconClient, tokenName)
		if err != nil {
			utils.CliError("Failed to delete the api token %s. ", err)
		}

		utils.CliInfo("API Token successfully deleted: %s", tokenName)
	},
}
