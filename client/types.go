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
	ID            string `json:"id"`
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Tags          string `json:"tags"`
	Description   string `json:"description"`
	NumGroups     int    `json:"num_groups"`
	UID           int    `json:"uid"`
	Shell         string `json:"shell"`
	HomeDirectory string `json:"home_directory"`
	IsActive      bool   `json:"is_active"`
	IsStaff       bool   `json:"is_staff"`
	IsSuperuser   bool   `json:"is_superuser"`
	IsLDAPUser    bool   `json:"is_ldap_user"`
	DateJoined    string `json:"date_joined"`
	LastLogin     string `json:"last_login"`
	LastLoginIP   string `json:"last_login_ip"`
	AddedAt       string `json:"added_at"`
	UpdatedAt     string `json:"updated_at"`
}
