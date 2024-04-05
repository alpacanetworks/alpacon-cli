package websh

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alpacanetworks/alpacon-cli/api/server"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/gorilla/websocket"
	"golang.org/x/term"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	createSessionURL = "/api/websh/sessions/"
)

func CreateWebshConnection(ac *client.AlpaconClient, serverName, username, groupname string) (SessionResponse, error) {
	var sessionResponse SessionResponse
	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return sessionResponse, err
	}

	return createWebshSession(ac, serverID, username, groupname)
}

// Create new websh session
func createWebshSession(ac *client.AlpaconClient, serverID, username, groupname string) (SessionResponse, error) {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return SessionResponse{}, err
	}

	sessionRequest := &SessionRequest{
		Server:    serverID,
		Username:  username,
		Groupname: groupname,
		Rows:      height,
		Cols:      width,
	}

	responseBody, err := ac.SendPostRequest(createSessionURL, sessionRequest)
	if err != nil {
		return SessionResponse{}, err
	}

	var response SessionResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return SessionResponse{}, nil
	}

	return response, nil
}

// Handles graceful termination of the websh terminal.
// Exits on error without further error handling.
func OpenNewTerminal(ac *client.AlpaconClient, sessionResponse SessionResponse) error {
	headers := http.Header{"Origin": []string{ac.BaseURL}}

	conn, _, err := websocket.DefaultDialer.Dial(sessionResponse.WebsocketURL, headers)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err = runWsClient(conn); err != nil {
		return err
	}
	return nil
}

func runWsClient(conn *websocket.Conn) error {
	oldState, err := checkTerminal()
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	done := make(chan error, 1)
	inputChan := make(chan string, 1)

	go readFromServer(conn, done)
	go readUserInput(inputChan, done)
	go writeToServer(conn, inputChan, done)

	return <-done
}

func checkTerminal() (*term.State, error) {
	utils.ShowLogo()
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, errors.New("websh command should be a terminal")
	}
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}

	return oldState, nil
}

func readFromServer(conn *websocket.Conn, done chan<- error) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			done <- err
			return
		}
		fmt.Print(string(message))
	}
}

func readUserInput(inputChan chan<- string, done chan<- error) {
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
		inputChan <- string(char)
	}
}

func writeToServer(conn *websocket.Conn, inputChan <-chan string, done chan<- error) {
	var inputBuffer []rune
	for {
		select {
		case input := <-inputChan:
			inputBuffer = append(inputBuffer, []rune(input)...)
		case <-time.After(time.Millisecond * 5):
			if len(inputBuffer) > 0 {
				err := conn.WriteMessage(websocket.TextMessage, []byte(string(inputBuffer)))
				if err != nil {
					done <- err
					return
				}
				inputBuffer = []rune{}
			}
		}
	}
}
