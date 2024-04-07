package utils

import (
	"fmt"
	"os"
)

var (
	gitIssueURL = "https://github.com/alpacanetworks/alpacon-cli/issues"
)

// TODO VersionCheck()

func reportCLIError() {
	fmt.Println("For issues, check the latest version or report on", gitIssueURL)
}

// CliError handles all error messages in the CLI.
func CliError(msg string, args ...interface{}) {
	errorMessage := fmt.Sprintf(msg, args...)
	fmt.Fprintf(os.Stderr, "%s: %s\n", Red("Error"), errorMessage)
	reportCLIError()
	os.Exit(1)
}

// CliInfo handles all informational messages in the CLI.
func CliInfo(msg string, args ...interface{}) {
	infoMessage := fmt.Sprintf(msg, args...)
	fmt.Fprintf(os.Stderr, "%s: %s\n", Blue("Info"), infoMessage)
}

// CliWarning handles all warning messages in the CLI.
func CliWarning(msg string, args ...interface{}) {
	warningMessage := fmt.Sprintf(msg, args...)
	fmt.Fprintf(os.Stderr, "%s: %s\n", Yellow("Warning"), warningMessage)
}

// CliInfoWithExit prints an informational message to stderr and exits the program with a status code of 0
func CliInfoWithExit(msg string, args ...interface{}) {
	infoMessage := fmt.Sprintf(msg, args...)
	fmt.Fprintf(os.Stderr, "%s: %s\n", Blue("Info"), infoMessage)
	os.Exit(0) // Use exit code 0 to indicate successful completion.
}
