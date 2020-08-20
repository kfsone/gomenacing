package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTestDir_Path(t *testing.T) {
	assert.Equal(t, "foo/bar", TestDir("foo/bar").Path())
}

func TestTestDir_String(t *testing.T) {
	assert.Equal(t, "foo/bar", TestDir("foo/bar").String())
}
