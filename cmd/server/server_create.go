package server

import (
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/iam"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
	"strings"
)

var serverCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new server",
	Long: `
	Create a new server with specific configurations. This command allows you to set up a server with a unique name, 
	choose a platform, and define access permissions for different groups. 
	`,
	Example: `alpacon create server`,
	Run: func(cmd *cobra.Command, args []string) {

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		groupList, err := iam.GetGroupList(alpaconClient)
		if err != nil {
			utils.CliError("Failed to retrieve the group list %s", err)
		}

		serverRequest := promptForServer(alpaconClient, groupList)

		response, err := server.CreateServer(alpaconClient, serverRequest)
		if err != nil {
			utils.CliError("Failed to create the new server %s", err)
		}

		installServerInfo(response)
	},
}

func promptForServer(ac *client.AlpaconClient, groupList []iam.GroupAttributes) server.ServerRequest {
	var serverRequest server.ServerRequest

	serverRequest.Name = utils.PromptForRequiredInput("Server Name: ")
	serverRequest.Platform = promptForPlatform()

	displayGroups(groupList)

	serverRequest.Groups = selectAndConvertGroups(ac, groupList)

	return serverRequest
}

func promptForPlatform() string {
	for {
		platform := utils.PromptForInput("Platform(debian, rhel): ")
		if strings.ToLower(platform) == "debian" || strings.ToLower(platform) == "rhel" {
			return platform
		}
		fmt.Println("Invalid platform. Please choose 'debian' or 'rhel'.")
	}
}

func displayGroups(groupList []iam.GroupAttributes) {
	fmt.Println("Groups:")
	for i, group := range groupList {
		fmt.Printf("[%d] %s\n", i+1, group.Name)
	}
}

func selectAndConvertGroups(ac *client.AlpaconClient, groupList []iam.GroupAttributes) []string {
	chosenGroups := utils.PromptForRequiredInput("Select groups that are authorized to access this server. (e.g., 1,2):")
	intGroups := utils.SplitAndParseInts(chosenGroups)

	var groupIDs []string

	for _, groupIndex := range intGroups {
		if groupIndex < 1 || groupIndex > len(groupList) {
			utils.CliError(fmt.Sprintf("Invalid group index: %d", groupIndex))
		}

		groupID, err := iam.GetGroupIDByName(ac, groupList[groupIndex-1].Name)
		if err != nil {
			utils.CliError("No group found with the given name")
		}

		groupIDs = append(groupIDs, groupID)
	}

	return groupIDs
}

func installServerInfo(response server.ServerCreatedResponse) {
	fmt.Println()
	utils.PrintHeader("Connecting server to alpacon")
	printIntro()
	printMethod("Simply use our install script:", response.Instruction1)
	printMethod("Or, do it manually (If you've followed the script above, this is not required):", response.Instruction2)
	utils.CliWarning("Please be aware that after leaving this page, you cannot obtain the script again for security.")
}

func printIntro() {
	fmt.Println("We provide two ways to connect your server to alpacon.")
	fmt.Println("Please follow one of the following steps to install the \"alpamon\" agent on your server.")
}

func printMethod(header, instruction string) {
	fmt.Println(utils.Green(header))
	fmt.Println()
	fmt.Println(instruction + "\n")
}
