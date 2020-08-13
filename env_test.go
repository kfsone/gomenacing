package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnv_DataFile(t *testing.T) {
	env := Env{dataFile: "hello/world"}
	assert.Equal(t, "hello/world", env.DataFile())
}

func TestEnv_DataPath(t *testing.T) {
	env := Env{dataPath: "c:/windows/user"}
	assert.Equal(t, "c:/windows/user", env.DataPath())
}

func TestNewEnv(t *testing.T) {
	// Verify that we can create an env for an existing directory.
	tempDir, err := ioutil.TempDir("", "menacing-test")
	assert.Nil(t, err)
	assert.NotEmpty(t, tempDir)
	defer os.RemoveAll(tempDir)

	oldDefault := DefaultDbPath
	DefaultDbPath = filepath.Join(tempDir, "data")
	defer func() { DefaultDbPath = oldDefault }()

	env, err := NewEnv("", "")
	require.Nil(t, err)
	require.NotNil(t, env)
	assert.Equal(t, filepath.Join(tempDir, "data"), env.DataPath())
	assert.Equal(t, filepath.Join(tempDir, "data", DefaultDbFile), env.DataFile())
	assert.False(t, env.ErrorOnDuplicate)
	assert.False(t, env.ErrorOnUnknown)
	_, err = os.Stat(DefaultDbPath)
	require.Nil(t, err)

	// We should be able to create another without issue.
	env2, err := NewEnv("", "")
	require.Nil(t, err)
	require.NotNil(t, env2)

	// If we specify arguments, they should override values.
	env, err = NewEnv(filepath.Join(tempDir, "second"), "fiddle.db")
	require.Nil(t, err)
	require.NotNil(t, env)
	assert.Equal(t, filepath.Join(tempDir, "second"), env.DataPath())
	assert.Equal(t, filepath.Join(env.DataPath(), "fiddle.db"), env.DataFile())
	_, err = os.Stat(env.DataPath())
	assert.Nil(t, err)

	// Check for an error if the specified data directory is a file.
	file, err := os.Create(env.DataFile())
	require.Nil(t, err)
	assert.Nil(t, file.Close())
	env, err = NewEnv(env.DataFile(), "foo.bar")
	assert.Nil(t, env)
	assert.True(t, os.IsExist(err))
}

func TestFilterError(t *testing.T) {
	// nil error should return nil with no output.
	env := Env{}
	result := captureLog(t, func(t *testing.T) {
		assert.Nil(t, env.FilterError(nil))
	})
	assert.Empty(t, result)

	// MissingCategory shouldn't be filtered.
	result = captureLog(t, func(t *testing.T) {
		assert.Equal(t, MissingCategory, env.FilterError(MissingCategory))
	})
	assert.Empty(t, result)

	// ErrDuplicateEntity should be filtered until we make it an error,
	// which means nil error returned but logging generated.
	result = captureLog(t, func(t *testing.T) {
		assert.Nil(t, env.FilterError(fmt.Errorf("test: %w", ErrDuplicateEntity)))
	})
	if assert.NotNil(t, result) {
		assert.Len(t, result, 1)
		assert.Contains(t, result[0], "test: duplicate")
	}

	// enabling duplicationErrors should change this to just returning an error/no log
	env.ErrorOnDuplicate = true
	result = captureLog(t, func(t *testing.T) {
		err := env.FilterError(fmt.Errorf("test: %w", ErrDuplicateEntity))
		assert.True(t, errors.Is(err, ErrDuplicateEntity))
	})
	assert.Nil(t, result)

	// ErrUnknownEntity should likewise be filtered/logged.
	result = captureLog(t, func(t *testing.T) {
		assert.Nil(t, env.FilterError(fmt.Errorf("test: %w", ErrUnknownEntity)))
	})
	if assert.NotNil(t, result) {
		assert.Len(t, result, 1)
		assert.Contains(t, result[0], "test: unknown")
	}
	env.ErrorOnUnknown = true
	result = captureLog(t, func(t *testing.T) {
		err := env.FilterError(fmt.Errorf("test: %w", ErrUnknownEntity))
		assert.True(t, errors.Is(err, ErrUnknownEntity))
	})
	assert.Nil(t, result)

	// Turning on SilenceWarnings should get nils but no outputs
	env = Env{SilenceWarnings: true}
	result = captureLog(t, func(t *testing.T) {
		assert.Nil(t, env.FilterError(fmt.Errorf("test: %w", ErrDuplicateEntity)))
		assert.Nil(t, env.FilterError(fmt.Errorf("test: %w", ErrUnknownEntity)))
	})
	assert.Nil(t, result)
}
