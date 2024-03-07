package token

import (
	"github.com/alpacanetworks/alpacon-cli/api/security"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var aclDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the specified command ACL from an API token.",
	Long: `
	Removes an existing command acl from the API token
	This command requires the command acl id to identify the command acl to be deleted.
	`,
	Example: `
	alpacon token acl delete [COMMAND_ACL_ID]
	alpacon token acl delete --token=[TOKEN_ID_OR_NAME] --command=[COMMAND]
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commandAclId := args[0]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		err = security.DeleteCommandAcl(alpaconClient, commandAclId)
		if err != nil {
			utils.CliError("Failed to delete the command acl %s.", err)
		}

		utils.CliInfo("Command ACL successfully deleted: %s", commandAclId)
	},
}
