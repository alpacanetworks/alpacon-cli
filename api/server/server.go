package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
)

const (
	serverURL      = "/api/servers/servers/"
	getServerIDURL = "/api/servers/servers/?name="
)

func GetServerList(ac *client.AlpaconClient) ([]ServerAttributes, error) {
	var serverList []ServerAttributes
	page := 1
	const pageSize = 100

	for {
		params := utils.CreatePaginationParams(page, pageSize)
		responseBody, err := ac.SendGetRequest(serverURL + "?" + params)
		if err != nil {
			return nil, err
		}

		var response ServerListResponse
		if err = json.Unmarshal(responseBody, &response); err != nil {
			return nil, err
		}

		for _, server := range response.Results {
			serverList = append(serverList, ServerAttributes{
				Name:      server.Name,
				IP:        server.RemoteIP,
				OS:        fmt.Sprintf("%s %s", server.OSName, server.OSVersion),
				Connected: server.IsConnected,
				Owner:     server.OwnerName,
			})
		}

		if len(response.Results) < pageSize {
			break
		}
		page++
	}

	return serverList, nil
}

func GetServerDetail(ac *client.AlpaconClient, serverName string) ([]byte, error) {
	serverID, err := GetServerIDByName(ac, serverName)
	if err != nil {
		return nil, err
	}

	body, err := ac.SendGetRequest(serverURL + serverID)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func DeleteServer(ac *client.AlpaconClient, serverName string) error {
	serverID, err := GetServerIDByName(ac, serverName)
	if err != nil {
		return err
	}

	_, err = ac.SendDeleteRequest(serverURL + serverID + "/")
	if err != nil {
		return err
	}

	return nil
}

func GetServerIDByName(ac *client.AlpaconClient, serverName string) (string, error) {
	body, err := ac.SendGetRequest(getServerIDURL + serverName)
	if err != nil {
		return "", err
	}

	var response ServerListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response.Count == 0 {
		return "", errors.New("no server found with the given name")
	}

	return response.Results[0].ID, nil
}

func GetServerNameByID(ac *client.AlpaconClient, serverID string) (string, error) {
	body, err := ac.SendGetRequest(serverURL + serverID)
	if err != nil {
		return "", err
	}

	var response ServerDetails
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Name, nil
}

func CreateServer(ac *client.AlpaconClient, serverRequest ServerRequest) (ServerCreatedResponse, error) {
	var response ServerCreatedResponse

	responseBody, err := ac.SendPostRequest(serverURL, serverRequest)
	if err != nil {
		return ServerCreatedResponse{}, err
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return ServerCreatedResponse{}, err
	}

	return response, nil
}
