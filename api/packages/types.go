package packages

type SystemPackage struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Arch     string `json:"arch"`
	Platform string `json:"platform"`
	Owner    string `json:"owner"`
}

type PythonPackage struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	PythonTarget string `json:"python_target"`
	ABI          string `json:"abi"`
	Platform     string `json:"platform"`
	Owner        string `json:"owner"`
}

type PythonPackageDetail struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Target      string `json:"target"`
	ABI         string `json:"abi"`
	Platform    string `json:"platform"`
	Filesize    int64  `json:"filesize"`
	Owner       string `json:"owner"`
	OwnerName   string `json:"owner_name"`
	AddedAt     string `json:"added_at"`
	DownloadURL string `json:"download_url"`
}

type SystemPackageDetail struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Platform    string `json:"platform"`
	Arch        string `json:"arch"`
	Filesize    int64  `json:"filesize"`
	Owner       string `json:"owner"`
	OwnerName   string `json:"owner_name"`
	AddedAt     string `json:"added_at"`
	DownloadURL string `json:"download_url"`
}
