// Parser for simple comma-separated-file translation.

package parsing

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

// getFieldOrder will identify the order of the comma-separated
// heading values in `line` as required to represent `fields`,
// as well as the number of headings in `line`.
func getFieldOrder(fields []string, line string) ([]int, int) {
	var fieldOrder = make([]int, 0, len(fields))
	if len(line) == 0 {
		return fieldOrder, 0
	}
	headers := strings.Split(line, ",")
	for _, fieldName := range fields {
		for headerNo, headerName := range headers {
			if strings.EqualFold(fieldName, headerName) {
				fieldOrder = append(fieldOrder, headerNo)
				break
			}
		}
	}

	return fieldOrder, len(headers)
}

// ParseCSV will iterate over comma-separated lines in `source`. The first
// line is presumed to provide headings for the subsequent lines. The columns
// specified in `fields` are extracted and dispatched to the channel in the
// order they are listed in `fields`.
//
// This if you specified fields a `{"win", "for"}` and the CSV heading is
// `for,the,win`, then the csv line `10,20,30` would map to `{30, 10}`.
//
// This is a very naive csv reader, it simply splits on commas, and removes
// leading/trailing double-quotes.
//
func ParseCSV(source io.Reader, fields []string) (<-chan [][]byte, error) {
	// Scan for a first line, which should contain the headers.
	scanner := bufio.NewScanner(source)
	if !scanner.Scan() {
		// If we couldn't scan one line, the file is empty. Close the channel
		// and return it with no error.
		return nil, io.EOF
	}

	// Map the csv column order to the requested fields order.
	fieldOrder, numHeadings := getFieldOrder(fields, scanner.Text())
	if len(fieldOrder) != len(fields) {
		return nil, errors.New("missing fields in header")
	}

	channel := make(chan [][]byte, 2)
	go func() {
		// Break remaining lines up by
		scanner, fieldOrder, channel := scanner, fieldOrder, channel
		var separator = []byte(",")
		defer close(channel)
		for scanner.Scan() {
			result := make([][]byte, len(fields))
			columns := bytes.Split(scanner.Bytes(), separator)
			if len(columns) < numHeadings {
				continue
			}
			for fieldNo, columnIdx := range fieldOrder {
				value := columns[columnIdx]
				if len(value) > 0 && value[0] == '"' {
					value = value[1:]
				}
				if len(value) > 0 && value[len(value)-1] == '"' {
					value = value[:len(value)-1]
				}
				result[fieldNo] = value
			}
			channel <- result
		}
	}()

	return channel, nil
}
