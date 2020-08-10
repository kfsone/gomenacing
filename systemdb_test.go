package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSystemDatabase(t *testing.T) {
	sdb := NewSystemDatabase()
	assert.NotNil(t, sdb.systemsById)
	assert.NotNil(t, sdb.systemIds)
	assert.NotNil(t, sdb.facilitiesById)
}

func TestSystemDatabase_registerSystem(t *testing.T) {
	var (
		err    error
		id     EntityID
		lookup *System
		ok     bool
	)
	sdb := NewSystemDatabase()
	first := System{DatabaseEntity: DatabaseEntity{1, "first"}}
	second := System{DatabaseEntity: DatabaseEntity{2, "second"}}

	// In the simplest case, adding the first system should just work.
	require.Nil(t, sdb.registerSystem(&first))
	assert.Equal(t, 1, len(sdb.systemsById))
	assert.Equal(t, 1, len(sdb.systemIds))
	lookup, ok = sdb.systemsById[1]
	assert.True(t, ok)
	assert.Equal(t, &first, lookup)
	id, ok = sdb.systemIds["first"]
	assert.True(t, ok)
	assert.Equal(t, EntityID(1), id)

	// Adding a second system should also be fine
	require.Nil(t, sdb.registerSystem(&second))
	assert.Equal(t, 2, len(sdb.systemsById))
	assert.Equal(t, 2, len(sdb.systemIds))

	// Check the first system is still correct
	lookup, ok = sdb.systemsById[1]
	assert.True(t, ok)
	assert.Equal(t, &first, lookup)
	id, ok = sdb.systemIds["first"]
	assert.True(t, ok)
	assert.Equal(t, EntityID(1), id)

	// Check the second system is registered correctly.
	lookup, ok = sdb.systemsById[2]
	assert.True(t, ok)
	assert.Equal(t, &second, lookup)
	id, ok = sdb.systemIds["second"]
	assert.True(t, ok)
	assert.Equal(t, EntityID(2), id)

	// Trying to register first again should cause this to fail.
	err = sdb.registerSystem(&first)
	assert.True(t, errors.Is(err, ErrDuplicateEntity))
	assert.Error(t, err, "FIRST (#1): duplicate: system id")

	err = sdb.registerSystem(&System{DatabaseEntity: DatabaseEntity{3, "first"}})
	assert.Error(t, err, "FIRST (#3): duplicate: system name")
}
