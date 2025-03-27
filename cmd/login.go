package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alpacanetworks/alpacon-cli/api/auth"
	"github.com/alpacanetworks/alpacon-cli/api/auth0"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/config"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
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
	alpacon login http://localhost:8000
	
	# Login via API Token
	alpacon login [WORKSPACE_URL] -t [TOKEN_KEY]
	`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var workspaceURL string

		// Determine the workspace URL to use
		if len(args) > 0 {
			workspaceURL = args[0]
		}

		if workspaceURL == "" {
			cfg, err := config.LoadConfig()
			if err == nil && cfg.WorkspaceURL != "" {
				workspaceURL = cfg.WorkspaceURL
			}
		}

		if workspaceURL == "" {
			workspaceURL = utils.PromptForRequiredInput("workspaceURL: ")
		}

		// Validate workspaceURL
		workspaceURL, err := validateAndFormatWorkspaceURL(workspaceURL)
		if err != nil {
			utils.CliError(err.Error())
		}

		// Check login method
		envInfo, err := auth0.FetchAuthEnv(workspaceURL)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				envInfo = &auth0.AuthEnvResponse{Method: "legacy"}
			} else {
				utils.CliError("Failed to patch environment variables from workspace. %v", err)
			}
		}

		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		token, _ := cmd.Flags().GetString("token")

		fmt.Printf("Logging in to %s\n", workspaceURL)
		if envInfo.Method == "auth0" && token == "" {
			deviceCode, err := auth0.RequestDeviceCode(workspaceURL, envInfo)
			if err != nil {
				utils.CliError("Device code request failed. %v", err)
			}

			highlight := "\033[1;34m" // blue + bold
			reset := "\033[0m"

			fmt.Println("\n==================== AUTHENTICATION REQUIRED ====================")
			fmt.Println("\nPlease authenticate by visiting the following URL:")
			fmt.Printf("%s%s%s\n\n", highlight, deviceCode.VerificationURIComplete, reset)
			fmt.Print("===============================================================\n\n")

			tokenRes, err := auth0.PollForToken(deviceCode, envInfo)
			if err != nil {
				utils.CliError(err.Error())
			}

			err = config.CreateConfig(workspaceURL, "", "", tokenRes.AccessToken, tokenRes.RefreshToken, tokenRes.ExpiresIn)
			if err != nil {
				utils.CliError("Failed to save config: %v", err)
			}

		} else {
			if (workspaceURL == "" || username == "" || password == "") && token == "" {
				workspaceURL, username, password = promptForCredentials(workspaceURL, username, password)
			}

			loginRequest := &auth.LoginRequest{
				WorkspaceURL: workspaceURL,
				Username:     username,
				Password:     password,
			}

			err = auth.LoginAndSaveCredentials(loginRequest, token)
			if err != nil {
				utils.CliError("Login failed %v. Please check your credentials and try again.", err)
			}

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

func validateAndFormatWorkspaceURL(workspaceURL string) (string, error) {
	if !strings.HasPrefix(workspaceURL, "http") {
		workspaceURL = "https://" + workspaceURL
	}

	resp, err := http.Get(workspaceURL)
	if err != nil || resp.StatusCode >= 400 {
		return "", fmt.Errorf("workspace URL is unreachable: %v", workspaceURL)
	}
	defer resp.Body.Close()

	return workspaceURL, nil
}
