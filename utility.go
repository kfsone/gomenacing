package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tidwall/gjson"
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

func GetTestDir(path ...string) TestDir {
	pathname := filepath.Join(path...)
	if dirname, err := ioutil.TempDir(pathname, "menace"); err != nil {
		panic(err)
	} else {
		return TestDir(dirname)
	}
}

func (td *TestDir) Close() {
	if err := os.RemoveAll(string(*td)); err != nil {
		panic(err)
	}
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

type JSONEntry struct {
	LineNo  int
	Results []gjson.Result
}

// GetJSONItemsFromFile reads a file and invokes the specified callback for each line in it.
// If an error occurs, the error will be decorated with the filename and line no.
func GetJSONItemsFromFile(filename string, source io.Reader, fieldNames []string, errorsCh chan<- error) chan JSONEntry {
	scanner := bufio.NewScanner(source)
	lineNo := 0

	linesCh := make(chan struct {
		int
		string
	}, 8)
	jsonsCh := make(chan JSONEntry, 8)

	go func() {
		defer close(linesCh)
		for scanner.Scan() {
			line := scanner.Text()
			lineNo++
			linesCh <- struct {
				int
				string
			}{lineNo, line}
		}
		if scanner.Err() != nil {
			errorsCh <- fmt.Errorf("%s:%d: parse error: %w", filename, lineNo, scanner.Err())
		}
	}()

	go func() {
		defer close(jsonsCh)
		for line := range linesCh {
			if !gjson.Valid(line.string) {
				errorsCh <- fmt.Errorf("%s:%d: invalid json: %s", filename, line.int, line.string)
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
				errorsCh <- fmt.Errorf("%s:%d: bad entry: %s", filename, line.int, line.string)
				continue
			}
			jsonsCh <- JSONEntry{line.int, results}
		}
	}()

	return jsonsCh
}
