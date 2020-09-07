package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSystemDatabase(t *testing.T) {
	db := Database{}
	sdb := NewSystemDatabase(&db)
	assert.Equal(t, &db, sdb.db)
	assert.NotNil(t, sdb.systemsByID)
	assert.NotNil(t, sdb.systemIDs)
	assert.NotNil(t, sdb.facilitiesByID)
	assert.NotNil(t, sdb.commoditiesByID)
	assert.NotNil(t, sdb.commodityIDs)
}

func TestSystemDatabase_registerCommodity(t *testing.T) {
	sdb := NewSystemDatabase(nil)
	assert.Len(t, sdb.commoditiesByID, 0)
	assert.Len(t, sdb.commodityIDs, 0)

	item1 := Commodity{DbEntity: DbEntity{1, "first"}}
	assert.Nil(t, sdb.registerCommodity(&item1))
	assert.Equal(t, sdb.commoditiesByID, map[EntityID]*Commodity{item1.ID: &item1})
	assert.Equal(t, sdb.commodityIDs, map[string]EntityID{"first": item1.ID})

	err := sdb.registerCommodity(&item1)
	if assert.Error(t, err) {
		assert.Equal(t, "first (#1): duplicate: item id", err.Error())
	}
	assert.Equal(t, sdb.commoditiesByID, map[EntityID]*Commodity{item1.ID: &item1})
	assert.Equal(t, sdb.commodityIDs, map[string]EntityID{"first": item1.ID})

	item2 := Commodity{DbEntity: DbEntity{2, "first"}}
	err = sdb.registerCommodity(&item2)
	if assert.Error(t, err) {
		assert.Equal(t, "first (#2): duplicate: item name", err.Error())
	}
	assert.Equal(t, sdb.commoditiesByID, map[EntityID]*Commodity{item1.ID: &item1})
	assert.Equal(t, sdb.commodityIDs, map[string]EntityID{"first": item1.ID})

	item2.DbName = "second"
	assert.Nil(t, sdb.registerCommodity(&item2))
	assert.Equal(t, sdb.commoditiesByID, map[EntityID]*Commodity{item1.ID: &item1, item2.ID: &item2})
	assert.Equal(t, sdb.commodityIDs, map[string]EntityID{"first": item1.ID, "second": item2.ID})

	err = sdb.registerCommodity(&item2)
	if assert.Error(t, err) {
		assert.Equal(t, "second (#2): duplicate: item id", err.Error())
	}
	assert.Equal(t, sdb.commoditiesByID, map[EntityID]*Commodity{item1.ID: &item1, item2.ID: &item2})
	assert.Equal(t, sdb.commodityIDs, map[string]EntityID{"first": item1.ID, "second": item2.ID})
}

func TestSystemDatabase_registerSystem(t *testing.T) {
	var (
		err    error
		id     EntityID
		lookup *System
		ok     bool
	)
	sdb := NewSystemDatabase(nil)
	first := System{DbEntity: DbEntity{1, "first"}}
	second := System{DbEntity: DbEntity{2, "second"}}

	// In the simplest case, adding the first system should just work.
	require.Nil(t, sdb.registerSystem(&first))
	assert.Len(t, sdb.systemsByID, 1)
	assert.Len(t, sdb.systemIDs, 1)
	lookup, ok = sdb.systemsByID[1]
	assert.True(t, ok)
	assert.Equal(t, &first, lookup)
	id, ok = sdb.systemIDs["first"]
	assert.True(t, ok)
	assert.Equal(t, EntityID(1), id)

	// Adding a second system should also be fine
	require.Nil(t, sdb.registerSystem(&second))
	assert.Len(t, sdb.systemsByID, 2)
	assert.Len(t, sdb.systemIDs, 2)

	// Check the first system is still correct
	lookup, ok = sdb.systemsByID[1]
	assert.True(t, ok)
	assert.Equal(t, &first, lookup)
	id, ok = sdb.systemIDs["first"]
	assert.True(t, ok)
	assert.Equal(t, EntityID(1), id)

	// Check the second system is registered correctly.
	lookup, ok = sdb.systemsByID[2]
	assert.True(t, ok)
	assert.Equal(t, &second, lookup)
	id, ok = sdb.systemIDs["second"]
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
	sdb := NewSystemDatabase(nil)

	// Register two star systems.
	sys1 := System{DbEntity: DbEntity{1, "first"}}
	require.Nil(t, sdb.registerSystem(&sys1))
	sys2 := System{DbEntity: DbEntity{2, "second"}}
	require.Nil(t, sdb.registerSystem(&sys2))

	err := sdb.registerFacility(&Facility{DbEntity: DbEntity{ID: 1, DbName: "first"}})
	if assert.Error(t, err) {
		assert.Equal(t, "first (#1): attempted to register facility without a system", err.Error())
	}

	facility1 := Facility{DbEntity: DbEntity{1, "first"}, System: &sys1}
	require.Nil(t, sdb.registerFacility(&facility1))
	assert.Len(t, sys1.facilities, 1)
	assert.Contains(t, sys1.facilities, &facility1)
	assert.Equal(t, &sys1, facility1.System)

	// Registering the same id twice should fail.
	err = sdb.registerFacility(&facility1)
	if assert.Error(t, err) {
		assert.Equal(t, "first/first (#1): duplicate: facility id", err.Error())
	}
	assert.Equal(t, []*Facility{&facility1}, sys1.facilities)
	assert.Empty(t, sys2.facilities)
	assert.Equal(t, &sys1, facility1.System)

	// Deliberately re-use id and name because they should be independent.
	facility2 := Facility{DbEntity: DbEntity{2, "first"}, System: &sys1}
	err = sdb.registerFacility(&facility2)
	if assert.Error(t, err) {
		assert.Equal(t, "first/first (#2): duplicate: facility name in system", err.Error())
	}
	assert.Equal(t, []*Facility{&facility1}, sys1.facilities)
	assert.Empty(t, sys2.facilities)
	assert.Equal(t, &sys1, facility2.System)

	// But registering it under system 2 should work.
	facility2.System = &sys2
	require.Nil(t, sdb.registerFacility(&facility2))
	assert.Equal(t, &sys2, facility2.System)
	assert.Equal(t, []*Facility{&facility1}, sys1.facilities)
	assert.Equal(t, []*Facility{&facility2}, sys2.facilities)

	// Registering either should fail at this point.
	err = sdb.registerFacility(&facility1)
	if assert.Error(t, err) {
		assert.Equal(t, "first/first (#1): duplicate: facility id", err.Error())
	}
	err = sdb.registerFacility(&facility2)
	facility2.System = nil
	if assert.Error(t, err) {
		assert.Equal(t, "second/first (#2): duplicate: facility id", err.Error())
	}
	assert.Nil(t, facility2.System)

	facility3 := Facility{DbEntity: DbEntity{3, "second"}, System: &sys1}
	require.Nil(t, sdb.registerFacility(&facility3))
	assert.Equal(t, []*Facility{&facility1, &facility3}, sys1.facilities)
	assert.Equal(t, []*Facility{&facility2}, sys2.facilities)
	assert.Equal(t, &sys1, facility3.System)

	err = sdb.registerFacility(&facility3)
	if assert.Error(t, err) {
		assert.Equal(t, "first/second (#3): duplicate: facility id", err.Error())
	}
	assert.Equal(t, &sys1, facility3.System)

	facility4 := Facility{DbEntity: DbEntity{4, "second"}, System: &sys1}
	err = sdb.registerFacility(&facility4)
	if assert.Error(t, err) {
		assert.Equal(t, "first/second (#4): duplicate: facility name in system", err.Error())
	}
	assert.Equal(t, &sys1, facility4.System)

	facility4.System = &sys2
	require.Nil(t, sdb.registerFacility(&facility4))
	assert.Equal(t, []*Facility{&facility1, &facility3}, sys1.facilities)
	assert.Equal(t, []*Facility{&facility2, &facility4}, sys2.facilities)

	err = sdb.registerFacility(&facility4)
	if assert.Error(t, err) {
		assert.Equal(t, "second/second (#4): duplicate: facility id", err.Error())
	}
	assert.Equal(t, []*Facility{&facility1, &facility3}, sys1.facilities)
	assert.Equal(t, []*Facility{&facility2, &facility4}, sys2.facilities)
}

func TestSystemDatabase_GetSystem(t *testing.T) {
	sdb := NewSystemDatabase(nil)
	t.Run("Query v empty db", func(t *testing.T) {
		assert.Nil(t, sdb.GetSystem(""))
		assert.Nil(t, sdb.GetSystem("first"))
	})

	first := System{DbEntity: DbEntity{1, "first"}}
	err := sdb.registerSystem(&first)
	require.Nil(t, err)

	second := System{DbEntity: DbEntity{2, "second"}}
	err = sdb.registerSystem(&second)
	require.Nil(t, err)

	assert.Nil(t, sdb.GetSystem(""))
	assert.Nil(t, sdb.GetSystem("invalid"))
	assert.Equal(t, &first, sdb.GetSystem("first"))
	assert.Equal(t, &first, sdb.GetSystem("fIrsT"))
	assert.Equal(t, &first, sdb.GetSystem("FIRST"))
	assert.Nil(t, sdb.GetSystem("third"))
	assert.Equal(t, &second, sdb.GetSystem("second"))
	assert.Equal(t, &second, sdb.GetSystem("SECOND"))
	assert.Nil(t, sdb.GetSystem("Firsts"))
	assert.Nil(t, sdb.GetSystem("Second."))
}

func Test_registerIDLookup(t *testing.T) {
	ids := make(map[string]EntityID)
	expected := make(map[string]EntityID)
	// We don't do validation, so registering an empty string can occur.
	require.True(t, registerIDLookup(&DbEntity{ID: 0, DbName: ""}, ids))
	expected[""] = 0
	assert.Equal(t, expected, ids)

	// It should fail if we try to do it again.
	require.False(t, registerIDLookup(&DbEntity{ID: 0, DbName: ""}, ids))
	assert.Equal(t, map[string]EntityID{"": EntityID(0)}, ids)
	assert.Equal(t, expected, ids)

	// Register something different, but it shouldn't do validation, so reuse id.
	require.True(t, registerIDLookup(&DbEntity{ID: 0, DbName: "Number 1"}, ids))
	// Note: registered IDs should have lowercase.
	expected["number 1"] = 0
	assert.Equal(t, expected, ids)

	// And register something completely different.
	require.True(t, registerIDLookup(&DbEntity{ID: 42, DbName: "Life, the Universe, EVERYTHING"}, ids))
	// again, lowercase IDs:
	expected["life, the universe, everything"] = 42
	assert.Equal(t, expected, ids)
}

func TestSystemDatabase_GetSystemByID(t *testing.T) {
	sdb := SystemDatabase{}
	assert.Nil(t, sdb.GetSystemByID(0))
	assert.Nil(t, sdb.GetSystemByID(1))
	assert.Nil(t, sdb.GetSystemByID(42))

	sdb.systemsByID = make(map[EntityID]*System)
	assert.Nil(t, sdb.GetSystemByID(0))
	assert.Nil(t, sdb.GetSystemByID(1))
	assert.Nil(t, sdb.GetSystemByID(42))

	// Confirm that querying did not alter the table.
	assert.Len(t, sdb.systemsByID, 0)

	system0 := System{}
	system7 := System{}
	sdb.systemsByID[7] = &system7
	sdb.systemsByID[0] = &system0
	assert.Nil(t, sdb.GetSystemByID(1))
	assert.Equal(t, &system0, sdb.GetSystemByID(0))
	assert.Equal(t, &system7, sdb.GetSystemByID(7))

	// Confirm that querying did not alter the table.
	assert.Len(t, sdb.systemsByID, 2)
}
