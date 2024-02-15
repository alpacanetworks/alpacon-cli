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

	commandRequest := &Command{
		Shell:     "system", // TODO Support osquery, alpamon
		Line:      command,
		Username:  username,
		Groupname: groupname,
		Server:    serverID,
		RunAfter:  []string{},
	}
	_, err = ac.SendPostRequest(getEventURL, commandRequest)
	if err != nil {
		return "", err
	}

	timer := time.NewTimer(5 * time.Minute)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timer.C:
			return "", errors.New("command execution timed out")
		case <-tick:
			responseBody, err := ac.SendGetRequest(buildPageURL(1, 1))
			if err != nil {
				continue
			}
			var response EventListResponse
			if err = json.Unmarshal(responseBody, &response); err != nil {
				return "", err
			}
			if len(response.Results) == 0 || utils.BoolPointerToString(response.Results[0].Success) == "null" {
				continue
			}
			if response.Results[0].Status["text"] == "Acked" {
				timer.Reset(5 * time.Minute)
				continue
			}
			if response.Results[0].Status["text"] == "Stuck" || response.Results[0].Status["text"] == "Error" {
				return response.Results[0].Status["message"].(string), nil
			}
			return response.Results[0].Result, nil
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

func buildPageURL(page int, pageSize int) string {
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("page_size", fmt.Sprintf("%d", pageSize))
	return getEventURL + "?" + params.Encode()
}
