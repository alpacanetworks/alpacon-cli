package cmd

import (
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/cmd/agent"
	"github.com/alpacanetworks/alpacon-cli/cmd/authority"
	"github.com/alpacanetworks/alpacon-cli/cmd/cert"
	"github.com/alpacanetworks/alpacon-cli/cmd/csr"
	"github.com/alpacanetworks/alpacon-cli/cmd/event"
	"github.com/alpacanetworks/alpacon-cli/cmd/ftp"
	"github.com/alpacanetworks/alpacon-cli/cmd/iam"
	"github.com/alpacanetworks/alpacon-cli/cmd/log"
	"github.com/alpacanetworks/alpacon-cli/cmd/note"
	"github.com/alpacanetworks/alpacon-cli/cmd/packages"
	"github.com/alpacanetworks/alpacon-cli/cmd/server"
	"github.com/alpacanetworks/alpacon-cli/cmd/token"
	"github.com/alpacanetworks/alpacon-cli/cmd/websh"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "alpacon",
	Aliases: []string{"ac"},
	Short:   "Alpacon CLI: Your Gateway to Alpacon Services",
	Long:    "Use this tool to interact with the alpacon service.",
	Run: func(cmd *cobra.Command, args []string) {
		utils.ShowLogo()
		fmt.Println("Welcome to Alpacon CLI! Use 'alpacon [command]' to execute a specific command or 'alpacon help' to see all available commands.")
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		utils.CliError("While executing the command: %s", err)
	}
}

func init() {
	// version
	RootCmd.AddCommand(versionCmd)

	// login
	RootCmd.AddCommand(loginCmd)

	// logout
	RootCmd.AddCommand(logoutCmd)

	// iam
	RootCmd.AddCommand(iam.UserCmd)
	RootCmd.AddCommand(iam.GroupCmd)

	// server
	RootCmd.AddCommand(server.ServerCmd)

	// agent
	RootCmd.AddCommand(agent.AgentCmd)

	// websh
	RootCmd.AddCommand(websh.WebshCmd)

	// ftp
	RootCmd.AddCommand(ftp.CpCmd)

	// packages
	RootCmd.AddCommand(packages.PackagesCmd)

	// log
	RootCmd.AddCommand(log.LogCmd)

	// event
	RootCmd.AddCommand(event.EventCmd)

	// note
	RootCmd.AddCommand(note.NoteCmd)

	// authority
	RootCmd.AddCommand(authority.AuthorityCmd)

	// csr
	RootCmd.AddCommand(csr.CsrCmd)

	// certificate
	RootCmd.AddCommand(cert.CertCmd)

	// token
	RootCmd.AddCommand(token.TokenCmd)
}
