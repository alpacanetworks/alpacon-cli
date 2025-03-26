package auth0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alpacanetworks/alpacon-cli/utils"
)

var path = struct {
	env        string
	deviceCode string
	prefetch   string
	token      string
}{
	env:        "/api/auth/env/?client=cli",
	deviceCode: "/oauth/device/code",
	prefetch:   "/api/workspaces/workspaces/-/prefetch",
	token:      "/oauth/token",
}

func FetchAuthEnv(workspaceURL string) (*AuthEnvResponse, error) {
	apiURL := workspaceURL + path.env

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", utils.GetUserAgent())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var env AuthEnvResponse
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &env, nil
}

func RequestDeviceCode(workspaceURL string, envInfo *AuthEnvResponse) (*DeviceCodeResponse, error) {
	apiURL := "https://" + envInfo.Domain + path.deviceCode

	subDomain, err := extractSubdomain(workspaceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract schema name: %w", err)
	}

	data := map[string]string{
		"client_id": envInfo.ClientID,
		"scope":     fmt.Sprintf("openid profile email offline_access cli org:%s", subDomain),
		"audience":  envInfo.Audience,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request data: %w", err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned error: %d", resp.StatusCode)
	}

	var deviceCode DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceCode); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("Received Device Code Response: %+v\n", deviceCode)

	return &deviceCode, nil
}

func PollForToken(deviceCodeRes *DeviceCodeResponse, envInfo *AuthEnvResponse) (*TokenResponse, error) {
	startTime := time.Now()

	for {
		if time.Since(startTime).Seconds() > float64(deviceCodeRes.ExpiresIn) {
			return nil, fmt.Errorf("authentication timed out. Please restart the login process")
		}

		tokenResponse, err := requestAccessToken(deviceCodeRes.DeviceCode, envInfo)
		if err != nil {
			if strings.Contains(err.Error(), "authorization_pending") {
				fmt.Println("Waiting for user to complete authentication...")
				time.Sleep(time.Duration(deviceCodeRes.Interval) * time.Second)
				continue
			}
			return nil, err
		}

		return tokenResponse, nil
	}
}

func RefreshAccessToken(workspaceURL string, refreshToken string) (*TokenResponse, error) {
	envInfo, err := FetchAuthEnv(workspaceURL)
	if err != nil {
		utils.CliError("Failed to fetch auth env: %v", err)
	}

	apiURL := "https://" + envInfo.Domain + path.token

	subDomain, err := extractSubdomain(workspaceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract schema name: %w", err)
	}

	data := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     envInfo.ClientID,
		"refresh_token": refreshToken,
		"scope":         fmt.Sprintf("cli org:%s", subDomain),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode refresh token request: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send refresh request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read refresh response: %w", err)
	}

	var tokenRes TokenResponse
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return nil, fmt.Errorf("failed to parse refresh response JSON: %w", err)
	}

	if tokenRes.Error != "" {
		return nil, fmt.Errorf("Auth0 error: %s - %s", tokenRes.Error, tokenRes.ErrorDesc)
	}

	return &tokenRes, nil
}

func requestAccessToken(deviceCode string, envInfo *AuthEnvResponse) (*TokenResponse, error) {

	apiURL := "https://" + envInfo.Domain + path.token

	data := map[string]string{
		"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
		"device_code": deviceCode,
		"client_id":   envInfo.ClientID,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode JSON: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var tokenResp TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	if tokenResp.Error != "" {
		return nil, fmt.Errorf("Auth0 error: %s - %s", tokenResp.Error, tokenResp.ErrorDesc)
	}

	return &tokenResp, nil
}

func extractSubdomain(workspaceURL string) (string, error) {
	parsedURL, err := url.Parse(workspaceURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}

	parts := strings.Split(parsedURL.Host, ".")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid workspace URL: Subdomain is required")
	}

	return parts[0], nil
}
