package iam

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"net/url"
)

var (
	userURL      = "/api/iam/users/"
	getUserIDURL = "/api/iam/users/?username="

	groupURL      = "/api/iam/groups/"
	getGroupIDURL = "/api/iam/groups/?name="

	getUserMembershipURL  = "/api/iam/memberships/?user="
	getGroupMembershipURL = "/api/iam/memberships/?group="
	membershipURL         = "/api/iam/memberships/"
)

func GetUserList(ac *client.AlpaconClient) ([]UserAttributes, error) {
	var userList []UserAttributes
	page := 1
	const pageSize = 100

	for {
		params := utils.CreatePaginationParams(page, pageSize)
		responseBody, err := ac.SendGetRequest(userURL + "?" + params)
		if err != nil {
			return nil, err
		}

		var response UserListResponse
		if err = json.Unmarshal(responseBody, &response); err != nil {
			return nil, err
		}

		for _, user := range response.Results {
			userList = append(userList, UserAttributes{
				Username:   user.Username,
				Name:       fmt.Sprintf("%s %s", user.LastName, user.FirstName),
				Email:      user.Email,
				Tags:       user.Tags,
				Groups:     user.NumGroups,
				UID:        user.UID,
				Status:     getUserStatus(user.IsActive, user.IsStaff, user.IsSuperuser),
				LDAPStatus: getLDAPStatus(user.IsLDAPUser),
			})
		}

		if len(response.Results) < pageSize {
			break
		}
		page++
	}
	return userList, nil
}

func GetGroupList(ac *client.AlpaconClient) ([]GroupAttributes, error) {
	var groupList []GroupAttributes
	page := 1
	const pageSize = 100

	for {
		params := utils.CreatePaginationParams(page, pageSize)
		responseBody, err := ac.SendGetRequest(groupURL + "?" + params)
		if err != nil {
			return nil, err
		}

		var response GroupListResponse
		if err = json.Unmarshal(responseBody, &response); err != nil {
			return nil, err
		}

		for _, group := range response.Results {
			groupList = append(groupList, GroupAttributes{
				Name:        group.Name,
				DisplayName: group.DisplayName,
				Tags:        group.Tags,
				Members:     group.NumMembers,
				Servers:     len(group.Servers),
				GID:         group.GID,
				LDAPStatus:  getLDAPStatus(group.IsLDAPGroup),
			})
		}

		if len(response.Results) < pageSize {
			break
		}
		page++
	}
	return groupList, nil
}

func GetUserDetail(ac *client.AlpaconClient, userName string) ([]byte, error) {
	var userDetails UserDetails

	userID, err := GetUserIDByName(ac, userName)
	if err != nil {
		return nil, err
	}

	responseBody, err := ac.SendGetRequest(userURL + userID)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &userDetails)
	if err != nil {
		return nil, err
	}

	groupNames, err := getGroupNames(ac, userID)
	if err != nil {
		return nil, err
	}

	userDetailAttributes := &UserDetailAttributes{
		Username:      userDetails.Username,
		Name:          fmt.Sprintf("%s %s", userDetails.LastName, userDetails.FirstName),
		Description:   userDetails.Description,
		Email:         userDetails.Email,
		Phone:         userDetails.Phone,
		UID:           userDetails.UID,
		Shell:         userDetails.Shell,
		HomeDirectory: userDetails.HomeDirectory,
		NumGroups:     userDetails.NumGroups,
		Groups:        groupNames,
		Tags:          userDetails.Tags,
		Status:        getUserStatus(userDetails.IsActive, userDetails.IsStaff, userDetails.IsSuperuser),
		LDAPStatus:    getLDAPStatus(userDetails.IsLDAPUser),
	}

	userDetailJSON, err := json.Marshal(userDetailAttributes)
	if err != nil {
		return nil, err
	}

	return userDetailJSON, nil
}

func GetGroupDetail(ac *client.AlpaconClient, groupName string) ([]byte, error) {
	var groupDetails GroupDetails

	groupID, err := GetGroupIDByName(ac, groupName)
	if err != nil {
		return nil, err
	}

	responseBody, err := ac.SendGetRequest(groupURL + groupID)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &groupDetails)
	if err != nil {
		return nil, err
	}

	memberNames, err := getMemberNames(ac, groupID)
	if err != nil {
		return nil, err
	}

	groupDetailAttributes := &GroupDetailAttributes{
		Name:         groupDetails.Name,
		DisplayName:  groupDetails.DisplayName,
		Tags:         groupDetails.Tags,
		Description:  groupDetails.Description,
		NumMembers:   groupDetails.NumMembers,
		Members:      memberNames,
		GID:          groupDetails.GID,
		LDAPStatus:   getLDAPStatus(groupDetails.IsLDAPGroup),
		Servers:      len(groupDetails.Servers),
		ServersNames: groupDetails.ServersNames,
	}

	groupDetailJson, err := json.Marshal(groupDetailAttributes)
	if err != nil {
		return nil, err
	}

	return groupDetailJson, nil
}

func CreateUser(ac *client.AlpaconClient, userRequest UserCreateRequest) error {
	userRequest.IsActive = true
	_, err := ac.SendPostRequest(userURL, userRequest)
	if err != nil {
		return err
	}

	return nil
}

func CreateGroup(ac *client.AlpaconClient, groupRequest GroupCreateRequest) error {
	_, err := ac.SendPostRequest(groupURL, groupRequest)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(ac *client.AlpaconClient, userName string) error {
	userID, err := GetUserIDByName(ac, userName)
	if err != nil {
		return err
	}

	_, err = ac.SendDeleteRequest(userURL + userID + "/")
	if err != nil {
		return err
	}

	return err
}

func DeleteGroup(ac *client.AlpaconClient, groupName string) error {
	groupID, err := GetGroupIDByName(ac, groupName)
	if err != nil {
		return err
	}

	_, err = ac.SendDeleteRequest(groupURL + groupID + "/")
	if err != nil {
		return err
	}

	return err
}

func AddMember(ac *client.AlpaconClient, memberRequest MemberAddRequest) error {
	var err error
	memberRequest.Group, err = GetGroupIDByName(ac, memberRequest.Group)
	if err != nil {
		return err
	}

	memberRequest.User, err = GetUserIDByName(ac, memberRequest.User)
	if err != nil {
		return err
	}

	_, err = ac.SendPostRequest(membershipURL, memberRequest)
	if err != nil {
		return err
	}

	return nil
}

func DeleteMember(ac *client.AlpaconClient, memberDeleteRequest MemberDeleteRequest) error {
	groupID, err := GetGroupIDByName(ac, memberDeleteRequest.Group)
	if err != nil {
		return err
	}

	memberID, err := GetUserIDByName(ac, memberDeleteRequest.User)
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("user", memberID)
	params.Add("group", groupID)

	responseBody, err := ac.SendGetRequest(membershipURL + "?" + params.Encode())
	if err != nil {
		return err
	}

	var memberDetails []MemberDetailResponse
	err = json.Unmarshal(responseBody, &memberDetails)
	if err != nil {
		return err
	}

	_, err = ac.SendDeleteRequest(membershipURL + memberDetails[0].ID + "/")
	if err != nil {
		return err
	}

	return err
}

func GetUserIDByName(ac *client.AlpaconClient, userName string) (string, error) {
	responseBody, err := ac.SendGetRequest(getUserIDURL + userName)
	if err != nil {
		return "", err
	}

	var response UserListResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", err
	}

	if response.Count == 0 {
		return "", errors.New("no user found with the given name")
	}

	return response.Results[0].ID, nil
}

func GetUserNameByID(ac *client.AlpaconClient, userID string) (string, error) {
	responseBody, err := ac.SendGetRequest(userURL + userID)
	if err != nil {
		return "", err
	}

	var response UserDetailAttributes
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", err
	}

	return response.Username, nil
}

func GetGroupIDByName(ac *client.AlpaconClient, groupName string) (string, error) {
	responseBody, err := ac.SendGetRequest(getGroupIDURL + groupName)
	if err != nil {
		return "", err
	}

	var response GroupListResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", err
	}

	if response.Count == 0 {
		return "", errors.New("no group found with the given name")
	}

	return response.Results[0].ID, nil
}

func getUserStatus(isActive bool, isStaff bool, isSuperuser bool) string {
	if isSuperuser {
		return "superuser"
	}
	if isStaff {
		return "staff"
	}
	if isActive {
		return "active"
	}
	return "inactive"
}

func getLDAPStatus(isLDAP bool) string {
	if isLDAP {
		return "ldap"
	}

	return "local"
}

func getGroupNames(ac *client.AlpaconClient, userID string) ([]string, error) {
	responseBody, err := ac.SendGetRequest(getUserMembershipURL + userID)
	if err != nil {
		return nil, err
	}

	var membershipResponse []Membership
	err = json.Unmarshal(responseBody, &membershipResponse)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, membership := range membershipResponse {
		result = append(result, membership.GroupName)
	}

	return result, nil
}

func getMemberNames(ac *client.AlpaconClient, groupID string) ([]string, error) {
	responseBody, err := ac.SendGetRequest(getGroupMembershipURL + groupID)
	if err != nil {
		return nil, err
	}

	var membershipResponse []Membership
	err = json.Unmarshal(responseBody, &membershipResponse)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, membership := range membershipResponse {
		result = append(result, membership.UserName)
	}

	return result, nil
}
