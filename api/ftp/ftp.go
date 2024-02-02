package ftp

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func UploadFile(ac *client.AlpaconClient, src []string, dest string) error {
	serverName, remotePath := splitPath(dest)

	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return err
	}

	for _, path := range src {
		var requestBody bytes.Buffer
		multiPartWriter := multipart.NewWriter(&requestBody)
		defer multiPartWriter.Close()

		if err = multiPartWriter.WriteField("path", remotePath); err != nil {
			return err
		}
		if err = multiPartWriter.WriteField("server", serverID); err != nil {
			return err
		}

		content, err := utils.ReadFileFromPath(path)
		if err != nil {
			return err
		}

		fileWriter, err := multiPartWriter.CreateFormFile("content", filepath.Base(path))
		if err != nil {
			return err
		}
		if _, err = fileWriter.Write(content); err != nil {
			return err
		}

		if err = multiPartWriter.Close(); err != nil {
			return err
		}

		if err = ac.SendMultipartRequest(uploadAPIURL, multiPartWriter, requestBody); err != nil {
			return err
		}
	}

	return nil
}

func DownloadFile(ac *client.AlpaconClient, src string, dest string) error {
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
			Path:   path,
			Server: serverID,
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

		utils.CliWarning(fmt.Sprintf("Awaits %s file transfer completion from Alpacon server. Transfer may timeout after 100 seconds. If the specified file is not found, it will not download even after 100 seconds.", path))

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
