package event

import (
	"encoding/json"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/iam"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"net/url"
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

func buildURL(serverID, userID string, pageSize int) string {
	params := url.Values{}
	params.Add("server", serverID)
	params.Add("requested_by", userID)
	params.Add("page_size", fmt.Sprintf("%d", pageSize))
	return getEventURL + "?" + params.Encode()
}
