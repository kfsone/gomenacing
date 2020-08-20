package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataFile(t *testing.T) {
	assert.Equal(t, "menace.db", DbFile())
}

func TestDataPath(t *testing.T) {
	assert.Equal(t, ".", DataPath())
}

func TestSetupEnv(t *testing.T) {
	assert.Equal(t, ".", *DefaultPath)
	assert.Equal(t, "menace.db", *DefaultDbFile)

	testDir := GetTestDir()
	defer testDir.Close()

	oldDefault := *DefaultPath
	defer func() { *DefaultPath = oldDefault }()

	*DefaultPath = filepath.Join(testDir.Path(), "upper", "lower", "dir")

	err := SetupEnv()
	assert.Nil(t, err)
	_, err = os.Stat(*DefaultPath)
	assert.Nil(t, err)

	// It shouldn't care if we do it again.
	err = SetupEnv()
	assert.Nil(t, err)

	// Check for an error if the specified data directory is a file.
	*DefaultPath = filepath.Join(testDir.Path(), "collide.file")
	file, err := os.Create(*DefaultPath)
	require.Nil(t, err)
	assert.Nil(t, file.Close())

	err = SetupEnv()
	if assert.Error(t, err) {
		assert.True(t, os.IsExist(err))
	}
}

func TestFilterError(t *testing.T) {
	// nil error should return nil with no output.
	result := captureLog(t, func(t *testing.T) {
		assert.Nil(t, FilterError(nil))
	})
	assert.Empty(t, result)

	// MissingCategory shouldn't be filtered.
	result = captureLog(t, func(t *testing.T) {
		assert.Equal(t, MissingCategory, FilterError(MissingCategory))
	})
	assert.Empty(t, result)

	t.Run("Check Duplicate Entry errors", func(t *testing.T) {
		// ErrDuplicateEntity should be filtered until we make it an error,
		// which means nil error returned and no logging generated, yet.
		result = captureLog(t, func(t *testing.T) {
			assert.Nil(t, FilterError(fmt.Errorf("test: %w", ErrDuplicateEntity)))
		})
		assert.Nil(t, result)

		// enabling duplicationErrors should change this to just returning an error/no log
		defer func() { *ErrorOnDuplicate = false }()
		*ErrorOnDuplicate = true
		result = captureLog(t, func(t *testing.T) {
			err := FilterError(fmt.Errorf("test: %w", ErrDuplicateEntity))
			assert.True(t, errors.Is(err, ErrDuplicateEntity))
		})
		assert.Nil(t, result)
	})

	t.Run("Check Unknown Entity errors", func(t *testing.T) {
		// ErrUnknownEntity should likewise be filtered/logged.
		result = captureLog(t, func(t *testing.T) {
			assert.Nil(t, FilterError(fmt.Errorf("test: %w", ErrUnknownEntity)))
		})
		assert.Nil(t, result)

		defer func() { *ErrorOnUnknown = false }()
		*ErrorOnUnknown = true
		result = captureLog(t, func(t *testing.T) {
			err := FilterError(fmt.Errorf("test: %w", ErrUnknownEntity))
			assert.True(t, errors.Is(err, ErrUnknownEntity))
		})
		assert.Nil(t, result)
	})

	t.Run("Check ShowWarnings", func(t *testing.T) {
		// Turning onShowWarnings should get nils but no outputs
		defer func() { *ShowWarnings = false }()
		*ShowWarnings = true
		result = captureLog(t, func(t *testing.T) {
			assert.Nil(t, FilterError(fmt.Errorf("test: %w", ErrDuplicateEntity)))
			assert.Nil(t, FilterError(fmt.Errorf("test: %w", ErrUnknownEntity)))
		})
		if assert.NotNil(t, result) {
			assert.Len(t, result, 2)
			assert.Contains(t, result[0], "test: duplicate")
			assert.Contains(t, result[1], "test: unknown")
		}
	})
}

func TestDataFilePath(t *testing.T) {
	var oldDefaultPath = *DefaultPath
	defer func() { *DefaultPath = oldDefaultPath }()
	*DefaultPath = "data/test"

	assert.Equal(t, filepath.Join("data/test", "foo.txt"), DataFilePath("foo.txt"))
	assert.Equal(t, filepath.Join("data/test", "foo", "bar", "baz.txt"), DataFilePath("foo", "bar", "baz.txt"))
}
