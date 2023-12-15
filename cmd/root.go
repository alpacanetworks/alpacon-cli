package cmd

import (
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/cmd/agent"
	"github.com/alpacanetworks/alpacon-cli/cmd/ftp"
	"github.com/alpacanetworks/alpacon-cli/cmd/iam"
	"github.com/alpacanetworks/alpacon-cli/cmd/log"
	"github.com/alpacanetworks/alpacon-cli/cmd/packages"
	"github.com/alpacanetworks/alpacon-cli/cmd/server"
	"github.com/alpacanetworks/alpacon-cli/cmd/websh"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

func showLogo() {
	alpaconLogo := `
    (` + "`" + `-')  _           _  (` + "`" + `-') (` + "`" + `-')  _                      <-. (` + "`" + `-')_
    (OO ).-/    <-.    \-.(OO ) (OO ).-/  _             .->      \( OO) )
    / ,---.   ,--. )   _.'    \ / ,---.   \-,-----.(` + "`" + `-')----. ,--./ ,--/
    | \ /` + ".`" + `\  |  (` + "`" + `-')(_...--'' | \ /` + ".`" + `\   |  .--./( OO).-.  '|   \ |  |
    '-'|_.' | |  |OO )|  |_.' | '-'|_.' | /_) (` + "`" + `-')( _) | |  ||  . '|  |)
    (|  .-.  |(|  '__ ||  .___.'(|  .-.  | ||  |OO ) \|  |)|  ||  |\    |
    |  | |  | |     |'|  |      |  | |  |(_'  '--'\  '  '-'  '|  | \   |
    ` + "`" + `--' ` + "`" + `--' ` + "`" + `-----' ` + "`" + `--'      ` + "`" + `--' ` + "`" + `--'   ` + "`" + `-----'   ` + "`" + `-----' ` + "`" + `--'  ` + "`" + `--'
    `
	fmt.Println(alpaconLogo)
}

var rootCmd = &cobra.Command{
	Use:   "alpacon",
	Short: "Alpacon CLI: Your Gateway to Alpacon Services",
	Long:  "Use this tool to interact with the alpacon service.",
	Run: func(cmd *cobra.Command, args []string) {
		showLogo()
		fmt.Println("Welcome to Alpacon CLI! Use 'alpacon [command]' to execute a specific command or 'alpacon help' to see all available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.CliError("While executing the command: %s", err)
	}
}

func init() {
	// version
	rootCmd.AddCommand(versionCmd)

	// login
	rootCmd.AddCommand(loginCmd)

	// iam
	rootCmd.AddCommand(iam.UserCmd)
	rootCmd.AddCommand(iam.GroupCmd)

	// server
	rootCmd.AddCommand(server.ServerCmd)

	// agent
	rootCmd.AddCommand(agent.AgentCmd)

	// websh
	rootCmd.AddCommand(websh.WebshCmd)

	// ftp
	rootCmd.AddCommand(ftp.UploadCmd)
	rootCmd.AddCommand(ftp.DownloadCmd)

	// packages
	rootCmd.AddCommand(packages.PackagesCmd)

	// log
	rootCmd.AddCommand(log.LogCmd)

}
