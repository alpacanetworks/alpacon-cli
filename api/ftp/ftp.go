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
	uploadAPIURL   = "/api/websh/uploads/"
	downloadAPIURL = "/api/websh/downloads/"
)

func UploadFile(ac *client.AlpaconClient, src []string, dest, username, groupname string) ([]string, error) {
	serverName, remotePath := splitPath(dest)

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
		err = writer.Close()
		if err != nil {
			return nil, err
		}
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
		result = append(result, status)
	}

	return result, nil
}

func DownloadFile(ac *client.AlpaconClient, src, dest, username, groupname string) error {
	serverName, remotePathStr := splitPath(src)

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

		maxAttempts := 100
		var resp *http.Response

		status, err := event.PollCommandExecution(ac, downloadResponse.Command)
		if err != nil {
			return err
		}

		utils.CliWarning(fmt.Sprintf("File Transfer Status: '%s'. Attempting to transfer '%s' from the Alpacon server. Note: Transfer may timeout after 100 seconds.", status, path))

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

		defer resp.Body.Close()

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

func splitPath(path string) (string, string) {
	parts := strings.SplitN(path, ":", 2)
	return parts[0], parts[1]
}
