package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func Test_captureLog(t *testing.T) {
	stringList := captureLog(t, func(t *testing.T) {
	})
	assert.Nil(t, stringList)
	stringList = captureLog(t, func(t *testing.T) {
		log.Print("!hello world!")
	})
	assert.Len(t, stringList, 1)
	assert.Contains(t, stringList[0], "!hello world!")
	stringList = captureLog(t, func(t *testing.T) {
		log.Print("line1")
		log.Print("line2")
		log.Print("line3")
	})
	assert.Len(t, stringList, 3)
	assert.Contains(t, stringList[0], "line1")
	assert.Contains(t, stringList[1], "line2")
	assert.Contains(t, stringList[2], "line3")
}

func Test_ensureDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tmpPath, err := ioutil.TempDir("", "go-dangerous-test")
	require.Nil(t, err)
	assert.NotEmpty(t, tmpPath)
	defer func() {
		if err := os.RemoveAll(tmpPath); err != nil {
			log.Fatal(err)
		}
	}()

	path := filepath.Join(tmpPath, "data")

	// Lets make real sure the directory does not exist.
	info, err := os.Stat(path)
	assert.NotNil(t, err)
	assert.Nil(t, info)
	assert.True(t, os.IsNotExist(err))

	// Use ensure directory to create it.
	created, err := ensureDirectory(path)
	assert.Nil(t, err)
	assert.True(t, created)
	info, err = os.Stat(path)
	require.Nil(t, err)
	assert.True(t, info.IsDir())

	// Trying to run it again should be a no-op.
	created, err = ensureDirectory(path)
	assert.Nil(t, err)
	assert.False(t, created)

	// Create something 2-folders down
	path = filepath.Join(tmpPath, "parent", "child", "leaf")
	created, err = ensureDirectory(path)
	assert.Nil(t, err)
	assert.True(t, created)

	info, err = os.Stat(path)
	require.Nil(t, err)
	assert.True(t, info.IsDir())

	info, err = os.Stat(filepath.Join(tmpPath, "parent"))
	require.Nil(t, err)
	assert.True(t, info.IsDir())

	// If the target name is a file, we should get an error.
	path = filepath.Join(tmpPath, "file.test")
	file, err := os.Create(path)
	require.Nil(t, err)
	require.NotNil(t, file)
	assert.Nil(t, file.Close())

	created, err = ensureDirectory(path)
	assert.False(t, created)
	assert.True(t, os.IsExist(err))
}

func Test_stringToFeaturePad(t *testing.T) {
	// Empty string -> empty pad
	assert.Equal(t, FacilityFeatureMask(0), stringToFeaturePad(""))

	// LMS -> values
	assert.Equal(t, FeatLargePad, stringToFeaturePad("L"))
	assert.Equal(t, FeatLargePad, stringToFeaturePad("l"))

	assert.Equal(t, FeatMediumPad, stringToFeaturePad("M"))
	assert.Equal(t, FeatMediumPad, stringToFeaturePad("m"))

	assert.Equal(t, FeatSmallPad, stringToFeaturePad("S"))
	assert.Equal(t, FeatSmallPad, stringToFeaturePad("s"))

	// anything else -> 0
	assert.Equal(t, FacilityFeatureMask(0), stringToFeaturePad("z"))
	assert.Equal(t, FacilityFeatureMask(0), stringToFeaturePad("1sml"))
	assert.Equal(t, FacilityFeatureMask(0), stringToFeaturePad("Large"))
	assert.Equal(t, FacilityFeatureMask(0), stringToFeaturePad("large"))
	assert.Equal(t, FacilityFeatureMask(0), stringToFeaturePad("Mm"))
	assert.Equal(t, FacilityFeatureMask(0), stringToFeaturePad("Ss"))
}

// A scanner that injects a non-EOF to ensure we test the handling of scanner.Err()
type MockReader struct {
	strings.Reader
	err error
}

func (m *MockReader) Read(p []byte) (n int, err error) {
	// Just forward calls to Read() until we reach EOF, then if we have
	// an error in the MockReader, return that and remove it, so that the
	// next call will return the original error
	n, err = m.Reader.Read(p)
	if err == io.EOF && m.err != nil {
		err = m.err
		m.err = nil
	}
	return
}

func TestGetJsonRowsFromFile(t *testing.T) {
	var fieldNames = []string{"b", "a"}
	const jsonl = "{\"a\":\"1\",\"b\":2}\n{\"b\":3, \"a\"  : \"4\"}\nintentional error\n{\"c\":\"3po\",\"a\":\"5\",\"b\":6}\n{\"a\":1}\n"
	jsonReader := &MockReader{*strings.NewReader(jsonl), fmt.Errorf("scan error")}
	errorCh := make(chan error, 2)

	jsonLines := GetJsonItemsFromFile("jsonl.jsonl", jsonReader, fieldNames, errorCh)
	require.NotNil(t, jsonLines)
	// We should get 3 lines total.
	var allLines bool
	type ResultType struct {
		lineNo int
		intVal int64
		strVal string
	}
	var results []ResultType
	go func() {
		defer close(errorCh)
		for line := range jsonLines {
			lineNo := line.LineNo
			intVal := line.Results[0].Int()
			strVal := line.Results[1].String()
			results = append(results, ResultType{lineNo, intVal, strVal})
		}
		allLines = true
	}()

	var allErrors bool
	var errorList []string
	go func() {
		for err := range errorCh {
			errorList = append(errorList, err.Error())
		}
		allErrors = true
	}()

	assert.Eventually(t, func() bool { return allLines }, time.Millisecond*50, time.Microsecond*10)
	assert.Len(t, results, 3)
	assert.Equal(t, ResultType{1, 2, "1"}, results[0])
	assert.Equal(t, ResultType{2, 3, "4"}, results[1])
	assert.Equal(t, ResultType{4, 6, "5"}, results[2])

	assert.Eventually(t, func() bool { return allErrors }, time.Millisecond*50, time.Microsecond*10)
	assert.Len(t, errorList, 3)
	assert.Contains(t, errorList, "jsonl.jsonl:3: invalid json: intentional error")
	assert.Contains(t, errorList, "jsonl.jsonl:5: bad entry: {\"a\":1}")
	assert.Contains(t, errorList, "jsonl.jsonl:5: parse error: scan error")
}
