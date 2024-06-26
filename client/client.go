package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/config"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const (
	checkAuthURL       = "/api/auth/is_authenticated/"
	checkPrivilegesURL = "/api/iam/users/-"
)

func NewAlpaconAPIClient() (*AlpaconClient, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	client := &AlpaconClient{
		HTTPClient: &http.Client{},
		BaseURL:    config.ServerAddress,
		Token:      config.Token,
		UserAgent:  utils.SetUserAgent(),
	}

	err = client.checkAuth()
	if err != nil {
		return nil, err
	}

	err = client.checkPrivileges()
	if err != nil {
		return nil, err
	}

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

func (ac *AlpaconClient) checkPrivileges() error {
	body, err := ac.SendGetRequest(checkPrivilegesURL)
	if err != nil {
		return err
	}

	var checkPrivilegesResponse CheckPrivilegesResponse

	err = json.Unmarshal(body, &checkPrivilegesResponse)
	if err != nil {
		return err
	}

	ac.Privileges = getUserPrivileges(checkPrivilegesResponse.IsStaff, checkPrivilegesResponse.IsSuperuser)

	return nil
}

func getUserPrivileges(isStaff, isSuperuser bool) string {
	if isSuperuser {
		return "superuser"
	}
	if isStaff {
		return "staff"
	}
	return "general"
}

func (ac *AlpaconClient) SetWebsocketHeader() http.Header {
	headers := http.Header{}
	headers.Set("Origin", ac.BaseURL)
	headers.Set("User-Agent", ac.UserAgent)

	return headers
}

func (ac *AlpaconClient) setHTTPHeader(req *http.Request) *http.Request {
	req.Header.Set("User-Agent", ac.UserAgent)
	req.Header.Set("Authorization", fmt.Sprintf("token=\"%s\"", ac.Token))

	return req
}

func (ac *AlpaconClient) createRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, ac.BaseURL+url, body)
	if err != nil {
		return nil, err
	}

	req = ac.setHTTPHeader(req)
	if method == "POST" || method == "PATCH" || method == "PUT" {
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

	contentType := resp.Header.Get("Content-Type")
	// Check for non-empty and non-JSON content types. Empty content type allowed for responses without content (e.g., from PATCH requests).
	if contentType != "" && !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("Server error or incorrect request detected")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if req.Method == "POST" && (resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK) {
		return nil, errors.New(string(respBody))
	} else if req.Method == "DELETE" && resp.StatusCode != http.StatusNoContent {
		return nil, errors.New(string(respBody))
	} else if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.New(string(respBody))
	}

	return respBody, nil
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
func (ac *AlpaconClient) SendPostRequest(url string, body interface{}) ([]byte, error) {
	jsonValue, err := json.Marshal(body)
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

// TODO PUT
func (ac *AlpaconClient) SendPatchRequest(url string, body interface{}) ([]byte, error) {
	jsonValue, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := ac.createRequest("PATCH", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	return ac.sendRequest(req)
}

func (ac *AlpaconClient) SendMultipartRequest(url string, multiPartWriter *multipart.Writer, body bytes.Buffer) ([]byte, error) {
	req, err := ac.createRequest("POST", url, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	resp, err := ac.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	// Check for non-empty and non-JSON content types. Empty content type allowed for responses without content (e.g., from PATCH requests).
	if contentType != "" && !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("Server error or incorrect request detected")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(string(respBody))
	}

	return respBody, nil
}

// This function returns response for custom error handling in each function, unlike direct error throwing in sendRequest.
func (ac *AlpaconClient) SendGetRequestForDownload(url string) (*http.Response, error) {
	req, err := ac.createRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := ac.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ac *AlpaconClient) IsUsingHTTPS() (bool, error) {
	parsedURL, err := url.Parse(ac.BaseURL)
	if err != nil {
		return false, err
	}

	if parsedURL.Scheme == "https" {
		return true, nil
	}

	return false, nil
}
