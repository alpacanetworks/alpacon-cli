package client

import "net/http"

type AlpaconClient struct {
	HTTPClient *http.Client
	BaseURL    string
	Token      string
	Privileges string
}

type CheckAuthResponse struct {
	Authenticated bool `json:"authenticated"`
}

type CheckPrivilegesResponse struct {
	IsStaff     bool `json:"is_staff"`
	IsSuperuser bool `json:"is_superuser"`
}
