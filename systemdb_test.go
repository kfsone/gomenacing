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
	first := System{DbEntity: DbEntity{1, "first"}}
	second := System{DbEntity: DbEntity{2, "second"}}

	// In the simplest case, adding the first system should just work.
	require.Nil(t, sdb.registerSystem(&first))
	assert.Len(t, sdb.systemsById, 1)
	assert.Len(t, sdb.systemIds, 1)
	lookup, ok = sdb.systemsById[1]
	assert.True(t, ok)
	assert.Equal(t, &first, lookup)
	id, ok = sdb.systemIds["first"]
	assert.True(t, ok)
	assert.Equal(t, EntityID(1), id)

	// Adding a second system should also be fine
	require.Nil(t, sdb.registerSystem(&second))
	assert.Len(t, sdb.systemsById, 2)
	assert.Len(t, sdb.systemIds, 2)

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
	if assert.True(t, errors.Is(err, ErrDuplicateEntity)) {
		assert.Equal(t, "first (#1): duplicate: system id", err.Error())
	}

	err = sdb.registerSystem(&System{DbEntity: DbEntity{3, "first"}})
	if assert.Error(t, err) {
		assert.Equal(t, "first (#3): duplicate: system name", err.Error())
	}
}

func TestSystemDatabase_registerFacility(t *testing.T) {
	sdb := NewSystemDatabase()

	// Register two star systems.
	sys1 := System{DbEntity: DbEntity{1, "first"}}
	require.Nil(t, sdb.registerSystem(&sys1))
	sys2 := System{DbEntity: DbEntity{2, "second"}}
	require.Nil(t, sdb.registerSystem(&sys2))

	// Check that a missing system is validated.
	err := sdb.registerFacility(&Facility{})
	if assert.Error(t, err) {
		assert.Equal(t, "attempted to register facility with nil system", err.Error())
	}

	facility1 := Facility{DbEntity: DbEntity{1, "first"}, System: &sys1}
	require.Nil(t, sdb.registerFacility(&facility1))
	assert.Len(t, sys1.Facilities, 1)
	assert.Contains(t, sys1.Facilities, &facility1)

	// Registering the same id twice should fail.
	err = sdb.registerFacility(&facility1)
	if assert.Error(t, err) {
		assert.Equal(t, "first/first (#1): duplicate: facility id", err.Error())
	}
	assert.Equal(t, []*Facility{&facility1}, sys1.Facilities)
	assert.Empty(t, sys2.Facilities)

	// Deliberately re-use id and name because they should be independent.
	facility2 := Facility{DbEntity: DbEntity{2, "first"}, System: &sys1}
	err = sdb.registerFacility(&facility2)
	if assert.Error(t, err) {
		assert.Equal(t, "first/first (#2): duplicate: facility name in system", err.Error())
	}
	assert.Equal(t, []*Facility{&facility1}, sys1.Facilities)
	assert.Empty(t, sys2.Facilities)

	// But registering it under system 2 should work.
	facility2.System = &sys2
	require.Nil(t, sdb.registerFacility(&facility2))
	assert.Equal(t, []*Facility{&facility1}, sys1.Facilities)
	assert.Equal(t, []*Facility{&facility2}, sys2.Facilities)

	// Registering either ID should fail at this point.
	err = sdb.registerFacility(&facility1)
	if assert.Error(t, err) {
		assert.Equal(t, "first/first (#1): duplicate: facility id", err.Error())
	}
	err = sdb.registerFacility(&facility2)
	if assert.Error(t, err) {
		assert.Equal(t, "second/first (#2): duplicate: facility id", err.Error())
	}

	facility3 := Facility{DbEntity: DbEntity{3, "second"}, System: &sys1}
	require.Nil(t, sdb.registerFacility(&facility3))
	assert.Equal(t, []*Facility{&facility1, &facility3}, sys1.Facilities)
	assert.Equal(t, []*Facility{&facility2}, sys2.Facilities)

	err = sdb.registerFacility(&facility3)
	if assert.Error(t, err) {
		assert.Equal(t, "first/second (#3): duplicate: facility id", err.Error())
	}

	facility4 := Facility{DbEntity: DbEntity{4, "second"}, System: &sys1}
	err = sdb.registerFacility(&facility4)
	if assert.Error(t, err) {
		assert.Equal(t, "first/second (#4): duplicate: facility name in system", err.Error())
	}

	facility4.System = &sys2
	require.Nil(t, sdb.registerFacility(&facility4))
	assert.Equal(t, []*Facility{&facility1, &facility3}, sys1.Facilities)
	assert.Equal(t, []*Facility{&facility2, &facility4}, sys2.Facilities)

	err = sdb.registerFacility(&facility4)
	if assert.Error(t, err) {
		assert.Equal(t, "second/second (#4): duplicate: facility id", err.Error())
	}
	assert.Equal(t, []*Facility{&facility1, &facility3}, sys1.Facilities)
	assert.Equal(t, []*Facility{&facility2, &facility4}, sys2.Facilities)
}

//func TestSystemDatabase_importFacilities(t *testing.T) {
//	env, err := NewEnv("", "")
//	require.Nil(t, err)
//	require.NotNil(t, env)
//
//	env.SilenceWarnings = true
//
//	sdb := NewSystemDatabase()
//	require.NotNil(t, sdb)
//
//	require.Nil(t, sdb.importSystems(env))
//	require.Nil(t, sdb.importFacilities(env))
//}

func Test_countErrors(t *testing.T) {
	var err error
	env := Env{}
	const filename = "test.me"

	// No errors, no count.
	errorCh := make(chan error, 8)
	close(errorCh)
	logged := captureLog(t, func(t *testing.T) {
		err = countErrors(&env, filename, errorCh)
	})
	assert.Nil(t, err)

	// Try a couple of errors.
	errorCh = make(chan error, 8)
	errorCh <- errors.New("[hello]")
	errorCh <- errors.New("[world]")
	close(errorCh)
	logged = captureLog(t, func(t *testing.T) {
		err = countErrors(&env, filename, errorCh)
	})
	assert.Error(t, err)
	assert.Equal(t, "failed because of 2 error(s)", err.Error())
	if assert.NotNil(t, logged) && assert.Len(t, logged, 2) {
		assert.Contains(t, logged[0], "[hello]")
		assert.Contains(t, logged[1], "[world]")
	}
}
