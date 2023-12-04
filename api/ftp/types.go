package ftp

type DownloadRequest struct {
	Path   string `json:"path"`
	Server string `json:"server"`
}

type DownloadResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        string `json:"size"`
	Server      string `json:"server"`
	User        string `json:"user"`
	ExpiresAt   string `json:"expires_at"`
	UploadURL   string `json:"upload_url"`
	DownloadURL string `json:"download_url"`
}
