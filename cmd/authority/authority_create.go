package authority

import (
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/cert"
	"github.com/alpacanetworks/alpacon-cli/api/iam"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var authorityCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new certificate authority",
	Long: `
  	Initializes a new certificate authority within the system, allowing you to sign and manage certificates and define their policies.
	`,
	Example: `alpacon authority create`,
	Run: func(cmd *cobra.Command, args []string) {

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		authorityRequest := promptForAuthority(alpaconClient)

		response, err := cert.CreateAuthority(alpaconClient, authorityRequest)
		if err != nil {
			utils.CliError("Failed to create the new authority: %s.", err)
		}

		installAuthorityInfo(response)
	},
}

func promptForAuthority(ac *client.AlpaconClient) cert.AuthorityRequest {
	var authorityRequest cert.AuthorityRequest

	authorityRequest.Name = utils.PromptForRequiredInput("Common name for the CA. (e.g., Alapca Networks' Root CA): ")
	authorityRequest.Organization = utils.PromptForRequiredInput("Organization name that this CA belongs to. (e.g., Alpaca Networks): ")
	authorityRequest.Domain = utils.PromptForRequiredInput("Domain name of the root certificate: ")
	authorityRequest.RootValidDays = utils.PromptForIntInput("Root certificate validity in days (10 years = 3650): ")
	authorityRequest.DefaultValidDays = utils.PromptForIntInput("Child certificate validity in days (3 months = 90, 1 year = 365): ")
	authorityRequest.MaxValidDays = utils.PromptForIntInput("Maximum valid days that users can request: ")

	agent := utils.PromptForRequiredInput("Name of sever to run this CA on: ")
	agentID, err := server.GetServerIDByName(ac, agent)
	if err != nil {
		utils.CliError("Failed to retrieve the server %s", err)
	}
	authorityRequest.Agent = agentID

	owner := utils.PromptForRequiredInput("Owner(username): ")
	ownerID, err := iam.GetUserIDByName(ac, owner)
	if err != nil {
		utils.CliError("Failed to retrieve the user %s", err)
	}
	authorityRequest.Owner = ownerID

	return authorityRequest
}

func installAuthorityInfo(response cert.AuthorityCreateResponse) {
	fmt.Println()
	utils.PrintHeader("Installation instruction")
	fmt.Println()
	fmt.Println(response.Instruction + "\n")
	utils.CliWarning("Please be aware that after leaving this page, you cannot obtain the script again for security.")
}
