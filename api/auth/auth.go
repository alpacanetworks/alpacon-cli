package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/config"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"io"
	"net/http"
)

const (
	loginURL      = "/api/auth/login/"
	tokenURL      = "/api/auth/tokens/"
	getTokenIDURL = "/api/auth/tokens/?=name"
	statusURL     = "/api/status/"
)

func LoginAndSaveCredentials(loginReq *LoginRequest, token string) error {
	if token != "" {

		client := &client.AlpaconClient{
			HTTPClient: &http.Client{},
			BaseURL:    loginReq.ServerAddress,
			Token:      token,
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

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusFound {
		return fmt.Errorf(resp.Status)
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

	for {
		params := utils.CreatePaginationParams(page, pageSize)
		responseBody, err := ac.SendGetRequest(tokenURL + "?" + params)
		if err != nil {
			return nil, err
		}

		var response APITokenListResponse
		if err = json.Unmarshal(responseBody, &response); err != nil {
			return nil, err
		}

		for _, token := range response.Results {
			tokenList = append(tokenList, APITokenAttributes{
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

func getAPITokenIDByName(ac *client.AlpaconClient, tokenName string) (string, error) {
	body, err := ac.SendGetRequest(getTokenIDURL + tokenName)
	if err != nil {
		return "", err
	}

	var response APITokenListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response.Count == 0 {
		return "", errors.New("no server found with the given name")
	}

	return response.Results[0].ID, nil
}

func DeleteAPIToken(ac *client.AlpaconClient, tokenName string) error {
	tokenID, err := getAPITokenIDByName(ac, tokenName)
	if err != nil {
		return err
	}

	_, err = ac.SendDeleteRequest(tokenURL + tokenID + "/")
	if err != nil {
		return err
	}

	return nil
}
