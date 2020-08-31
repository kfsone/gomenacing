package main

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
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
		assert.Equal(t, "first (#1): attempted to register facility without a system id", err.Error())
	}

	err = sdb.registerFacility(&Facility{DbEntity: DbEntity{ID: 1, DbName: "first"}, SystemID: 42})
	if assert.Error(t, err) {
		assert.Equal(t, "first (#1): system id: unknown: 42", err.Error())
	}

	facility1 := Facility{DbEntity: DbEntity{1, "first"}, System: &sys1}
	require.Nil(t, sdb.registerFacility(&facility1))
	assert.Len(t, sys1.Facilities, 1)
	assert.Contains(t, sys1.Facilities, &facility1)
	assert.Equal(t, sys1.ID, facility1.SystemID)
	assert.Equal(t, &sys1, facility1.System)

	// Registering the same id twice should fail.
	err = sdb.registerFacility(&facility1)
	if assert.Error(t, err) {
		assert.Equal(t, "first/first (#1): duplicate: facility id", err.Error())
	}
	assert.Equal(t, []*Facility{&facility1}, sys1.Facilities)
	assert.Empty(t, sys2.Facilities)
	assert.Equal(t, sys1.ID, facility1.SystemID)
	assert.Equal(t, &sys1, facility1.System)

	// Deliberately re-use id and name because they should be independent.
	facility2 := Facility{DbEntity: DbEntity{2, "first"}, System: &sys1}
	err = sdb.registerFacility(&facility2)
	if assert.Error(t, err) {
		assert.Equal(t, "first/first (#2): duplicate: facility name in system", err.Error())
	}
	assert.Equal(t, []*Facility{&facility1}, sys1.Facilities)
	assert.Empty(t, sys2.Facilities)
	assert.EqualValues(t, 0, facility2.SystemID)
	assert.Equal(t, &sys1, facility2.System)

	// But registering it under system 2 should work. Check that ID registration works.
	facility2.System = nil
	facility2.SystemID = sys2.ID
	require.Nil(t, sdb.registerFacility(&facility2))
	assert.Equal(t, &sys2, facility2.System)
	assert.Equal(t, sys2.ID, facility2.SystemID)
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
	assert.Equal(t, sys2.ID, facility2.SystemID)

	facility3 := Facility{DbEntity: DbEntity{3, "second"}, System: &sys1, SystemID: sys1.ID}
	require.Nil(t, sdb.registerFacility(&facility3))
	assert.Equal(t, []*Facility{&facility1, &facility3}, sys1.Facilities)
	assert.Equal(t, []*Facility{&facility2}, sys2.Facilities)
	assert.Equal(t, &sys1, facility3.System)
	assert.Equal(t, sys1.ID, facility3.SystemID)

	err = sdb.registerFacility(&facility3)
	if assert.Error(t, err) {
		assert.Equal(t, "first/second (#3): duplicate: facility id", err.Error())
	}
	assert.Equal(t, &sys1, facility3.System)
	assert.Equal(t, sys1.ID, facility3.SystemID)

	facility4 := Facility{DbEntity: DbEntity{4, "second"}, System: &sys1}
	err = sdb.registerFacility(&facility4)
	if assert.Error(t, err) {
		assert.Equal(t, "first/second (#4): duplicate: facility name in system", err.Error())
	}
	assert.Equal(t, &sys1, facility4.System)
	assert.EqualValues(t, 0, facility4.SystemID)

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

func TestImportJsonlFile(t *testing.T) {
	t.Run("Error on invalid file", func(t *testing.T) {
		c, err := ImportJsonlFile("testdata/DOESNOTEXIST", nil, nil)
		assert.Error(t, err)
		assert.Nil(t, c)
	})

	basename := filepath.Join("testdata", "TestImportJsonFile")
	fields := []string{"int", "str", "array", "e"}

	t.Run("Handle empty file", func(t *testing.T) {
		c, err := ImportJsonlFile(basename+".empty.jsonl", fields, func(entry *JSONEntry) error {
			return errors.New("unexpected row")
		})
		if assert.Nil(t, err) {
			assert.NotNil(t, c)
			// Check we see the channel close without sending anything
			assert.Eventually(t, func() bool {
				select {
				case _, ok := <-c:
					return !ok
				default:
					return false
				}
			}, time.Millisecond*20, time.Microsecond*10)
		}
	})

	t.Run("Basic parsing", func(t *testing.T) {
		var results = make([][]gjson.Result, 0, 4)
		c, err := ImportJsonlFile(basename+".basic.jsonl", fields, func(entry *JSONEntry) error {
			results = append(results, entry.Results)
			if entry.Results[3].Str != "" {
				return errors.New(entry.Results[3].String())
			}
			return nil
		})
		// Consume the errors channel
		errorsDrained := false
		var errs = make([]string, 0, 8)

		go func() {
			for {
				select {
				case err, ok := <-c:
					if ok {
						errs = append(errs, err.Error())
					} else {
						errorsDrained = true
						break
					}
				}
			}
		}()
		assert.Eventually(t, func() bool { return errorsDrained }, time.Millisecond*100, time.Microsecond*20)

		error1 := basename + `.basic.jsonl:2: bad entry: {"int": 2, "str": "two", "e":"", "ignored": "Missing array"}`
		error2 := basename + `.basic.jsonl:4: bad entry: {"int": 3, "array": [], "ignored": "Missing str and e"}`
		error3 := basename + `.basic.jsonl:5: expected error`
		error4 := basename + `.basic.jsonl:6: invalid json: invalid`
		assert.Len(t, errs, 4)
		assert.Contains(t, errs, error1)
		assert.Contains(t, errs, error2)
		assert.Contains(t, errs, error3)
		assert.Contains(t, errs, error4)

		if assert.Nil(t, err) {
			assert.NotNil(t, c)
			assert.Len(t, results, 3)
			// Check the types
			ints := []int64{1, 2, 3}
			strs := []string{"first", "second", "third"}
			for lineNo, result := range results {
				for idx, typeNo := range []int{2, 3, 5, 3} {
					assert.EqualValues(t, typeNo, result[idx].Type)
				}
				assert.Equal(t, ints[lineNo], result[0].Int())
				assert.Equal(t, strs[lineNo], result[1].String())
			}
		}
	})
}