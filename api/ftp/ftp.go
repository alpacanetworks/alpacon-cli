package ftp

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/event"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	uploadAPIURL   = "/api/webftp/uploads/"
	downloadAPIURL = "/api/webftp/downloads/"
)

func uploadToS3(uploadUrl string, file io.Reader) error {
	req, err := http.NewRequest("PUT", uploadUrl, file)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return err
	}
	return nil
}

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

		fileWriter, err := writer.CreateFormFile("content", filepath.Base(filePath))
		if err != nil {
			return nil, err
		}
		_, err = fileWriter.Write(file)
		if err != nil {
			return nil, err
		}
		_ = writer.Close()

		body := map[string]string{
			"id":        uuid.New().String(),
			"name":      filepath.Base(filePath),
			"path":      remotePath,
			"server":    serverID,
			"username":  username,
			"groupname": groupname,
		}

		respBody, err := ac.SendPostRequest(uploadAPIURL, body)
		if err != nil {
			return nil, err
		}

		var response UploadResponse
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return nil, err
		}

		if response.UploadUrl != "" {
			err = uploadToS3(response.UploadUrl, bytes.NewReader(file))
			if err != nil {
				return nil, err
			}
		}

		_, err = ac.SendGetRequest(uploadAPIURL + response.Id + "/upload")
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

func UploadFolder(ac *client.AlpaconClient, src []string, dest, username, groupname string) ([]string, error) {
	serverName, remotePath := utils.SplitPath(dest)

	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, folderPath := range src {
		zipBytes, err := compressFolder(folderPath)
		if err != nil {
			return nil, err
		}
		zipName := filepath.Base(folderPath) + ".zip"

		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)

		fileWriter, err := writer.CreateFormFile("content", zipName)
		if err != nil {
			return nil, err
		}
		_, err = fileWriter.Write(zipBytes)
		if err != nil {
			return nil, err
		}
		_ = writer.Close()

		body := map[string]string{
			"id":          uuid.New().String(),
			"allow_unzip": "true",
			"name":        zipName,
			"path":        remotePath,
			"server":      serverID,
			"username":    username,
			"groupname":   groupname,
		}

		respBody, err := ac.SendPostRequest(uploadAPIURL, body)
		if err != nil {
			return nil, err
		}

		var response UploadResponse
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return nil, err
		}

		if response.UploadUrl != "" {
			err = uploadToS3(response.UploadUrl, bytes.NewReader(zipBytes))
			if err != nil {
				return nil, err
			}
		}

		_, err = ac.SendGetRequest(uploadAPIURL + response.Id + "/upload")
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
			Path:         path,
			Name:         filepath.Base(path),
			Server:       serverID,
			Username:     username,
			Groupname:    groupname,
			ResourceType: "file",
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
			resp, err = http.Get(downloadResponse.DownloadURL)
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

func compressFolder(folderPath string) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(folderPath, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			relPath += "/"
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		if !info.IsDir() {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		zipWriter.Close()
		return nil, err
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
