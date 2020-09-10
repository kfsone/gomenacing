// Parser for simple comma-separated-file translation.

package parsing

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"strconv"
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
func ParseCSV(source io.Reader, fields []string) (<-chan []string, error) {
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

	channel := make(chan []string, 4)
	go func() {
		// Break remaining lines up by
		scanner, fieldOrder, channel := scanner, fieldOrder, channel
		var separator = []byte(",")
		defer close(channel)
		for scanner.Scan() {
			result := make([]string, len(fields))
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
				result[fieldNo] = string(value)
			}
			channel <- result
		}
	}()

	return channel, nil
}

func stringsToUint64s(from []string) (into []uint64, err error) {
	into = make([]uint64, len(from))
	for idx, value := range from {
		if into[idx], err = strconv.ParseUint(value, 10, 64); err != nil {
			return nil, err
		}
	}
	return into, nil
}

func ParseCSVToUint64s(source io.Reader, fields []string) (<-chan []uint64, error) {
	incoming, err := ParseCSV(source, fields)
	if err != nil {
		return nil, err
	}
	arraysOut := make(chan []uint64, 2)
	go func() {
		defer close(arraysOut)
		for stringsIn := range incoming {
			if row, err := stringsToUint64s(stringsIn); err != nil {
				log.Printf("expected numeric value: ")
			} else {
				arraysOut <- row
			}
		}
	}()
	return arraysOut, nil
}
