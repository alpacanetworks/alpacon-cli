package ftp

import "time"

type DownloadRequest struct {
	Path         string `json:"path"`
	Name         string `json:"name"`
	Server       string `json:"server"`
	Username     string `json:"username"`
	Groupname    string `json:"groupname"`
	ResourceType string `json:"resource_type"`
}

type DownloadResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int    `json:"size"`
	Server      string `json:"server"`
	User        string `json:"user"`
	Username    string `json:"username"`
	Groupname   string `json:"groupname"`
	ExpiresAt   string `json:"expires_at"`
	UploadURL   string `json:"upload_url"`
	DownloadURL string `json:"download_url"`
	Command     string `json:"command"`
}

type UploadRequest struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	Server         string `json:"server"`
	Username       string `json:"username"`
	Groupname      string `json:"groupname"`
	AllowUnzip     string `json:"allow_unzip"`
	AllowOverwrite string `json:"allow_overwrite"`
}

type UploadResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int       `json:"size"`
	Server    string    `json:"server"`
	User      string    `json:"user"`
	Username  string    `json:"username"`
	Groupname string    `json:"groupname"`
	ExpiresAt time.Time `json:"expires_at"`
	UploadUrl string    `json:"upload_url"`
	Command   string    `json:"command"`
}
