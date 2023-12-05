package utils

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func PromptForInput(promptText string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptText)
	input, err := reader.ReadString('\n')
	if err != nil {
		CliError("Error during input. Please try again.")
	}
	return strings.TrimSpace(input)
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
