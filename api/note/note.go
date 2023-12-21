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
		serverName, err = server.GetServerNameByID(ac, note.Server)
		if err != nil {
			return nil, err
		}

		userName, err := iam.GetUserNameByID(ac, note.Author)
		if err != nil {
			return nil, err
		}

		noteList = append(noteList, NoteDetails{
			ID:      note.ID,
			Server:  serverName,
			Author:  userName,
			Content: note.Content,
			Private: note.Private,
		})
	}

	return noteList, nil
}

func CreateNote(ac *client.AlpaconClient, noteRequest NoteCreateRequest) error {
	serverID, err := server.GetServerIDByName(ac, noteRequest.Server)
	if err != nil {
		return err
	}

	noteRequest.Server = serverID
	noteRequest.Pinned = false // The default value for the alpacon API server is currently false

	_, err = ac.SendPostRequest(noteURL, noteRequest)
	if err != nil {
		return err
	}

	return nil
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
