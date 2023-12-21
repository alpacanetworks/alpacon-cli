package note

import (
	"encoding/json"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/iam"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"net/url"
)

var (
	noteURL = "/api/servers/notes/"
)

func GetNoteList(ac *client.AlpaconClient, serverName string, pageSize int) ([]NoteDetails, error) {
	var noteList []NoteDetails
	var serverID string
	var err error

	if serverName != "" {
		serverID, err = server.GetServerIDByName(ac, serverName)
		if err != nil {
			return nil, err
		}
	}

	url := buildURL(serverID, pageSize)

	responseBody, err := ac.SendGetRequest(url)
	if err != nil {
		return nil, err
	}

	var response NoteListResponse
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return nil, err
	}

	for _, note := range response.Results {
		userName, err := iam.GetUserNameByID(ac, note.Author)
		if err != nil {
			return nil, err
		}

		noteList = append(noteList, NoteDetails{
			ID:      note.ID,
			Server:  serverName,
			Author:  userName,
			Content: note.Content,
		})
	}

	return noteList, nil
}

func DeleteNote(ac *client.AlpaconClient, noteID string) error {
	_, err := ac.SendDeleteRequest(noteURL + noteID + "/")
	if err != nil {
		return err
	}

	return err
}

func buildURL(serverID string, pageSize int) string {
	params := url.Values{}
	params.Add("server", serverID)
	params.Add("page_size", fmt.Sprintf("%d", pageSize))
	return noteURL + "?" + params.Encode()
}
