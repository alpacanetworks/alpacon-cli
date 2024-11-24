package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ShowLogo() {
	alpaconLogo := `
     (` + "`" + `-')  _           _  (` + "`" + `-') (` + "`" + `-')  _                      <-. (` + "`" + `-')_
     (OO ).-/    <-.    \-.(OO ) (OO ).-/  _             .->      \( OO) )
     / ,---.   ,--. )   _.'    \ / ,---.   \-,-----.(` + "`" + `-')----. ,--./ ,--/
     | \ /` + ".`" + `\  |  (` + "`" + `-')(_...--'' | \ /` + ".`" + `\   |  .--./( OO).-.  '|   \ |  |
     '-'|_.' | |  |OO )|  |_.' | '-'|_.' | /_) (` + "`" + `-')( _) | |  ||  . '|  |)
    (|  .-.  |(|  '__ ||  .___.'(|  .-.  | ||  |OO ) \|  |)|  ||  |\    |
     |  | |  | |     |'|  |      |  | |  |(_'  '--'\  '  '-'  '|  | \   |
     ` + "`" + `--' ` + "`" + `--' ` + "`" + `-----' ` + "`" + `--'      ` + "`" + `--' ` + "`" + `--'   ` + "`" + `-----'   ` + "`" + `-----' ` + "`" + `--'  ` + "`" + `--'
    `
	fmt.Println(alpaconLogo)
}

func ReadFileFromPath(filePath string) ([]byte, error) {
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(absolutePath)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func SetUserAgent() string {
	return fmt.Sprintf("%s/%s", "alpacon-cli", GetCLIVersion())
}

func SplitAndParseInt(input string) []int {
	var intValues []int

	stringValues := strings.Split(input, ",")

	for _, stringValue := range stringValues {
		trimmedString := strings.TrimSpace(stringValue)

		intValue, err := strconv.Atoi(trimmedString)
		if err != nil {
			CliError("Invalid input: only integers allowed.")
		}
		intValues = append(intValues, intValue)
	}

	return intValues
}

func TimeUtils(t time.Time) string {
	if t.IsZero() {
		return "None"
	}

	now := time.Now()
	diff := t.Sub(now)

	if diff >= 0 {
		switch {
		case diff < time.Minute:
			return "in a few seconds"
		case diff < time.Hour:
			return fmt.Sprintf("in %d minutes", diff/time.Minute)
		case diff < 24*time.Hour:
			return fmt.Sprintf("in %d hours", diff/time.Hour)
		case diff < 48*time.Hour:
			return "tomorrow"
		default:
			return fmt.Sprintf("in %d days", diff/(24*time.Hour))
		}
	} else {
		diff = -diff
		switch {
		case diff < time.Minute:
			return "just now"
		case diff < time.Hour:
			return fmt.Sprintf("%d minutes ago", diff/time.Minute)
		case diff < 24*time.Hour:
			return fmt.Sprintf("%d hours ago", diff/time.Hour)
		case diff < 48*time.Hour:
			return "yesterday"
		default:
			return fmt.Sprintf("%d days ago", diff/(24*time.Hour))
		}
	}
}

func TimeFormat(value int) *string {
	expiresAt := time.Now().Add(time.Hour * 24 * time.Duration(value))
	formattedExpiresAt := expiresAt.Format(time.RFC3339)

	return &formattedExpiresAt
}

func TruncateString(str string, num int) string {
	if len(str) > num {
		return str[:num] + "..."
	}
	return str
}

func RemovePrefixBeforeAPI(url string) string {
	apiIndex := strings.Index(url, "/api/")
	if apiIndex == -1 {
		return url
	}
	return url[apiIndex:]
}

func SaveFile(fileName string, data []byte) error {
	dir := filepath.Dir(fileName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer func() { _ = file.Close() }()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}

func BoolPointerToString(value *bool) string {
	if value == nil {
		return "null"
	}
	if *value {
		return "true"
	}
	return "false"
}

func BuildURL(basePath, relativePath string, params map[string]string) string {
	u, err := url.Parse(basePath)
	if err != nil {
		CliError("Failed to parse base URL")
	}

	u.Path = path.Join(u.Path, relativePath)
	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}
	q := u.Query()

	for key, value := range params {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

func IsUUID(str string) bool {
	_, err := uuid.Parse(str)
	return err == nil
}

// ProcessEditedData facilitates user modifications to original data,
// formats it, supports editing via a temp file, compares the edited data against the original,
// and parses it into JSON. If no changes are made, the update is aborted and an error is returned.
func ProcessEditedData(originalData []byte) (interface{}, error) {
	prettyJSON, err := PrettyJSON(originalData)
	if err != nil {
		return nil, err
	}

	tmpFile, err := CreateAndEditTempFile(prettyJSON.Bytes())
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.Remove(tmpFile) }()

	editedContent, err := os.ReadFile(tmpFile)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(prettyJSON.Bytes(), editedContent) {
		CliInfoWithExit("No changes made. Aborting update.")
	}

	var jsonData interface{}
	err = json.Unmarshal(editedContent, &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func CreateAndEditTempFile(data []byte) (string, error) {
	tmpl, err := os.CreateTemp("", "example.*.json")
	if err != nil {
		return "", errors.New("Failed to create temp file for update")
	}
	defer func() { _ = tmpl.Close() }()

	if _, err = tmpl.Write(data); err != nil {
		return "", err
	}

	if err = tmpl.Close(); err != nil {
		return "", err
	}

	cmd := exec.Command("vi", tmpl.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		return "", err
	}

	return tmpl.Name(), nil
}

func SplitPath(path string) (string, string) {
	parts := strings.SplitN(path, ":", 2)
	return parts[0], parts[1]
}
