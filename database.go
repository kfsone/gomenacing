package main

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/akrylysov/pogreb"
	"golang.org/x/sync/errgroup"
)

type Database struct {
	storePath string
}

// Helper that opens a specific pogreb schema
func getSchema(path string, name string) (*pogreb.DB, error) {
	return pogreb.Open(filepath.Join(path, name), nil)
}

// Returns an open handle to the commodity schema
func (db Database) Commodities() (*pogreb.DB, error) {
	return getSchema(db.storePath, "commodities")
}

// Returns an open handle to the facility schema
func (db Database) Facilities() (*pogreb.DB, error) {
	return getSchema(db.storePath, "facilities")
}

// Returns an open handle to the system schema
func (db Database) Systems() (*pogreb.DB, error) {
	return getSchema(db.storePath, "systems")
}

func GetDatabase(path string) (*Database, error) {
	database := Database{storePath: filepath.Join(path, "database")}
	return &database, nil
}

func loadData(name string, store *pogreb.DB, handler func(val []byte) error) error {
	defer func() { failOnError(store.Close()) }()

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

func (db *Database) loadCommodities(sdb *SystemDatabase) error {
	store, err := db.Commodities()
	if err == nil {
		sdb.commodityIds = make(map[string]EntityID, store.Count())
		sdb.commoditiesById = make(map[EntityID]*Commodity, store.Count())
		err = loadData("Commodities", store, func(val []byte) error {
			var item = &Commodity{}
			if err = json.Unmarshal(val, item); err == nil {
				err = sdb.registerCommodity(item)
			}
			return err
		})
	}
	return err
}

func (db *Database) LoadData(sdb *SystemDatabase) (err error) {
	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)

	// Systems and facilities need to be loaded synchronously for now.
	///TODO: Evaluate making loadFacilities not associate facility->system
	eg.Go(func() error {
		if err = db.loadSystems(sdb); err == nil {
			err = db.loadFacilities(sdb)
		}
		return err
	})
	// We can load the commodity list in parallel
	eg.Go(func() error {
		return db.loadCommodities(sdb)
	})
	if err := eg.Wait(); err != nil {
		return err
	}

	// Now import any prices.
	return nil
}
