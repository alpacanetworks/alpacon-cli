package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"reflect"
)

func PrintTable(slice interface{}) {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		CliError("Parsing data: Expected a list format.")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetHeaderLine(false)
	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	headers := make([]string, s.Type().Elem().NumField())
	for i := 0; i < s.Type().Elem().NumField(); i++ {
		headers[i] = s.Type().Elem().Field(i).Name
	}
	table.SetHeader(headers)

	for i := 0; i < s.Len(); i++ {
		row := make([]string, s.Type().Elem().NumField())
		for j := 0; j < s.Type().Elem().NumField(); j++ {
			value := s.Index(i).Field(j)
			row[j] = fmt.Sprintf("%v", value)
		}
		table.Append(row)
	}

	table.Render()
}

func PrintJson(body []byte) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "    ")
	if err != nil {
		CliError("Parsing data: Expected a json format: %s", err)
	}

	fmt.Println(prettyJSON.String())
}
