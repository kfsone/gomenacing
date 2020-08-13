package main

import (
	"encoding/json"
	"fmt"
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

func importSystems(db *Database) error {
	store, err := db.Systems()
	if err != nil {
		return err
	}
	defer store.Close()

	var system *System
	var data []byte
	loaded := 0
	filename := DataFilePath(*EddbPath, EddbSystems)

	errorsCh, err := ImportJsonFile(filename, systemFields, func(jsonLine JsonLine) (err error) {
		if system, err = NewSystemFromJson(jsonLine.Results); err == nil {
			if data, err = json.Marshal(*system); err == nil {
				if err = store.Put([]byte(jsonLine.Results[0].Raw), data); err == nil {
					loaded += 1
					return nil
				}
			}
		}
		return fmt.Errorf("%s:%d: %w", filename, jsonLine.LineNo, err)
	})
	if err == nil {
		err = countErrors(errorsCh)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Imported %d systems.\n", loaded)

	return nil
}

func importFacilities(db *Database) error {
	store, err := db.Facilities()
	if err != nil {
		return err
	}
	defer store.Close()

	filename := DataFilePath(*EddbPath, EddbFacilities)
	loaded := 0
	var facility *Facility
	var data []byte

	errorsCh, err := ImportJsonFile(filename, facilityFields, func(jsonLine JsonLine) (err error) {
		if facility, err = NewFacilityFromJson(jsonLine.Results); err == nil {
			if data, err = json.Marshal(*facility); err == nil {
				if err = store.Put([]byte(jsonLine.Results[0].Raw), data); err == nil {
					loaded += 1
					return nil
				}
			}
		}
		return fmt.Errorf("%s:%d: %w", filename, jsonLine.LineNo, err)
	})
	if err == nil {
		err = countErrors(errorsCh)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Imported %d facilities.\n", loaded)

	return nil
}

func (db *Database) loadSystems(sdb *SystemDatabase) error {
	store, err := db.Systems()
	if err != nil {
		return err
	}
	defer store.Close()

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
		var system = &System{}
		if err = json.Unmarshal(val, system); err != nil {
			store.Delete(key)
			return err
		}
		if err = sdb.registerSystem(system); err != nil {
			store.Delete(key)
			if FilterError(err) != nil {
				return err
			}
		}
		loaded += 1
	}
	fmt.Printf("Loaded %d systems\n", len(sdb.systemsById))
	return nil
}

func (db *Database) loadFacilities(sdb *SystemDatabase) error {
	store, err := db.Facilities()
	if err != nil {
		return err
	}
	defer store.Close()

	it := store.Items()
	loaded := 0
	for {
		_, val, err := it.Next()
		if err != nil {
			if err == pogreb.ErrIterationDone {
				break
			}
			return err
		}
		var facility = &Facility{}
		if err = json.Unmarshal(val, facility); err != nil {
			return err
		}
		if err = FilterError(sdb.registerFacility(facility)); err != nil {
			return err
		}
		loaded += 1
	}

	fmt.Printf("Loaded %d facilities\n", len(sdb.facilitiesById))

	return nil
}

func (db *Database) LoadData(sdb *SystemDatabase) (err error) {
	if err = db.loadSystems(sdb); err == nil {
		err = db.loadFacilities(sdb)
	}
	return err
}
