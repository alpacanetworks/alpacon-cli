package cert

type SubjectInfo struct {
	Country            []string
	Province           []string
	Locality           []string
	Organization       []string
	OrganizationalUnit []string
	CommonName         string
	EmailAddresses     []string
}

type CertificatePath struct {
	PrivateKeyPath string
	CSRPath        string
}
