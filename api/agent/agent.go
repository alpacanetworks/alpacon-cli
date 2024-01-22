package agent

import (
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"net/url"
)

const (
	baseURL        = "/api/servers/servers"
	upgradeAction  = "upgrade_agent"
	restartAction  = "restart_agent"
	shutdownAction = "shutdown_agent"
)

var actionMap = map[string]string{
	"upgrade":  upgradeAction,
	"restart":  restartAction,
	"shutdown": shutdownAction,
}

func RequestAgentAction(ac *client.AlpaconClient, serverName string, action string) error {
	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return err
	}

	url, err := buildURL(serverID)
	if err != nil {
		return err
	}

	actionValue, ok := actionMap[action]
	if !ok {
		return fmt.Errorf("invalid action: %s. Valid actions are: upgrade, restart, shutdown", action)
	}

	request := RequestAgent{
		Action: actionValue,
	}

	_, err = ac.SendPostRequest(url, request)
	return err
}

func buildURL(serverID string) (string, error) {
	url, err := url.JoinPath(baseURL, serverID, "/actions/")
	if err != nil {
		return "", err
	}

	return url, nil
}
