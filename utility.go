package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
