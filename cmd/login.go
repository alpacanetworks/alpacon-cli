package cmd

import (
	"bufio"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/auth"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"os"
	"strings"
)

var (
	loginRequest auth.LoginRequest
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Alpacon Server",
	Long:  "Log in to Alpacon Server.\n To access Alpacon Server, server address is must specified",
	Run: func(cmd *cobra.Command, args []string) {
		performLogin()
	},
}

func init() {
	loginCmd.Flags().StringVarP(&loginRequest.ServerAddress, "server", "s", "", "URL of the server to login")
	loginCmd.Flags().StringVarP(&loginRequest.Username, "username", "u", "", "Username for login")
	loginCmd.Flags().StringVarP(&loginRequest.Password, "password", "p", "", "Password for login")
}

func performLogin() {
	if loginRequest.Username == "" || loginRequest.Password == "" || loginRequest.ServerAddress == "" {
		promptForCredentials()
	}

	err := auth.LoginAndSaveCredentials(&loginRequest)
	if err != nil {
		fmt.Printf("Login failed: %v. Please check your credentials and try again.\n", err)
		os.Exit(1)
	}

	_, err = client.NewAlpaconAPIClient()
	if err != nil {
		fmt.Printf("Error creating Alpacon client: %v. Please retry the login process.\n", err)
		os.Exit(1)
	}
	fmt.Println("Login succeeded!")
}

// promptForCredentials
func promptForCredentials() {
	reader := bufio.NewReader(os.Stdin)

	if loginRequest.Username == "" || loginRequest.Password == "" || loginRequest.ServerAddress == "" {
		if loginRequest.Username == "" {
			fmt.Print("Username: ")
			username, _ := reader.ReadString('\n')
			loginRequest.Username = strings.TrimSpace(username)
		}
		if loginRequest.Password == "" {
			fmt.Print("Password: ")
			bytePassword, _ := term.ReadPassword(0)
			loginRequest.Password = strings.TrimSpace(string(bytePassword))
			fmt.Println()
		}
		if loginRequest.ServerAddress == "" {
			fmt.Print("Server Address: ")
			serverAddress, _ := reader.ReadString('\n')
			loginRequest.ServerAddress = strings.TrimSpace(serverAddress)
		}
	}
}
