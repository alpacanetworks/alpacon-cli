package cmd

import (
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/auth"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var (
	loginRequest auth.LoginRequest
)

const defaultServerURL = "https://alpacon.io"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Alpacon Server",
	Long:  "Log in to Alpacon Server.\n To access Alpacon Server, server address is must specified",
	Run: func(cmd *cobra.Command, args []string) {
		if loginRequest.Username == "" || loginRequest.Password == "" || loginRequest.ServerAddress == "" {
			promptForCredentials()
		}

		err := auth.LoginAndSaveCredentials(&loginRequest)
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
	loginCmd.Flags().StringVarP(&loginRequest.ServerAddress, "server", "s", "defaultServerURL", "URL of the server to login, default: https://alpacon.io")
	loginCmd.Flags().StringVarP(&loginRequest.Username, "username", "u", "", "Username for login")
	loginCmd.Flags().StringVarP(&loginRequest.Password, "password", "p", "", "Password for login")
}

func promptForCredentials() {
	if loginRequest.Username == "" {
		loginRequest.Username = utils.PromptForRequiredInput("Username: ")
	}
	if loginRequest.Password == "" {
		loginRequest.Password = utils.PromptForPassword("Password: ")
	}

	loginRequest.ServerAddress = utils.PromptForInput("Server Address[https://alpacon.io]: ")
	if loginRequest.ServerAddress == "" {
		loginRequest.ServerAddress = defaultServerURL
	}
}
