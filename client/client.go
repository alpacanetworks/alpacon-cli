package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/config"
	"io"
	"net/http"
)

type AlpaconClient struct {
	HTTPClient *http.Client
	BaseURL    string
	UserAgent  string
	CSRFToken  string
	SessionID  string
}

func NewAlpaconAPIClient() (*AlpaconClient, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	client := &AlpaconClient{
		BaseURL:    config.ServerAddress,
		CSRFToken:  config.CSRFToken,
		SessionID:  config.SessionID,
		HTTPClient: &http.Client{},
	}

	// TODO version check
	// res, err := utils.VersionCheck()

	return client, nil
}

func (ac *AlpaconClient) createRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, ac.BaseURL+url, body)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: ac.CSRFToken})
	req.AddCookie(&http.Cookie{Name: "sessionid", Value: ac.SessionID})
	req.Header.Add("X-CSRFToken", ac.CSRFToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (ac *AlpaconClient) sendRequest(req *http.Request) ([]byte, error) {
	resp, err := ac.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("server error: %d %s", resp.StatusCode, resp.Status)
	}

	return body, nil
}

// Get Request to Alpacon Server
func (ac *AlpaconClient) SendGetRequest(url string) ([]byte, error) {
	req, err := ac.createRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return ac.sendRequest(req)
}

// POST Request to Alpacon Server
func (ac *AlpaconClient) SendPostRequest(url string, params interface{}) ([]byte, error) {
	jsonValue, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := ac.createRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	return ac.sendRequest(req)
}

// TODO DELETE, PUT, PATCH
