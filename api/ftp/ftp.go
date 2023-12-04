package ftp

import (
	"bytes"
	"encoding/json"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"mime/multipart"
)

var (
	uploadAPIURL   = "/api/websh/uploads/"
	downloadAPIURL = "/api/websh/downloads/"
)

func UploadFile(ac *client.AlpaconClient, file string, serverName string, path string) error {
	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return err
	}

	content, err := utils.ReadFileFromPath(file)
	if err != nil {
		return err
	}

	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)

	err = multiPartWriter.WriteField("path", path)
	if err != nil {
		return err
	}

	fileWriter, err := multiPartWriter.CreateFormFile("content", file)
	if err != nil {
		return err
	}
	_, err = fileWriter.Write(content)
	if err != nil {
		return err
	}

	err = multiPartWriter.WriteField("server", serverID)
	if err != nil {
		return err
	}
	multiPartWriter.Close()

	err = ac.SendMultipartRequest(uploadAPIURL, multiPartWriter, requestBody)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFile(ac *client.AlpaconClient, serverName string, path string) (string, error) {
	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return "", err
	}

	downloadRequest := &DownloadRequest{
		Path:   path,
		Server: serverID,
	}

	body, err := ac.SendPostRequest(downloadAPIURL, downloadRequest)
	if err != nil {
		return "", err
	}

	var downloadResponse DownloadResponse
	err = json.Unmarshal(body, &downloadResponse)
	if err != nil {
		return "", err
	}

	return downloadResponse.DownloadURL, nil
}
