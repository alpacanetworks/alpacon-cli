package iam

import (
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/config"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type LoginRequest struct {
	ServerAddress string `json:"server_address"`
	Username      string `json:"username"`
	Password      string `json:"password"`
}

var (
	loginURL = "/iam/login/?next=/"
)

func getCSRFToken(serverURL string) (string, error) {
	resp, err := http.Get(serverURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "csrftoken" {
			return cookie.Value, nil
		}
	}

	return "", err
}

func LoginAndSaveCredentials(req *LoginRequest) error {
	serverAddress := req.ServerAddress

	// Before request login, get CSRF token from alpacon server
	csrfToken, err := getCSRFToken(serverAddress)
	if err != nil {
		return fmt.Errorf("csrf token not found %s", err)
	}

	data := url.Values{}
	data.Set("username", req.Username)
	data.Set("password", req.Password)

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	httpClient := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.AddCookie(&http.Cookie{Name: "csrftoken", Value: csrfToken})
			return nil
		},
	}

	// Log in to Alpacon server
	httpReq, err := http.NewRequest("POST", serverAddress+loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	httpReq.AddCookie(&http.Cookie{Name: "csrftoken", Value: csrfToken})
	httpReq.Header.Add("X-CSRFToken", csrfToken)
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = config.CreateConfig(serverAddress, httpClient)
	if err != nil {
		return fmt.Errorf("failed to create credential config file : %s", err)
	}

	return nil
}
