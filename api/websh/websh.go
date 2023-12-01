package websh

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/gorilla/websocket"
	"golang.org/x/term"
	"io"
	"net/http"
	"os"
)

var (
	createSessionURL = "/api/websh/sessions/"
)

func CreateWebshConnection(ac *client.AlpaconClient, serverName string, root bool) error {
	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return err
	}

	sessionResponse, err := createWebshSession(ac, serverID, root)
	if err != nil {
		return err
	}

	err = openNewTerminal(ac, sessionResponse)
	if err != nil {
		return err
	}

	return nil
}

// Create new websh session
func createWebshSession(ac *client.AlpaconClient, serverID string, root bool) (SessionResponse, error) {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return SessionResponse{}, err
	}

	sessionRequest := &SessionRequest{
		Server: serverID,
		Root:   root,
		Rows:   height,
		Cols:   width,
	}

	responseBody, err := ac.SendPostRequest(createSessionURL, sessionRequest)
	if err != nil {
		return SessionResponse{}, err
	}

	var response SessionResponse
	err = json.Unmarshal(responseBody, &response)
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return response, err
	}

	return response, nil
}

// Open terminal
func openNewTerminal(ac *client.AlpaconClient, sessionResponse SessionResponse) error {
	headers := http.Header{"Origin": []string{ac.BaseURL}}

	conn, _, err := websocket.DefaultDialer.Dial(sessionResponse.WebsocketURL, headers)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err = websocketClient(conn); err != nil {
		return err
	}
	return nil
}

func websocketClient(conn *websocket.Conn) error {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("websh command should be a terminal")
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	done := make(chan error)

	// Goroutine for reading messages from the server
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				done <- err
				return
			}
			fmt.Print(string(message))
		}
	}()

	// Goroutine for reading user input and sending it to the server
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			char, _, err := reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					done <- nil
					return
				}
				done <- err

				return
			}
			err = conn.WriteMessage(websocket.TextMessage, []byte(string(char)))
			if err != nil {
				done <- err
				return
			}
		}
	}()

	return <-done
}