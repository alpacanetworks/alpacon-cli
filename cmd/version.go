package cmd

import (
	"context"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
)

var VersionCli = "v0.0.2"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the current CLI version.",
	Long:  "Displays the current version of the CLI and checks if there is an available update.",
	Run: func(cmd *cobra.Command, args []string) {
		release, skip := versionCheck()
		if skip == false {
			utils.CliWarning("Upgrade available. Current version: %s. Latest version: %s \n"+
				"Visit %s for update instructions and release notes.", VersionCli, release.GetTagName(), release.GetHTMLURL())
			return
		}

		utils.CliInfo("Current version (%s) is the latest.", release.GetTagName())
		return
	},
}

func versionCheck() (*github.RepositoryRelease, bool) {
	client := github.NewClient(nil)
	ctx := context.Background()

	release, _, err := client.Repositories.GetLatestRelease(ctx, "alpacanetworks", "alpacon-cli")
	if err != nil {
		utils.CliError("Checking for a newer version failed with %s \n", err)
		return nil, true
	}

	if release.GetTagName() != VersionCli {
		return release, false
	}

	return release, true
}
