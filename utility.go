package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

// Miscellaneous utility functions.

// Captures any log output produced by a test function and returns it as a string.
func captureLog(t *testing.T, test func(t *testing.T)) string {
	// https://stackoverflow.com/questions/44119951/how-to-check-a-log-output-in-go-test
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	test(t)
	return buf.String()
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

// GenerateLinesFromReader exposes the contents of e.g a file, line-at-a-time, over a generator.
func GenerateLinesFromReader(source io.Reader) *Generator {
	scanner := bufio.NewScanner(source)
	outputCh := make(chan interface{})
	generator := NewGenerator(outputCh, func(sink GeneratorSink, generator *Generator) {
		defer close(outputCh)
		defer generator.Close(scanner.Err())
		for !generator.Cancelled() && scanner.Scan() {
			sink <- scanner.Text()
		}
	})

	return generator
}

// Reads a file and invokes the specified callback for each line in it.
// If an error occurs, the error will be decorated with the filename and line no.
func IterateLinesInFile(filename string, file io.Reader, callback func(string) error) error {
	outputCh := make(chan interface{})
	defer close(outputCh)

	lineNo := 0
	lines := GenerateLinesFromReader(file)
	for line := range lines.OutputCh {
		lineNo += 1
		if err := callback(line.(string)); err != nil {
			lines.Cancel(fmt.Errorf("%s:%d: %w", filename, lineNo, err))
			break
		}
	}

	return lines.Error()
}
