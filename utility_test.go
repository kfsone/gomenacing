package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestCaptureLog(t *testing.T) {
	str := captureLog(t, func(t *testing.T) {
		log.Print("!hello world!")
	})
	assert.Contains(t, str, "!hello world!")
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

func TestGenerateLinesFromReader(t *testing.T) {
	const testData = "testing\n1, 2, testing\n"

	// Create a fake reader
	reader := strings.NewReader(testData)
	generator := GenerateLinesFromReader(reader)
	require.NotNil(t, generator)

	// The generator should send us two lines and then end.
	var received string
	go func() {
		for line := range generator.OutputCh {
			received += line.(string) + "\n"
		}
	}()

	assert.Eventually(t, func() bool { return received == testData }, time.Millisecond*20, time.Microsecond*10)

	// Check that we can collect the error, and that it is nil
	var done bool
	go func() {
		done = generator.Error() == nil
	}()

	assert.Eventually(t, func() bool { return done }, time.Millisecond*10, time.Microsecond)
}

func TestIterateLinesInFile(t *testing.T) {
	const testData = "line 1\nline 2\n"
	var received string
	var err error

	file := strings.NewReader(testData + "X\n")

	// Run the iterator in the background in-case it deadlocks
	go func() {
		err = IterateLinesInFile("fakefile", file, func(line string) error {
			if line != "X" {
				received += line + "\n"
				return nil
			} else {
				return io.ErrClosedPipe
			}
		})
	}()

	assert.Eventually(t, func() bool { return err != nil }, time.Second, time.Microsecond)
	assert.True(t, errors.Is(err, io.ErrClosedPipe))
	assert.Equal(t, "fakefile:3: io: read/write on closed pipe", err.Error())
	assert.Equal(t, testData, received)
}
