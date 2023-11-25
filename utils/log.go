package utils

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/savioxavier/termlink"
	"os"
)

var (
	gitIssueURL = "https://github.com/alpacanetworks/alpacon-cli/issues"
)

// TODO VersionCheck()

func reportCLIError() {
	gitIssueLink := termlink.ColorLink("", gitIssueURL, "blue")
	fmt.Println("For issues, check the latest version or report on", gitIssueLink)
}

// CliError handles all error messages in the CLI.
func CliError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", color.Red.Sprintf("Error"), fmt.Sprintf(msg, args...))
	reportCLIError()
	os.Exit(1)
}

// CliInfo handles all informational messages in the CLI.
func CliInfo(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", color.Blue.Sprintf("Info"), fmt.Sprintf(msg, args...))
}

// CliWarning handles all warning messages in the CLI.
func CliWarning(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", color.Yellow.Sprintf("Warning"), fmt.Sprintf(msg, args...))
}
