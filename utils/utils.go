package utils

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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

func PromptForPassword(promptText string) string {
	fmt.Print(promptText)
	bytePassword, err := term.ReadPassword(0)
	if err != nil {
		return ""
	}
	fmt.Println()
	return strings.TrimSpace(string(bytePassword))
}

func PromptForInput(promptText string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptText)
	input, err := reader.ReadString('\n')
	if err != nil {
		CliError("Error during input. Please try again.")
	}
	return strings.TrimSpace(input)
}

func PromptForRequiredInput(promptText string) string {
	for {
		input := PromptForInput(promptText)
		if input != "" {
			return input
		}
		fmt.Println("This field is required. Please enter a value.")
	}
}

func PromptForIntInput(promptText string) int {
	inputStr := PromptForInput(promptText)
	inputInt, err := strconv.Atoi(inputStr)
	if err != nil {
		fmt.Println("Only integers are allowed. Please try again.")
		return PromptForIntInput(promptText)
	}
	return inputInt
}

func PromptForListInput(promptText string) []string {
	inputStr := PromptForInput(promptText)
	inputList := strings.Split(inputStr, ",")
	for i, item := range inputList {
		inputList[i] = strings.TrimSpace(item)
	}
	return inputList
}

func PromptForRequiredListInput(promptText string) []string {
	for {
		inputList := PromptForListInput(promptText)
		if len(inputList) > 0 && inputList[0] != "" {
			return inputList
		}
		fmt.Println("This field is required. Please enter a value.")
	}
}

func PromptForBool(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("Invalid input. Please enter 'y' (yes) or 'n' (no).")
		}
	}
}

func SplitAndParseInts(input string) []int {
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

func CreatePaginationParams(page int, pageSize int) string {
	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	params.Add("page_size", strconv.Itoa(pageSize))
	return params.Encode()
}

func TimeUtils(t time.Time) string {
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
	defer file.Close()

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

func StringToStringPointer(value string) *string {
	if value == "" {
		return nil
	} else {
		return &value
	}
}
