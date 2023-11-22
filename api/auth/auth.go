package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/config"
	"io"
	"net/http"
)

type LoginRequest struct {
	ServerAddress string `json:"server_address"`
	Username      string `json:"username"`
	Password      string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

var (
	loginURL = "/api/auth/login/"
)

func LoginAndSaveCredentials(loginReq *LoginRequest) error {
	serverAddress := loginReq.ServerAddress

	reqBody, err := json.Marshal(loginReq)
	if err != nil {
		return err
	}

	httpClient := &http.Client{}

	// Log in to Alpacon server
	httpReq, err := http.NewRequest("POST", serverAddress+loginURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	httpReq.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("server error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var loginResponse LoginResponse
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		return err
	}

	err = config.CreateConfig(serverAddress, loginResponse.Token, loginResponse.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to create credential config file : %s", err)
	}

	return nil
}
