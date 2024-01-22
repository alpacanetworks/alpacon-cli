package packages

import (
	"bytes"
	"encoding/json"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"mime/multipart"
)

const (
	systemPackageEntryURL = "/api/packages/system/entries/"
	pythonPackageEntryURL = "/api/packages/python/entries/"
)

func GetSystemPackageEntry(ac *client.AlpaconClient) ([]SystemPackage, error) {
	var packageList []SystemPackage
	page := 1
	const pageSize = 100

	for {
		params := utils.CreatePaginationParams(page, pageSize)
		responseBody, err := ac.SendGetRequest(systemPackageEntryURL + "?" + params)
		if err != nil {
			return nil, err
		}

		var response SystemPackageListResponse
		if err = json.Unmarshal(responseBody, &response); err != nil {
			return nil, err
		}

		for _, packages := range response.Results {
			packageList = append(packageList, SystemPackage{
				Name:     packages.Name,
				Version:  packages.Version,
				Arch:     packages.Arch,
				Platform: packages.Platform,
				Owner:    packages.OwnerName,
			})
		}

		if len(response.Results) < pageSize {
			break
		}
		page++
	}
	return packageList, nil
}

func GetPythonPackageEntry(ac *client.AlpaconClient) ([]PythonPackage, error) {
	var packageList []PythonPackage
	page := 1
	const pageSize = 100

	for {
		params := utils.CreatePaginationParams(page, pageSize)
		responseBody, err := ac.SendGetRequest(pythonPackageEntryURL + "?" + params)
		if err != nil {
			return nil, err
		}

		var response PythonPackageListResponse
		if err = json.Unmarshal(responseBody, &response); err != nil {
			return nil, err
		}

		for _, packages := range response.Results {
			packageList = append(packageList, PythonPackage{
				Name:         packages.Name,
				Version:      packages.Version,
				PythonTarget: packages.Target,
				ABI:          packages.ABI,
				Platform:     packages.Platform,
				Owner:        packages.OwnerName,
			})
		}

		if len(response.Results) < pageSize {
			break
		}
		page++
	}
	return packageList, nil
}

func UploadPackage(ac *client.AlpaconClient, file string, packageType string) error {
	content, err := utils.ReadFileFromPath(file)
	if err != nil {
		return err
	}

	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)

	fileWriter, err := multiPartWriter.CreateFormFile("content", file)
	if err != nil {
		return err
	}
	_, err = fileWriter.Write(content)
	if err != nil {
		return err
	}
	multiPartWriter.Close()

	var requestURL string
	if packageType == "python" {
		requestURL = pythonPackageEntryURL
	} else {
		requestURL = systemPackageEntryURL
	}

	err = ac.SendMultipartRequest(requestURL, multiPartWriter, requestBody)
	if err != nil {
		return err
	}

	return nil
}
