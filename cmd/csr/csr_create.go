package csr

import (
	"fmt"
	certApi "github.com/alpacanetworks/alpacon-cli/api/cert"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/pkg/cert"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
	"os/user"
	"path/filepath"
)

var (
	defaultPrivateKeyDir string
	defaultCSRDir        string
)

const (
	infoMessage = "Please specify the paths for the private key and CSR files.\n" +
		"If an existing key is found at the specified path, it will be used.\n" +
		"Otherwise, a new key will be generated.\n" +
		"Note: Root permission may be required for certain paths."
)

var csrCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a CSR",
	Long: `
 	Generates a new Certificate Signing Request based on provided information, 
	which can then be submitted for signing to a certificate authority.
	`,
	Example: `alpacon csr create`,
	Run: func(cmd *cobra.Command, args []string) {

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		signRequest, certPath := promptForCert()

		// 1. sign request (domain, ip, valid days only)
		response, err := certApi.CreateSignRequest(alpaconClient, signRequest)
		if err != nil {
			utils.CliError("Failed to send sign request to server: %s", err)
		}

		// 2. 받아온 template 으로 CSR 파일 생성
		csr, err := cert.CreateCSR(response, certPath)
		if err != nil {
			utils.CliError("Failed to create CSR file %s", err)
		}

		// 3. 생성한 CSR 파일로 submit action
		_, err = certApi.SubmitCSR(alpaconClient, csr, response.SubmitURL)
		if err != nil {
			utils.CliError("Failed to submit CSR file to server %s", err)
		}

		// 4. holding?
	},
}

func init() {
	usr, err := user.Current()
	if err != nil {
		utils.CliError("Failed to obtain the current user information: " + err.Error())
	}

	defaultPrivateKeyDir = filepath.Join(usr.HomeDir, "tmp/private/")
	defaultCSRDir = filepath.Join(usr.HomeDir, "tmp/")
}

func promptForCert() (certApi.SignRequest, cert.CertificatePath) {
	var signRequest certApi.SignRequest
	var certPath cert.CertificatePath

	signRequest.DomainList = utils.PromptForRequiredListInput("domain list (e.g., domain1.com, domain2.com): ")
	signRequest.IpList = utils.PromptForRequiredListInput("ip list (e.g., 192.168.1.1, 10.0.0.1): ")
	signRequest.ValidDays = utils.PromptForInput("valid days: ")

	domainName := signRequest.DomainList[0]
	defaultKeyPath := fmt.Sprintf("%s/%s.key", defaultPrivateKeyDir, domainName)
	defaultCSRPath := fmt.Sprintf("%s/%s.csr", defaultCSRDir, domainName)

	utils.CliInfo(infoMessage)

	certPath.PrivateKeyPath = utils.PromptForInput(fmt.Sprintf("Path for the Private Key file[`%s`]: ", defaultKeyPath))
	if certPath.PrivateKeyPath == "" {
		certPath.PrivateKeyPath = defaultKeyPath
	}

	certPath.CSRPath = utils.PromptForInput(fmt.Sprintf("Path for the CSR file[`%s`]: ", defaultCSRPath))
	if certPath.CSRPath == "" {
		certPath.CSRPath = defaultCSRPath
	}

	return signRequest, certPath
}
