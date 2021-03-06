package main

import (
	"context"
	"fmt"
	"github.com/akrylysov/pogreb"
	"github.com/kfsone/gomenacing/pkg/gomschema"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"log"
	"path/filepath"
)

type Database struct {
	storePath string
}

func OpenDatabase(path string, dbName string) (*Database, error) {
	database := Database{storePath: filepath.Join(path, dbName)}
	if _, err := ensureDirectory(database.Path()); err != nil {
		return nil, err
	}
	return &database, nil
}

func (db Database) Close() {
	// noop
}

func (db Database) Path() string {
	return db.storePath
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

func getSchemaForMessage(db *Database, message proto.Message) (*Schema, error) {
	switch v := message.(type) {
	case *gomschema.Commodity:
		return db.Commodities()

	case *gomschema.System:
		return db.Systems()

	case *gomschema.Facility:
		return db.Facilities()

	case *gomschema.FacilityListing:
		return db.Listings()

	default:
		return nil, fmt.Errorf("%w: message type: %t", ErrUnknownEntity, v)
	}
}

func (db *Database) loadSystems(sdb *SystemDatabase) error {
	schema, err := db.Systems()
	var loaded int
	if err == nil {
		sdb.systemIDs = make(map[string]EntityID, schema.Count())
		sdb.systemsByID = make(map[EntityID]*System, schema.Count())
		temporary := &gomschema.System{}
		loader, err := NewTypedDataLoader("gom", temporary, func() error { return sdb.newSystem(temporary) })
		if err == nil {
			loaded, err = schema.LoadData(loader)
			if err == nil {
				log.Printf("Loaded %d Systems.", loaded)
			}
		}
	}
	return err
}

func (db *Database) loadFacilities(sdb *SystemDatabase) error {
	schema, err := db.Facilities()
	var loaded int
	if err == nil {
		sdb.facilitiesByID = make(map[EntityID]*Facility, schema.Count())
		temporary := &gomschema.Facility{}
		loader, err := NewTypedDataLoader("gom", temporary, func() error { return sdb.newFacility(temporary) })
		if err == nil {
			loaded, err = schema.LoadData(loader)
			if err == nil {
				log.Printf("Loaded %d Facilities.", loaded)
			}
		}
	}
	return err
}

func (db *Database) loadCommodities(sdb *SystemDatabase) error {
	schema, err := db.Commodities()
	var loaded int
	if err == nil {
		sdb.commodityIDs = make(map[string]EntityID, schema.Count())
		sdb.commoditiesByID = make(map[EntityID]*Commodity, schema.Count())
		temporary := &gomschema.Commodity{}
		loader, err := NewTypedDataLoader("gom", temporary, func() error { return sdb.newCommodity(temporary) })
		if err == nil {
			loaded, err = schema.LoadData(loader)
			if err == nil {
				log.Printf("Loaded %d Commodities.", loaded)
			}
		}
	}
	return err
}

func (db *Database) loadListings(sdb *SystemDatabase) error {
	schema, err := db.Listings()
	loaded, lineItems := 0, 0
	if err == nil {
		temporary := &gomschema.FacilityListing{}
		loader, err := NewTypedDataLoader("gom", temporary, func() error {
			lineItems += len(temporary.Listings)
			return sdb.newListings(temporary)
		})
		if err == nil {
			loaded, err = schema.LoadData(loader)
			if err == nil {
				log.Printf("Loaded %d Listings at %d Facilities.", lineItems, loaded)
			}
		}
	}
	return err
}

func (db *Database) LoadDatabase(sdb *SystemDatabase) (err error) {
	///TODO: Speed up import-load by making import a part of load.
	/// Thus we can load commodities as soon as we've imported them.
	/// We could also break up checks into resource channels, so work
	/// can be done out of order.
	///  Load-from-file -> unmarshal -> conditions -> register.
	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)

	// Systems and facilities need to be loaded synchronously for now.
	///TODO: Evaluate making loadFacilities not associate facility->system
	/// Make sdb.facilitiesByID be map[EntityID]*Facility
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
