package ftp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/event"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	connectAPIURL  = "/api/webftp/sessions/"
	uploadAPIURL   = "/api/webftp/uploads/"
	downloadAPIURL = "/api/webftp/downloads/"
)

func connectWebSocket(ac *client.AlpaconClient, serverId, username string) (*websocket.Conn, string, error) {
	reqBody := &ConnectRequest{
		Server:   serverId,
		Username: username,
	}

	respBody, err := ac.SendPostRequest(connectAPIURL, reqBody)
	if err != nil {
		return nil, "", err
	}

	var cr connectResponse
	if err := json.Unmarshal(respBody, &cr); err != nil {
		return nil, "", fmt.Errorf("failed to parse websocket_url: %w", err)
	}

	wsConn, _, err := websocket.DefaultDialer.Dial(cr.WebsocketURL, ac.SetWebsocketHeader())
	if err != nil {
		return nil, "", fmt.Errorf("dial error: %w", err)
	}

	return wsConn, cr.UserChannel, nil
}

func ManageWebSocketSession(ac *client.AlpaconClient, serverID, username string) (*websocket.Conn, string, error) {

	// Establish the raw WebSocket connection
	wsConn, channel, err := connectWebSocket(ac, serverID, username)
	if err != nil {
		return nil, "", fmt.Errorf("failed to connect websocket: %w", err)
	}

	// Read initial handshake message to establish session
	if _, _, err := wsConn.ReadMessage(); err != nil {
		wsConn.Close()
		return nil, "", fmt.Errorf("handshake failed: %w", err)
	}

	return wsConn, channel, nil
}

func FinalizeWebSocketSession(wsConn *websocket.Conn) ([]byte, error) {
	// Read the final message containing status
	_, msg, err := wsConn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("failed to read final WS message: %w", err)
	}

	// Close the connection
	if err := wsConn.Close(); err != nil {
		return nil, fmt.Errorf("failed to close WS connection: %w", err)
	}

	return msg, nil
}

func uploadToS3(uploadUrl string, file io.Reader) error {
	req, err := http.NewRequest(http.MethodPut, uploadUrl, file)
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

	var wsConn *websocket.Conn
	var channel string
	wsConn, channel, err = ManageWebSocketSession(ac, serverID, username)
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

		uploadRequest := &UploadRequest{
			Id:             uuid.New().String(),
			Name:           filepath.Base(filePath),
			Path:           remotePath,
			Server:         serverID,
			Username:       username,
			Groupname:      groupname,
			Channel:        channel,
			AllowOverwrite: "true",
		}

		respBody, err := ac.SendPostRequest(uploadAPIURL, uploadRequest)
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

		fullURL := utils.BuildURL(uploadAPIURL, path.Join(response.Id, "upload"), nil)
		_, err = ac.SendGetRequest(fullURL)
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

	msg, err := FinalizeWebSocketSession(wsConn)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(msg, &statMsg); err != nil {
		return nil, err
	}

	if !statMsg.Data.Success {
		return nil, fmt.Errorf("%s", statMsg.Data.Message)
	}

	return result, nil
}

func UploadFolder(ac *client.AlpaconClient, src []string, dest, username, groupname string) ([]string, error) {
	serverName, remotePath := utils.SplitPath(dest)

	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return nil, err
	}

	var wsConn *websocket.Conn
	var channel string
	wsConn, channel, err = ManageWebSocketSession(ac, serverID, username)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, folderPath := range src {
		zipBytes, err := utils.Zip(folderPath)
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

		uploadRequest := &UploadRequest{
			Id:             uuid.New().String(),
			AllowUnzip:     "true",
			AllowOverwrite: "true",
			Name:           zipName,
			Path:           remotePath,
			Server:         serverID,
			Username:       username,
			Groupname:      groupname,
			Channel:        channel,
		}

		respBody, err := ac.SendPostRequest(uploadAPIURL, uploadRequest)
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

	msg, err := FinalizeWebSocketSession(wsConn)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(msg, &statMsg); err != nil {
		return nil, err
	}

	if !statMsg.Data.Success {
		return nil, fmt.Errorf("%s", statMsg.Data.Message)
	}

	return result, nil
}

func DownloadFile(ac *client.AlpaconClient, src, dest, username, groupname string, recursive bool) error {
	serverName, remotePathStr := utils.SplitPath(src)

	var remotePaths []string
	var resourceType string

	trimmedPathStr := strings.Trim(remotePathStr, "\"")
	remotePaths = strings.Fields(trimmedPathStr)

	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return err
	}

	if recursive {
		resourceType = "folder"
	} else {
		resourceType = "file"
	}

	for _, path := range remotePaths {
		downloadRequest := &DownloadRequest{
			Path:         path,
			Name:         filepath.Base(path),
			Server:       serverID,
			Username:     username,
			Groupname:    groupname,
			ResourceType: resourceType,
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

		var fileName string
		if recursive {
			fileName = filepath.Base(path) + ".zip"
		} else {
			fileName = filepath.Base(path)
		}
		err = utils.SaveFile(filepath.Join(dest, fileName), respBody)
		if err != nil {
			return err
		}
		err = utils.Unzip(filepath.Join(dest, fileName), dest)
		if err != nil {
			return err
		}
		err = utils.DeleteFile(filepath.Join(dest, fileName))
		if err != nil {
			return err
		}
	}

	return nil
}
