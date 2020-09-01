package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
	testDir := GetTestDir()
	defer testDir.Close()
	tmpPath := testDir.Path()

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

func TestTestDir(t *testing.T) {
	t.Run("Panic on illegal path", func(t *testing.T) {
		require.Panics(t, func() {
			GetTestDir("/c:/invalid/\\never")
		})
	})

	t.Run("Validate GetTestDir", func(t *testing.T) {
		testDir := GetTestDir()
		defer testDir.Close()
		if assert.True(t, strings.HasPrefix(testDir.Path(), filepath.Join(os.TempDir(), "menace"))) {
			if assert.DirExists(t, testDir.Path()) {
				assert.NotPanics(t, func() { testDir.Close() })
				assert.NoDirExists(t, testDir.Path())
				assert.NotPanics(t, func() { testDir.Close() })
			}
		}
	})

	t.Run("Validate Close error handling", func(t *testing.T) {
		// In order to test the handling of removeall throwing an error, we need to
		// get underhanded. On Windows, having a file open in a directory will help
		// us block deletion. On Linux we need to mess with permissions. This is a
		// combination of both.
		testDir := GetTestDir()
		defer testDir.Close()

		innerName := filepath.Join(testDir.Path(), "block")
		file, err := os.OpenFile(innerName, os.O_RDONLY|os.O_CREATE|os.O_EXCL, 0400)
		require.Nil(t, err)

		assert.Nil(t, os.Chmod(testDir.Path(), 0400))
		savedPath := testDir.Path()
		testDir = TestDir(filepath.Join(savedPath, "."))
		assert.Panics(t, func() { testDir.Close() })
		testDir = TestDir(savedPath)

		file.Close()

		assert.Nil(t, os.Chmod(savedPath, 0700))
		assert.Nil(t, os.Chmod(innerName, 0700))

		assert.NotPanics(t, func() { testDir.Close() })

		assert.NoDirExists(t, testDir.Path())
	})
}
