package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/config"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"io"
	"net/http"
	"strconv"
)

const (
	loginURL  = "/api/auth/login/"
	logoutURL = "/api/auth/logout/"
	tokenURL  = "/api/auth/tokens/"
	statusURL = "/api/status/"
)

func LoginAndSaveCredentials(loginReq *LoginRequest, token string) error {
	if token != "" {

		client := &client.AlpaconClient{
			HTTPClient: &http.Client{},
			BaseURL:    loginReq.ServerAddress,
			Token:      token,
			UserAgent:  utils.SetUserAgent(),
		}

		_, err := client.SendGetRequest(statusURL)
		if err != nil {
			return err
		}

		err = config.CreateConfig(loginReq.ServerAddress, token, "")
		if err != nil {
			return err
		}

		return nil
	}

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
	httpReq.Header.Set("User-Agent", utils.SetUserAgent())

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusFound {
		return fmt.Errorf("response status: %s", resp.Status)
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
		return err
	}

	return nil
}

func CreateAPIToken(ac *client.AlpaconClient, tokenRequest APITokenRequest) (string, error) {
	resp, err := ac.SendPostRequest(tokenURL, tokenRequest)
	if err != nil {
		return "", err
	}

	var response APITokenResponse
	if err = json.Unmarshal(resp, &response); err != nil {
		return "", err
	}

	return response.Key, nil
}

func GetAPITokenList(ac *client.AlpaconClient) ([]APITokenAttributes, error) {
	var tokenList []APITokenAttributes
	page := 1
	const pageSize = 100

	params := map[string]string{
		"page":      strconv.Itoa(page),
		"page_size": fmt.Sprintf("%d", pageSize),
	}
	for {
		responseBody, err := ac.SendGetRequest(utils.BuildURL(tokenURL, "", params))
		if err != nil {
			return nil, err
		}

		var response api.ListResponse[APITokenResponse]
		if err = json.Unmarshal(responseBody, &response); err != nil {
			return nil, err
		}

		for _, token := range response.Results {
			tokenList = append(tokenList, APITokenAttributes{
				ID:        token.ID,
				Name:      token.Name,
				Enabled:   token.Enabled,
				UpdatedAt: utils.TimeUtils(token.UpdatedAt),
				ExpiresAt: utils.TimeUtils(token.ExpiresAt),
			})
		}

		if len(response.Results) < pageSize {
			break
		}
		page++
	}
	return tokenList, nil
}

func GetAPITokenIDByName(ac *client.AlpaconClient, tokenName string) (string, error) {
	params := map[string]string{
		"name": tokenName,
	}
	body, err := ac.SendGetRequest(utils.BuildURL(tokenURL, "", params))
	if err != nil {
		return "", err
	}

	var response api.ListResponse[APITokenResponse]
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response.Count == 0 {
		return "", errors.New("no token found with the given name")
	}

	return response.Results[0].ID, nil
}

func DeleteAPIToken(ac *client.AlpaconClient, tokenID string) error {
	_, err := ac.SendDeleteRequest(utils.BuildURL(tokenURL, tokenID, nil))
	if err != nil {
		return err
	}

	return nil
}

func LogoutAndDeleteCredentials(ac *client.AlpaconClient) error {
	_, err := ac.SendPostRequest(logoutURL, nil)
	if err != nil {
		return err
	}

	err = config.DeleteConfig()
	if err != nil {
		return err
	}
	return nil
}
