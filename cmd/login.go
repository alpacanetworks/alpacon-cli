package cmd

import (
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/auth"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/config"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
	"strings"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Alpacon",
	Long:  "Log in to Alpacon. To access Alpacon, workspace url is must specified",
	Example: `
	alpacon login

	alpacon login [WORKSPACE_URL] -u [USERNAME] -p [PASSWORD]
	alpacon login example.alpacon.io
	
	# Include http if using localhost.
	alppacon login http://localhost:8000
	
	# Login via API Token
	alpacon login -w [WORKSPACE_URL] -t [TOKEN_KEY]
	`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var workspaceURL string
		if len(args) > 0 {
			workspaceURL = args[0]
		} else {
			workspaceURL = ""
		}

		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		token, _ := cmd.Flags().GetString("token")

		if (workspaceURL == "" || username == "" || password == "") && token == "" {
			workspaceURL, username, password = promptForCredentials(workspaceURL, username, password)
		}

		if !strings.HasPrefix(workspaceURL, "http") {
			workspaceURL = "https://" + workspaceURL
		}

		loginRequest := &auth.LoginRequest{
			WorkspaceURL: workspaceURL,
			Username:     username,
			Password:     password,
		}

		fmt.Printf("Logging in to %s...\n", workspaceURL)
		err := auth.LoginAndSaveCredentials(loginRequest, token)
		if err != nil {
			utils.CliError("Login failed %v. Please check your credentials and try again.\n", err)
		}

		_, err = client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		fmt.Println("Login succeeded!")
	},
}

func init() {
	var username, password, token string

	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Username for login")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Password for login")
	loginCmd.Flags().StringVarP(&token, "token", "t", "", "API token for login")
}

func promptForCredentials(workspaceURL, username, password string) (string, string, string) {
	if workspaceURL == "" {
		configFile, err := config.LoadConfig()
		if err == nil && configFile.WorkspaceURL != "" {
			workspaceURL = configFile.WorkspaceURL
			fmt.Printf("Using Workspace URL %s from config file.\n", configFile.WorkspaceURL)
			fmt.Println("If you want to change the workspace, specify workspace url: alpacon login [WORKSPACE_URL] -u [USERNAME] -p [PASSWORD]")
			fmt.Println()
		}
	}

	if username == "" {
		username = utils.PromptForRequiredInput("Username: ")
	}

	if password == "" {
		password = utils.PromptForPassword("Password: ")
	}

	return workspaceURL, username, password
}
