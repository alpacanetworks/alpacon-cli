package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/config"
	"io"
	"mime/multipart"
	"net/http"
)

var (
	checkAuthURL = "/api/auth/is_authenticated/"
)

type AlpaconClient struct {
	HTTPClient *http.Client
	BaseURL    string
	Token      string
}

type CheckAuthResponse struct {
	Authenticated bool `json:"authenticated"`
}

func NewAlpaconAPIClient() (*AlpaconClient, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	client := &AlpaconClient{
		HTTPClient: &http.Client{},
		BaseURL:    config.ServerAddress,
		Token:      config.Token,
	}

	err = client.checkAuth()
	if err != nil {
		return nil, err
	}

	// TODO CLI version check
	// res, err := utils.VersionCheck()

	return client, nil
}

func (ac *AlpaconClient) checkAuth() error {
	body, err := ac.SendGetRequest(checkAuthURL)
	if err != nil {
		return err
	}

	var checkAuthResponse CheckAuthResponse

	err = json.Unmarshal(body, &checkAuthResponse)
	if err != nil {
		return err
	}
	if !checkAuthResponse.Authenticated {
		return errors.New("authenticated failed")
	}

	return nil
}

func (ac *AlpaconClient) createRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, ac.BaseURL+url, body)
	if err != nil {
		return nil, err
	}

	authHeaderValue := fmt.Sprintf("token=\"%s\"", ac.Token)
	req.Header.Add("Authorization", authHeaderValue)
	if method == "POST" {
		req.Header.Add("Content-Type", "application/json")
	}

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

	// Alpacon API Server always return 201 on POST request
	if req.Method == "POST" && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("POST request failed: expected status 201, got %d (%s)", resp.StatusCode, resp.Status)
	} else if req.Method == "DELETE" && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("DELETE request failed: expected status 204, got %d (%s). Verify permissions.", resp.StatusCode, resp.Status)
	} else if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected status code: %d (%s)", resp.StatusCode, resp.Status)
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

func (ac *AlpaconClient) SendDeleteRequest(url string) ([]byte, error) {
	req, err := ac.createRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return ac.sendRequest(req)
}

// TODO PUT, PATCH

func (ac *AlpaconClient) SendMultipartRequest(url string, multiPartWriter *multipart.Writer, body bytes.Buffer) error {
	req, err := ac.createRequest("POST", url, &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	resp, err := ac.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New(resp.Status + string(respBody))
	}

	return nil
}
