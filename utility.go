package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

// Miscellaneous utility functions.

// Captures any log output produced by a test function and returns it as a string.
func captureLog(t *testing.T, test func(t *testing.T)) []string {
	// https://stackoverflow.com/questions/44119951/how-to-check-a-log-output-in-go-test
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	test(t)
	if buf.Len() == 0 {
		return nil
	}
	return strings.Split(strings.TrimSpace(buf.String()), "\n")
}

// ensureDirectory creates a directory and its parents if it does not already exist.
func ensureDirectory(path string) (created bool, err error) {
	// Check whether it exists, and if it does check it's a directory.
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		err = os.ErrExist
	}
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0750); err == nil {
			created = true
		}
	}
	return
}

func stringToFeaturePad(padSize string) FacilityFeatureMask {
	if len(padSize) == 1 {
		switch strings.ToUpper(padSize)[0] {
		case 'L':
			return FeatLargePad
		case 'M':
			return FeatMediumPad
		case 'S':
			return FeatSmallPad
		default:
			break
		}
	}
	return FacilityFeatureMask(0)
}

type JsonEntry struct {
	LineNo  int
	Results []gjson.Result
}

// Reads a file and invokes the specified callback for each line in it.
// If an error occurs, the error will be decorated with the filename and line no.
func GetJsonItemsFromFile(filename string, source io.Reader, fieldNames []string, errorCh chan<- error) chan JsonEntry {
	scanner := bufio.NewScanner(source)
	lineNo := 0

	linesCh := make(chan struct {
		int
		string
	}, 8)
	jsonsCh := make(chan JsonEntry, 8)

	go func() {
		defer close(linesCh)
		for scanner.Scan() {
			line := scanner.Text()
			lineNo += 1
			linesCh <- struct {
				int
				string
			}{lineNo, line}
		}
		if scanner.Err() != nil {
			errorCh <- fmt.Errorf("%s:%d: parse error: %w", filename, lineNo, scanner.Err())
		}
	}()

	go func() {
		defer close(jsonsCh)
		for line := range linesCh {
			if !gjson.Valid(line.string) {
				errorCh <- fmt.Errorf("%s:%d: invalid json: %s", filename, line.int, line.string)
				continue
			}
			results := gjson.GetMany(line.string, fieldNames...)
			// Check if any were bad
			var badData bool
			for _, field := range results {
				if field.Type == 0 && len(field.Raw) == 0 {
					badData = true
					break
				}
			}
			if badData == true {
				errorCh <- fmt.Errorf("%s:%d: bad entry: %s", filename, line.int, line.string)
				continue
			}
			jsonsCh <- JsonEntry{line.int, results}
		}
	}()

	return jsonsCh
}
