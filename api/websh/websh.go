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
	"net/url"
	"os"
	"path"
	"time"
)

const (
	createSessionURL = "/api/websh/sessions/"
)

func JoinWebshSession(ac *client.AlpaconClient, sharedURL, password string) (SessionResponse, error) {
	parsedURL, err := url.Parse(sharedURL)
	if err != nil {
		return SessionResponse{}, err
	}

	sessionID := parsedURL.Query().Get("session")
	if sessionID == "" {
		return SessionResponse{}, errors.New("Invalid URL format")
	}

	joinRequest := &JoinRequest{
		Password: password,
	}

	responseBody, err := ac.SendPostRequest(utils.BuildURL(createSessionURL, path.Join("", sessionID, "join"), nil), joinRequest)
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

// Create new websh session
func CreateWebshSession(ac *client.AlpaconClient, serverName, username, groupname string, share, readOnly bool) (SessionResponse, error) {
	serverID, err := server.GetServerIDByName(ac, serverName)
	if err != nil {
		return SessionResponse{}, err
	}

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

	if share {
		shareRequest := &ShareRequest{
			ReadOnly: readOnly,
		}
		var shareResponse ShareResponse
		responseBody, err = ac.SendPostRequest(utils.BuildURL(createSessionURL, path.Join(response.ID, "share"), nil), shareRequest)
		err = json.Unmarshal(responseBody, &shareResponse)
		if err != nil {
			return SessionResponse{}, nil
		}
		sharingInfo(shareResponse)
	}

	return response, nil
}

// Handles graceful termination of the websh terminal.
// Exits on error without further error handling.
func OpenNewTerminal(ac *client.AlpaconClient, sessionResponse SessionResponse) error {
	headers := ac.SetWebsocketHeader()

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

func sharingInfo(response ShareResponse) {
	header := `Share the following URL to allow access for the current session to someone else.
**Note: The invitee will be required to enter the provided password to access the websh terminal.**`

	instructions := `
To join the shared session:
1. Execute the following command in a terminal:
   $ alpacon websh join --url="%s" --password="%s"
	
2. Or, directly access the session via the shared URL in a web browser.`

	fmt.Println(header)
	fmt.Printf(instructions, response.SharedURL, response.Password)
	fmt.Println()
	fmt.Println("Session Details:")
	fmt.Println("Share URL:    ", response.SharedURL)
	fmt.Println("Password:     ", response.Password)
	fmt.Println("Read Only:    ", response.ReadOnly)
	fmt.Println("Expiration:   ", utils.TimeUtils(response.Expiration))
}
