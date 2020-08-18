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
	assert.NotNil(t, sdb.commoditiesById)
	assert.NotNil(t, sdb.commodityIds)
}

func TestSystemDatabase_registerCommodity(t *testing.T) {
	sdb := NewSystemDatabase()
	assert.Len(t, sdb.commoditiesById, 0)
	assert.Len(t, sdb.commodityIds, 0)

	item1 := Commodity{DbEntity: DbEntity{1, "first"}}
	assert.Nil(t, sdb.registerCommodity(&item1))
	assert.Equal(t, sdb.commoditiesById, map[EntityID]*Commodity{item1.Id: &item1})
	assert.Equal(t, sdb.commodityIds, map[string]EntityID{"first": item1.Id})

	err := sdb.registerCommodity(&item1)
	if assert.Error(t, err) {
		assert.Equal(t, "first (#1): duplicate: item id", err.Error())
	}
	assert.Equal(t, sdb.commoditiesById, map[EntityID]*Commodity{item1.Id: &item1})
	assert.Equal(t, sdb.commodityIds, map[string]EntityID{"first": item1.Id})

	item2 := Commodity{DbEntity: DbEntity{2, "first"}}
	err = sdb.registerCommodity(&item2)
	if assert.Error(t, err) {
		assert.Equal(t, "first (#2): duplicate: item name", err.Error())
	}
	assert.Equal(t, sdb.commoditiesById, map[EntityID]*Commodity{item1.Id: &item1})
	assert.Equal(t, sdb.commodityIds, map[string]EntityID{"first": item1.Id})

	item2.DbName = "second"
	assert.Nil(t, sdb.registerCommodity(&item2))
	assert.Equal(t, sdb.commoditiesById, map[EntityID]*Commodity{item1.Id: &item1, item2.Id: &item2})
	assert.Equal(t, sdb.commodityIds, map[string]EntityID{"first": item1.Id, "second": item2.Id})

	err = sdb.registerCommodity(&item2)
	if assert.Error(t, err) {
		assert.Equal(t, "second (#2): duplicate: item id", err.Error())
	}
	assert.Equal(t, sdb.commoditiesById, map[EntityID]*Commodity{item1.Id: &item1, item2.Id: &item2})
	assert.Equal(t, sdb.commodityIds, map[string]EntityID{"first": item1.Id, "second": item2.Id})
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

	err := sdb.registerFacility(&Facility{DbEntity: DbEntity{Id: 1, DbName: "first"}})
	if assert.Error(t, err) {
		assert.Equal(t, "first (#1): attempted to register facility without a system id", err.Error())
	}

	err = sdb.registerFacility(&Facility{DbEntity: DbEntity{Id: 1, DbName: "first"}, SystemId: 42})
	if assert.Error(t, err) {
		assert.Equal(t, "first (#1): system id: unknown: 42", err.Error())
	}

	facility1 := Facility{DbEntity: DbEntity{1, "first"}, System: &sys1}
	require.Nil(t, sdb.registerFacility(&facility1))
	assert.Len(t, sys1.Facilities, 1)
	assert.Contains(t, sys1.Facilities, &facility1)
	assert.Equal(t, sys1.Id, facility1.SystemId)
	assert.Equal(t, &sys1, facility1.System)

	// Registering the same id twice should fail.
	err = sdb.registerFacility(&facility1)
	if assert.Error(t, err) {
		assert.Equal(t, "first/first (#1): duplicate: facility id", err.Error())
	}
	assert.Equal(t, []*Facility{&facility1}, sys1.Facilities)
	assert.Empty(t, sys2.Facilities)
	assert.Equal(t, sys1.Id, facility1.SystemId)
	assert.Equal(t, &sys1, facility1.System)

	// Deliberately re-use id and name because they should be independent.
	facility2 := Facility{DbEntity: DbEntity{2, "first"}, System: &sys1}
	err = sdb.registerFacility(&facility2)
	if assert.Error(t, err) {
		assert.Equal(t, "first/first (#2): duplicate: facility name in system", err.Error())
	}
	assert.Equal(t, []*Facility{&facility1}, sys1.Facilities)
	assert.Empty(t, sys2.Facilities)
	assert.EqualValues(t, 0, facility2.SystemId)
	assert.Equal(t, &sys1, facility2.System)

	// But registering it under system 2 should work. Check that ID registration works.
	facility2.System = nil
	facility2.SystemId = sys2.Id
	require.Nil(t, sdb.registerFacility(&facility2))
	assert.Equal(t, &sys2, facility2.System)
	assert.Equal(t, sys2.Id, facility2.SystemId)
	assert.Equal(t, []*Facility{&facility1}, sys1.Facilities)
	assert.Equal(t, []*Facility{&facility2}, sys2.Facilities)

	// Registering either ID should fail at this point.
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
	assert.Equal(t, sys2.Id, facility2.SystemId)

	facility3 := Facility{DbEntity: DbEntity{3, "second"}, System: &sys1, SystemId: sys1.Id}
	require.Nil(t, sdb.registerFacility(&facility3))
	assert.Equal(t, []*Facility{&facility1, &facility3}, sys1.Facilities)
	assert.Equal(t, []*Facility{&facility2}, sys2.Facilities)
	assert.Equal(t, &sys1, facility3.System)
	assert.Equal(t, sys1.Id, facility3.SystemId)

	err = sdb.registerFacility(&facility3)
	if assert.Error(t, err) {
		assert.Equal(t, "first/second (#3): duplicate: facility id", err.Error())
	}
	assert.Equal(t, &sys1, facility3.System)
	assert.Equal(t, sys1.Id, facility3.SystemId)

	facility4 := Facility{DbEntity: DbEntity{4, "second"}, System: &sys1}
	err = sdb.registerFacility(&facility4)
	if assert.Error(t, err) {
		assert.Equal(t, "first/second (#4): duplicate: facility name in system", err.Error())
	}
	assert.Equal(t, &sys1, facility4.System)
	assert.EqualValues(t, 0, facility4.SystemId)

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
//	sdb := NewSystemDatabase()
//	require.NotNil(t, sdb)
//
//	require.Nil(t, sdb.importSystems())
//	require.Nil(t, sdb.importFacilities())
//}

func Test_countErrors(t *testing.T) {
	var err error

	// No errors, no count.
	errorCh := make(chan error, 8)
	close(errorCh)
	logged := captureLog(t, func(t *testing.T) {
		err = countErrors(errorCh)
	})
	assert.Nil(t, err)

	// Try a couple of errors.
	errorCh = make(chan error, 8)
	errorCh <- errors.New("[hello]")
	errorCh <- errors.New("[world]")
	close(errorCh)
	logged = captureLog(t, func(t *testing.T) {
		err = countErrors(errorCh)
	})
	assert.Error(t, err)
	assert.Equal(t, "failed because of 2 error(s)", err.Error())
	if assert.NotNil(t, logged) && assert.Len(t, logged, 2) {
		assert.Contains(t, logged[0], "[hello]")
		assert.Contains(t, logged[1], "[world]")
	}
}

func TestSystemDatabase_GetSystem(t *testing.T) {
	sdb := NewSystemDatabase()
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
