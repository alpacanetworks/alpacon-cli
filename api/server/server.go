package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/client"
)

var (
	getServerURL   = "/api/servers/servers/"
	getServerIDURL = "/api/servers/servers/?name="
)

func GetServerList(ac *client.AlpaconClient) ([]ServerAttributes, error) {
	responseBody, err := ac.SendGetRequest(getServerURL)
	if err != nil {
		return nil, err
	}

	var response ServerListResponse
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return nil, err
	}

	var serverList []ServerAttributes
	for _, server := range response.Results {
		serverList = append(serverList, ServerAttributes{
			Name:      server.Name,
			IP:        server.RemoteIP,
			OS:        fmt.Sprintf("%s %s", server.OSName, server.OSVersion),
			Connected: server.IsConnected,
			Owner:     server.OwnerName,
		})
	}

	return serverList, nil
}

func GetServerDetail(ac *client.AlpaconClient, serverName string) ([]byte, error) {
	serverID, err := GetServerIDByName(ac, serverName)
	if err != nil {
		return nil, err
	}

	body, err := ac.SendGetRequest(getServerURL + serverID)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetServerIDByName(ac *client.AlpaconClient, name string) (string, error) {
	body, err := ac.SendGetRequest(getServerIDURL + name)
	if err != nil {
		return "", err
	}

	var serverListResponse ServerListResponse
	err = json.Unmarshal(body, &serverListResponse)
	if err != nil {
		return "", err
	}

	if serverListResponse.Count == 0 {
		return "", errors.New("no server found with the given name")
	}

	return serverListResponse.Results[0].ID, nil
}
