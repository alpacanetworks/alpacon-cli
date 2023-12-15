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

func PromptForRequiredInput(prompt string) string {
	for {
		input := PromptForInput(prompt)
		if input != "" {
			return input
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

func CreatePaginationParams(page, pageSize int) string {
	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	params.Add("page_size", strconv.Itoa(pageSize))
	return params.Encode()
}

func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

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
