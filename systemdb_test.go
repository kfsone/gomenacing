package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSystemDatabase(t *testing.T) {
	sdb := NewSystemDatabase()
	assert.NotNil(t, sdb.systemsById)
	assert.NotNil(t, sdb.systemIds)
	assert.NotNil(t, sdb.facilitiesById)
}
