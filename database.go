package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/akrylysov/pogreb"
	"golang.org/x/sync/errgroup"
	"path/filepath"
)

type Database struct {
	storePath string
}

func (db Database) Path() string {
	return db.storePath
}

func (db Database) Close() {
	// noop
}

func GetDatabase(path string) (*Database, error) {
	database := Database{storePath: filepath.Join(path, "database")}
	if _, err := ensureDirectory(database.Path()); err != nil {
		return nil, err
	}
	return &database, nil
}

func (db *Database) GetSchema(name string) (schema *Schema, err error) {
	path := filepath.Join(db.Path(), name)
	store, err := pogreb.Open(path, nil)
	if err != nil {
		return nil, err
	}
	schema = &Schema{db, name, store}
	return schema, nil
}

// Returns an open handle to the commodity schema
func (db *Database) Commodities() (*Schema, error) {
	return db.GetSchema("commodities")
}

// Returns an open handle to the facility schema
func (db *Database) Facilities() (*Schema, error) {
	return db.GetSchema("facilities")
}

func (db *Database) Listings() (*Schema, error) {
	return db.GetSchema("listings")
}

// Returns an open handle to the system schema
func (db *Database) Systems() (*Schema, error) {
	return db.GetSchema("systems")
}

func (db *Database) loadSystems(sdb *SystemDatabase) error {
	schema, err := db.Systems()
	if err == nil {
		sdb.systemIds = make(map[string]EntityID, schema.Count())
		sdb.systemsById = make(map[EntityID]*System, schema.Count())
		err = schema.LoadData("Systems", func(val []byte) error {
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
	schema, err := db.Facilities()
	if err == nil {
		sdb.facilitiesById = make(map[EntityID]*Facility, schema.Count())
		err = schema.LoadData("Facilities", func(val []byte) error {
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
	schema, err := db.Commodities()
	if err == nil {
		sdb.commodityIds = make(map[string]EntityID, schema.Count())
		sdb.commoditiesById = make(map[EntityID]*Commodity, schema.Count())
		err = schema.LoadData("Commodities", func(val []byte) error {
			var item = &Commodity{}
			if err = json.Unmarshal(val, item); err == nil {
				err = sdb.registerCommodity(item)
			}
			return err
		})
	}
	return err
}

func (db *Database) loadListings(sdb *SystemDatabase) error {
	schema, err := db.Listings()
	if err == nil {
		err = schema.LoadData("Listings", func(val []byte) error {
			_ = bytes.Split(val, []byte(" "))
			return nil
		})
	}
	return err
}

func (db *Database) LoadData(sdb *SystemDatabase) (err error) {
	///TODO: Speed up import-load by making import a part of load.
	/// Thus we can load commodities as soon as we've imported them.
	/// We could also break up checks into resource channels, so work
	/// can be done out of order.
	///  Load-from-file -> unmarshal -> conditions -> register.
	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)

	// Systems and facilities need to be loaded synchronously for now.
	///TODO: Evaluate making loadFacilities not associate facility->system
	/// Make sdb.facilitiesById be map[EntityID]*Facility
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
	return db.loadListings(sdb)
}
