package websh

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCommandParsing(t *testing.T) {
	tests := []struct {
		testName          string
		args              []string
		expectRoot        bool
		expectUsername    string
		expectGroupname   string
		expectServerName  string
		expectCommandArgs []string
	}{
		{
			testName:          "RootAccessToServer",
			args:              []string{"-r", "prod-server", "df", "-h"},
			expectRoot:        true,
			expectUsername:    "",
			expectGroupname:   "",
			expectServerName:  "prod-server",
			expectCommandArgs: []string{"df", "-h"},
		},
		{
			testName:          "ExecuteUpdateAsAdminSysadmin",
			args:              []string{"-u", "admin", "-g", "sysadmin", "update-server", "sudo", "apt-get", "update"},
			expectRoot:        false,
			expectUsername:    "admin",
			expectGroupname:   "sysadmin",
			expectServerName:  "update-server",
			expectCommandArgs: []string{"sudo", "apt-get", "update"},
		},
		{
			testName:          "DockerComposeDeploymentWithFlags",
			args:              []string{"deploy-server", "docker-compose", "-f", "/home/admin/deploy/docker-compose.yml", "up", "-d"},
			expectRoot:        false,
			expectUsername:    "",
			expectGroupname:   "",
			expectServerName:  "deploy-server",
			expectCommandArgs: []string{"docker-compose", "-f", "/home/admin/deploy/docker-compose.yml", "up", "-d"},
		},
		{
			testName:          "VerboseListInFileServer",
			args:              []string{"file-server", "ls", "-l", "/var/www"},
			expectRoot:        false,
			expectUsername:    "",
			expectGroupname:   "",
			expectServerName:  "file-server",
			expectCommandArgs: []string{"ls", "-l", "/var/www"},
		},
		{
			testName:          "MisplacedFlagOrderWithRoot",
			args:              []string{"-r", "df", "-h"},
			expectRoot:        true,
			expectUsername:    "",
			expectGroupname:   "",
			expectServerName:  "df",
			expectCommandArgs: []string(nil),
		},
		{
			testName:          "UnrecognizedFlagWithEchoCommand",
			args:              []string{"-x", "unknown-server", "echo", "Hello World"},
			expectRoot:        false,
			expectUsername:    "",
			expectGroupname:   "",
			expectServerName:  "-x",
			expectCommandArgs: []string{"unknown-server", "echo", "Hello World"},
		},
		{
			testName:          "AdminSysadminAccessToMultiFlagServer",
			args:              []string{"--username=admin", "--groupname=sysadmin", "multi-flag-server", "uptime"},
			expectRoot:        false,
			expectUsername:    "admin",
			expectGroupname:   "sysadmin",
			expectServerName:  "multi-flag-server",
			expectCommandArgs: []string{"uptime"},
		},
		{
			testName:          "CommandLineArgsResembleFlags",
			args:              []string{"--username", "admin", "server-name", "--fake-flag", "value"},
			expectRoot:        false,
			expectUsername:    "admin",
			expectGroupname:   "",
			expectServerName:  "server-name",
			expectCommandArgs: []string{"--fake-flag", "value"},
		},
		{
			testName:          "SysadminGroupWithMixedSyntax",
			args:              []string{"-g=sysadmin", "server-name", "echo", "hello world"},
			expectRoot:        false,
			expectUsername:    "",
			expectGroupname:   "sysadmin",
			expectServerName:  "server-name",
			expectCommandArgs: []string{"echo", "hello world"},
		},
		{
			testName:          "HelpRequestedViaCombinedFlags",
			args:              []string{"-rh"},
			expectRoot:        false,
			expectUsername:    "",
			expectGroupname:   "",
			expectServerName:  "-rh",
			expectCommandArgs: nil,
		},
		{
			testName:          "InvalidUsageDetected",
			args:              []string{"-u", "user", "-x", "unknown-flag", "server-name", "cmd"},
			expectRoot:        false,
			expectUsername:    "user",
			expectGroupname:   "",
			expectServerName:  "-x",
			expectCommandArgs: []string{"unknown-flag", "server-name", "cmd"},
		},
		{
			testName:          "ValidFlagsFollowedByInvalidFlag",
			args:              []string{"-u", "user", "-g", "group", "-x", "server-name", "cmd"},
			expectRoot:        false,
			expectUsername:    "user",
			expectGroupname:   "group",
			expectServerName:  "-x",
			expectCommandArgs: []string{"server-name", "cmd"},
		},
		{
			testName:          "FlagsIntermixedWithCommandArgs",
			args:              []string{"server-name", "-u", "user", "cmd", "-g", "group"},
			expectRoot:        false,
			expectUsername:    "user",
			expectGroupname:   "",
			expectServerName:  "server-name",
			expectCommandArgs: []string{"cmd", "-g", "group"},
		},
		{
			testName:          "FlagsAndCommandArgsIntertwined",
			args:              []string{"server-name", "-u", "user", "cmd", "-g", "group"},
			expectRoot:        false,
			expectUsername:    "user",
			expectGroupname:   "",
			expectServerName:  "server-name",
			expectCommandArgs: []string{"cmd", "-g", "group"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			root, username, groupname, serverName, commandArgs := executeTestCommand(tc.args)

			assert.Equal(t, tc.expectRoot, root, "Mismatch in root flag")
			assert.Equal(t, tc.expectUsername, username, "Mismatch in username")
			assert.Equal(t, tc.expectGroupname, groupname, "Mismatch in groupname")
			assert.Equal(t, tc.expectServerName, serverName, "Mismatch in server name")
			assert.Equal(t, tc.expectCommandArgs, commandArgs, "Mismatch in command arguments")

		})
	}
}

func executeTestCommand(args []string) (bool, string, string, string, []string) {
	var root bool
	var username, groupname, serverName string
	var commandArgs []string

	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "-r" || args[i] == "--root":
			root = true
		case args[i] == "-h" || args[i] == "--help":
			return root, username, groupname, serverName, commandArgs
		case strings.HasPrefix(args[i], "-u") || strings.HasPrefix(args[i], "--username"):
			username, i = extractValue(args, i)
		case strings.HasPrefix(args[i], "-g") || strings.HasPrefix(args[i], "--groupname"):
			groupname, i = extractValue(args, i)
		default:
			if serverName == "" {
				serverName = args[i]
			} else {
				commandArgs = append(commandArgs, args[i:]...)
				i = len(args)
			}
		}
	}

	return root, username, groupname, serverName, commandArgs
}
