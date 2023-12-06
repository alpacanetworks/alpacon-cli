package agent

import (
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
)

var (
	baseURL        = "/api/servers/servers/"
	upgradeAction  = "upgrade_agent"
	restartAction  = "restart_agent"
	shutdownAction = "shutdown_agent"
)

func RequestAgentAction(ac *client.AlpaconClient, serverName string, action string) error {
	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return err
	}

	url := buildURL(serverID)
	actionMap := map[string]string{
		"upgrade":  upgradeAction,
		"restart":  restartAction,
		"shutdown": shutdownAction,
	}

	actionValue, ok := actionMap[action]
	if !ok {
		return fmt.Errorf("invalid action: %s", action)
	}

	request := RequestAgent{
		Action: actionValue,
	}

	_, err = ac.SendPostRequest(url, request)
	if err != nil {
		return err
	}

	return nil
}

func buildURL(serverID string) string {
	return baseURL + serverID + "/actions/"
}
