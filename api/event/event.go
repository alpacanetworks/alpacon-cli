package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/iam"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"net/url"
	"time"
)

const (
	getEventURL = "/api/events/commands/"
)

func GetEventList(ac *client.AlpaconClient, pageSize int, serverName string, userName string) ([]EventAttributes, error) {
	var serverID, userID string
	var err error
	if serverName != "" {
		serverID, err = server.GetServerIDByName(ac, serverName)
		if err != nil {
			return nil, err
		}
	}
	if userName != "" {
		userID, err = iam.GetUserIDByName(ac, userName)
		if err != nil {
			return nil, err
		}
	}
	responseBody, err := ac.SendGetRequest(buildURL(serverID, userID, pageSize))
	if err != nil {
		return nil, err
	}

	var response EventListResponse
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return nil, err
	}

	var eventList []EventAttributes
	for _, event := range response.Results {
		eventList = append(eventList, EventAttributes{
			Server:      event.ServerName,
			Shell:       event.Shell,
			Command:     event.Line,
			Result:      utils.TruncateString(event.Result, 70),
			Status:      utils.BoolPointerToString(event.Success),
			Operator:    event.RequestedByName,
			RequestedAt: utils.TimeUtils(event.AddedAt),
		})
	}
	return eventList, nil
}

func RunCommand(ac *client.AlpaconClient, serverName, command string, username, groupname string) (string, error) {
	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return "", err
	}

	commandRequest := &CommandRequest{
		Shell:     "system", // TODO Support osquery, alpamon
		Line:      command,
		Username:  username,
		Groupname: groupname,
		Server:    serverID,
		RunAfter:  []string{},
	}
	respBody, err := ac.SendPostRequest(getEventURL, commandRequest)
	if err != nil {
		return "", err
	}

	var cmdResponse CommandResponse

	err = json.Unmarshal(respBody, &cmdResponse)
	if err != nil {
		return "", err
	}

	return PollCommandExecution(ac, cmdResponse.Id)
}

func PollCommandExecution(ac *client.AlpaconClient, cmdId string) (string, error) {
	var response EventDetails

	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timer.C:
			return "", errors.New("command execution timed out")
		case <-ticker.C:
			responseBody, err := ac.SendGetRequest(getEventURL + cmdId)
			if err != nil {
				continue
			}
			if err = json.Unmarshal(responseBody, &response); err != nil {
				return "", err
			}

			switch response.Status["text"] {
			case "Acked":
				continue
			case "Stuck", "Error":
				return response.Status["message"].(string), nil
			default:
				return response.Result, nil
			}
		}
	}
}

func buildURL(serverID, userID string, pageSize int) string {
	params := url.Values{}
	params.Add("server", serverID)
	params.Add("requested_by", userID)
	params.Add("page_size", fmt.Sprintf("%d", pageSize))
	return getEventURL + "?" + params.Encode()
}
