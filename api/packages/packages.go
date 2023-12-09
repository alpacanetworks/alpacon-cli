package packages

import (
	"encoding/json"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
)

var (
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
