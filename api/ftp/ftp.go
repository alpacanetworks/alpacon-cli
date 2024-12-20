package ftp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/event"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

const (
	uploadAPIURL   = "/api/webftp/uploads/"
	downloadAPIURL = "/api/webftp/downloads/"
)

func UploadFile(ac *client.AlpaconClient, src []string, dest, username, groupname string) ([]string, error) {
	serverName, remotePath := utils.SplitPath(dest)

	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, filePath := range src {
		file, err := utils.ReadFileFromPath(filePath)
		if err != nil {
			return nil, err
		}

		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)

		params := map[string]string{
			"path":      remotePath,
			"server":    serverID,
			"username":  username,
			"groupname": groupname,
		}
		for key, value := range params {
			if err := writer.WriteField(key, value); err != nil {
				return nil, err
			}
		}

		fileWriter, err := writer.CreateFormFile("content", filepath.Base(filePath))
		if err != nil {
			return nil, err
		}
		_, err = fileWriter.Write(file)
		if err != nil {
			return nil, err
		}
		_ = writer.Close()

		respBody, err := ac.SendMultipartRequest(uploadAPIURL, writer, requestBody)
		if err != nil {
			return nil, err
		}

		var response UploadResponse
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return nil, err
		}

		status, err := event.PollCommandExecution(ac, response.Command)
		if err != nil {
			return nil, err
		}
		if status.Status["text"] == "Stuck" || status.Status["text"] == "Error" {
			result = append(result, status.Status["message"].(string))
		} else {
			result = append(result, status.Result)
		}
	}

	return result, nil
}

func DownloadFile(ac *client.AlpaconClient, src, dest, username, groupname string) error {
	serverName, remotePathStr := utils.SplitPath(src)

	var remotePaths []string

	trimmedPathStr := strings.Trim(remotePathStr, "\"")
	remotePaths = strings.Fields(trimmedPathStr)

	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return err
	}

	for _, path := range remotePaths {
		downloadRequest := &DownloadRequest{
			Path:      path,
			Server:    serverID,
			Username:  username,
			Groupname: groupname,
		}

		postBody, err := ac.SendPostRequest(downloadAPIURL, downloadRequest)
		if err != nil {
			return err
		}

		var downloadResponse DownloadResponse
		err = json.Unmarshal(postBody, &downloadResponse)
		if err != nil {
			return err
		}

		status, err := event.PollCommandExecution(ac, downloadResponse.Command)
		if err != nil {
			return err
		}

		if status.Status["text"] == "Stuck" || status.Status["text"] == "Error" {
			utils.CliError("%s", status.Status["message"].(string))
		}
		if status.Status["text"] == "Failed" {
			utils.CliError("%s", status.Result)
		}
		utils.CliWarning("File Transfer Status: '%s'. Attempting to transfer '%s' from the Alpacon server. Note: Transfer may timeout after 100 seconds.", status.Result, path)

		maxAttempts := 100
		var resp *http.Response
		for count := 0; count < maxAttempts; count++ {
			resp, err = ac.SendGetRequestForDownload(utils.RemovePrefixBeforeAPI(downloadResponse.DownloadURL))
			if err != nil {
				return err
			}

			if resp.StatusCode == http.StatusOK {
				break
			} else {
				time.Sleep(time.Second * 1)
			}

			if count == maxAttempts-1 {
				return fmt.Errorf("%d attempts", maxAttempts)
			}
		}

		defer func() { _ = resp.Body.Close() }()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = utils.SaveFile(filepath.Join(dest, filepath.Base(path)), respBody)
		if err != nil {
			return err
		}
	}

	return nil
}
