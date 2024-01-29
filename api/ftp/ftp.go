package ftp

import (
	"bytes"
	"encoding/json"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

const (
	uploadAPIURL   = "/api/websh/uploads/"
	downloadAPIURL = "/api/websh/downloads/"
)

func UploadFile(ac *client.AlpaconClient, src string, dest string) error {
	parts := strings.SplitN(dest, ":", 2)
	serverName := parts[0]
	remotePath := parts[1]

	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return err
	}

	content, err := utils.ReadFileFromPath(src)
	if err != nil {
		return err
	}

	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)

	err = multiPartWriter.WriteField("path", remotePath)
	if err != nil {
		return err
	}

	fileWriter, err := multiPartWriter.CreateFormFile("content", src)
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

func DownloadFile(ac *client.AlpaconClient, src string, dest string) (string, error) {
	parts := strings.SplitN(src, ":", 2)
	serverName := parts[0]
	remotePath := parts[1]

	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return "", err
	}

	downloadRequest := &DownloadRequest{
		Path:   parts[1],
		Server: serverID,
	}

	postBody, err := ac.SendPostRequest(downloadAPIURL, downloadRequest)
	if err != nil {
		return "", err
	}

	var downloadResponse DownloadResponse
	err = json.Unmarshal(postBody, &downloadResponse)
	if err != nil {
		return "", err
	}

	var data []byte
	maxAttempts := 100

	utils.CliWarning("Awaiting file transfer completion from Alpacon server. Transfer may timeout if it exceeds 100 seconds.")

	for count := 0; count < maxAttempts; count++ {
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}

	err = utils.SaveFile(filepath.Join(dest, filepath.Base(remotePath)), data)

	return downloadResponse.DownloadURL, nil
}

func splitPath(path string) (string, string) {
	parts := strings.SplitN(path, ":", 2)
	return parts[0], parts[1]
}
