package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"path/filepath"

	"github.com/akrylysov/pogreb"
)

type Database struct {
	storePath string
}

func (db Database) Systems() (*pogreb.DB, error) {
	return pogreb.Open(filepath.Join(db.storePath, "systems"), nil)
}

func (db Database) Facilities() (*pogreb.DB, error) {
	return pogreb.Open(filepath.Join(db.storePath, "facilities"), nil)
}

func GetDatabase(path string) (*Database, error) {
	database := Database{storePath: filepath.Join(path, "database")}
	return &database, nil
}

func importWrapper(store *pogreb.DB, source string, fields []string, convertFn func([]gjson.Result) (interface{}, error)) error {
	defer failOnError(store.Close())
	filename := DataFilePath(*EddbPath, source)
	loaded := 0
	var data []byte

	errorsCh, err := ImportJsonFile(filename, fields, func(jsonLine *JsonLine) (err error) {
		item, err := convertFn(jsonLine.Results)
		if err == nil {
			if data, err = json.Marshal(item); err == nil {
				if err = store.Put([]byte(jsonLine.Results[0].Raw), data); err == nil {
					loaded += 1
					return nil
				}
			}
		}
		return fmt.Errorf("%s:%d: %w", source, jsonLine.LineNo, err)
	})
	if err == nil {
		err = countErrors(errorsCh)
	}
	if err != nil {
		return err
	}

	fmt.Printf("%s: imported %d items\n", source, loaded)
	return nil
}

func importSystems(db *Database) error {
	store, err := db.Systems()
	if err == nil {
		err = importWrapper(store, EddbSystems, systemFields, func(results []gjson.Result) (interface{}, error) { return NewSystemFromJson(results) })
	}
	return err
}

func importFacilities(db *Database) error {
	store, err := db.Facilities()
	if err == nil {
		err = importWrapper(store, EddbFacilities, facilityFields, func(results []gjson.Result) (interface{}, error) { return NewFacilityFromJson(results) })
	}
	return err
}

func loadData(name string, store *pogreb.DB, handler func(val []byte) error) error {
	defer failOnError(store.Close())

	it := store.Items()
	loaded := 0
	for {
		key, val, err := it.Next()
		if err != nil {
			if err == pogreb.ErrIterationDone {
				break
			}
			return err
		}
		if err := handler(val); err != nil {
			if FilterError(err) != nil {
				failOnError(store.Delete(key))
				return err
			}
		}
		loaded += 1
	}

	fmt.Printf("Loaded %d %s.\n", loaded, name)
	return nil
}

func (db *Database) loadSystems(sdb *SystemDatabase) error {
	store, err := db.Systems()
	if err == nil {
		sdb.systemIds = make(map[string]EntityID, store.Count())
		sdb.systemsById = make(map[EntityID]*System, store.Count())
		err = loadData("Systems", store, func(val []byte) error {
			var system = &System{}
			if err = json.Unmarshal(val, system); err == nil {
				err = sdb.registerSystem(system)
			}
			return err
		})
	}
	return err
}

func (db *Database) loadFacilities(sdb *SystemDatabase) error {
	store, err := db.Facilities()
	if err == nil {
		sdb.facilitiesById = make(map[EntityID]*Facility, store.Count())
		err = loadData("Facilities", store, func(val []byte) error {
			var facility = &Facility{}
			if err = json.Unmarshal(val, facility); err == nil {
				err = sdb.registerFacility(facility)
			}
			return err
		})
	}
	return err
}

func (db *Database) LoadData(sdb *SystemDatabase) (err error) {
	if err = db.loadSystems(sdb); err == nil {
		err = db.loadFacilities(sdb)
	}
	return err
}
