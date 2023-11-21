package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
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
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// login
	rootCmd.AddCommand(loginCmd)
}
