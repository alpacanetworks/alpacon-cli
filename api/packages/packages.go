package packages

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"io"
	"mime/multipart"
	"path/filepath"
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

func GetPackageIDByName(ac *client.AlpaconClient, fileName string, packageType string) (string, error) {
	var url string

	if packageType == "python" {
		url = pythonPackageEntryURL
	} else {
		url = systemPackageEntryURL
	}

	body, err := ac.SendGetRequest(url + "?name=" + fileName)
	if err != nil {
		return "", err
	}

	var response PythonPackageListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response.Count == 0 {
		return "", errors.New("no server found with the given name")
	}

	return response.Results[0].ID, nil
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

	_, err = ac.SendMultipartRequest(requestURL, multiPartWriter, requestBody)
	if err != nil {
		return err
	}

	return nil
}

func DownloadPackage(ac *client.AlpaconClient, fileName string, dest string, packageType string) error {
	packageID, err := GetPackageIDByName(ac, fileName, packageType)
	if err != nil {
		return err
	}

	var url string
	if packageType == "python" {
		url = pythonPackageEntryURL
	} else {
		url = systemPackageEntryURL
	}

	type DownloadURL struct {
		DownloadURL string `json:"download_url"`
	}

	var downloadURL DownloadURL

	respBody, err := ac.SendGetRequest(url + packageID)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, &downloadURL)
	if err != nil {
		return err
	}

	resp, err := ac.SendGetRequestForDownload(utils.RemovePrefixBeforeAPI(downloadURL.DownloadURL))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	savePath := filepath.Join(dest, filepath.Base(fileName))
	err = utils.SaveFile(savePath, respBody)
	if err != nil {
		return err
	}

	return nil
}
