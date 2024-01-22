package cert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"net/url"
	"strconv"
)

var (
	authorityURL   = "/api/cert/authorities/"
	signRequestURL = "/api/cert/sign_requests/"
	certURL        = "/api/cert/certificates/"
)

func CreateSignRequest(ac *client.AlpaconClient, signRequest SignRequest) (SignRequestResponse, error) {
	var response SignRequestResponse

	responseBody, err := ac.SendPostRequest(signRequestURL, signRequest)
	if err != nil {
		return SignRequestResponse{}, err
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return SignRequestResponse{}, err
	}

	return response, nil
}

func SubmitCSR(ac *client.AlpaconClient, csr []byte, submitURL string) ([]byte, error) {
	var request CSRSubmit
	request.CsrText = string(csr)

	responseBody, err := ac.SendPatchRequest(submitURL, request)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func CreateAuthority(ac *client.AlpaconClient, authorityRequest AuthorityRequest) (AuthorityCreateResponse, error) {
	var response AuthorityCreateResponse

	responseBody, err := ac.SendPostRequest(authorityURL, authorityRequest)
	if err != nil {
		return AuthorityCreateResponse{}, err
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return AuthorityCreateResponse{}, err
	}

	return response, nil
}

func GetCSRList(ac *client.AlpaconClient, state string) ([]CSRAttributes, error) {
	var csrList []CSRAttributes
	page := 1
	const pageSize = 100

	for {
		responseBody, err := ac.SendGetRequest(buildURL(state, page, pageSize))
		if err != nil {
			return nil, err
		}

		var response CSRListResponse
		if err = json.Unmarshal(responseBody, &response); err != nil {
			return nil, err
		}

		for _, csr := range response.Results {
			csrList = append(csrList, CSRAttributes{
				Id:            csr.Id,
				Name:          csr.CommonName,
				Authority:     csr.AuthorityName,
				DomainList:    csr.DomainList,
				IpList:        csr.IpList,
				State:         csr.State,
				RequestedIp:   csr.RequestedIp,
				RequestedBy:   csr.RequestedByName,
				RequestedDate: utils.TimeUtils(csr.AddedAt),
			})
		}

		if len(response.Results) < pageSize {
			break
		}
		page++
	}

	return csrList, nil
}

func GetAuthorityList(ac *client.AlpaconClient) ([]AuthorityAttributes, error) {
	var authorityList []AuthorityAttributes
	page := 1
	const pageSize = 100

	for {
		params := utils.CreatePaginationParams(page, pageSize)
		responseBody, err := ac.SendGetRequest(authorityURL + "?" + params)
		if err != nil {
			return nil, err
		}

		var response AuthorityListResponse
		if err = json.Unmarshal(responseBody, &response); err != nil {
			return nil, err
		}

		for _, authority := range response.Results {
			authorityList = append(authorityList, AuthorityAttributes{
				Id:               authority.Id,
				Name:             authority.Name,
				Organization:     authority.Organization,
				Domain:           authority.Domain,
				RootValidDays:    authority.RootValidDays,
				DefaultValidDays: authority.DefaultValidDays,
				MaxValidDays:     authority.MaxValidDays,
				Server:           authority.AgentName,
				Owner:            authority.OwnerName,
				SignedAt:         utils.TimeUtils(authority.SignedAt),
			})
		}

		if len(response.Results) < pageSize {
			break
		}
		page++
	}

	return authorityList, nil
}

func GetAuthorityDetail(ac *client.AlpaconClient, authorityId string) ([]byte, error) {
	body, err := ac.SendGetRequest(authorityURL + authorityId)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetCSRDetail(ac *client.AlpaconClient, csrId string) ([]byte, error) {
	body, err := ac.SendGetRequest(signRequestURL + csrId)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetCertificateDetail(ac *client.AlpaconClient, certId string) ([]byte, error) {
	body, err := ac.SendGetRequest(certURL + certId)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func ApproveCSR(ac *client.AlpaconClient, csrId string) ([]byte, error) {
	responseBody, err := ac.SendPostRequest(signRequestURL+csrId+"/approve/", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func DenyCSR(ac *client.AlpaconClient, csrId string) ([]byte, error) {
	responseBody, err := ac.SendPostRequest(signRequestURL+csrId+"/deny/", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func DeleteCSR(ac *client.AlpaconClient, csrId string) error {
	_, err := ac.SendDeleteRequest(signRequestURL + csrId + "/")
	if err != nil {
		return err
	}

	return err
}

func GetCertificateList(ac *client.AlpaconClient) ([]CertificateAttributes, error) {
	var certList []CertificateAttributes
	page := 1
	const pageSize = 100

	for {
		params := utils.CreatePaginationParams(page, pageSize)
		responseBody, err := ac.SendGetRequest(certURL + "?" + params)
		if err != nil {
			return nil, err
		}

		var response CertificateListResponse
		if err = json.Unmarshal(responseBody, &response); err != nil {
			return nil, err
		}

		for _, cert := range response.Results {
			certList = append(certList, CertificateAttributes{
				Id:        cert.Id,
				Authority: cert.Authority,
				Csr:       cert.Csr,
				ValidDays: cert.ValidDays,
				SignedAt:  utils.TimeUtils(cert.SignedAt),
				ExpiresAt: utils.TimeUtils(cert.ExpiresAt),
				SignedBy:  cert.SignedByName,
			})
		}

		if len(response.Results) < pageSize {
			break
		}
		page++
	}

	return certList, nil
}

func buildURL(state string, page int, pageSize int) string {
	params := url.Values{}
	params.Add("state", state)
	params.Add("page", strconv.Itoa(page))
	params.Add("page_size", fmt.Sprintf("%d", pageSize))
	return signRequestURL + "?" + params.Encode()
}
