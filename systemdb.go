package main

import (
	"fmt"
	"strings"

	flag "github.com/spf13/pflag"
)

var eddbPath = flag.StringP("eddbdir", "e", "", "Path to EDDB json files to import.")

const (
	EddbSystems     string = "systems_populated.jsonl"
	EddbFacilities  string = "stations.jsonl"
	EddbCommodities string = "commodities.json"
)

type SystemDatabase struct {
	db *Database
	// Index of Systems by their database ids.
	systemsByID map[EntityID]*System
	// Look-up a system's EntityID by it's name.
	systemIDs map[string]EntityID
	// Index of Facilities by their database ids.
	facilitiesByID map[EntityID]*Facility
	// Index of Commodities by their database ids.
	commoditiesByID map[EntityID]*Commodity
	// Look-up a commodity's EntityID by it's name.
	commodityIDs map[string]EntityID
}

func NewSystemDatabase(db *Database) *SystemDatabase {
	return &SystemDatabase{
		db:              db,
		systemsByID:     make(map[EntityID]*System),
		systemIDs:       make(map[string]EntityID),
		facilitiesByID:  make(map[EntityID]*Facility),
		commoditiesByID: make(map[EntityID]*Commodity),
		commodityIDs:    make(map[string]EntityID),
	}
}

func registerIDLookup(entity DbEntity, ids map[string]EntityID) bool {
	name := strings.ToLower(entity.DbName)
	if _, present := ids[name]; present != false {
		return false
	}
	ids[name] = entity.ID
	return true
}

func (sdb *SystemDatabase) registerCommodity(commodity *Commodity) (err error) {
	if _, present := sdb.commoditiesByID[commodity.ID]; present == false {
		if registerIDLookup(commodity.DbEntity, sdb.commodityIDs) {
			sdb.commoditiesByID[commodity.ID] = commodity
			return nil
		}
		err = fmt.Errorf("%w: item name", ErrDuplicateEntity)
	} else {
		err = fmt.Errorf("%w: item id", ErrDuplicateEntity)
	}
	return fmt.Errorf("%s (#%d): %w", commodity.DbName, commodity.ID, err)
}

func (sdb *SystemDatabase) registerSystem(system *System) (err error) {
	if _, present := sdb.systemsByID[system.ID]; present == false {
		if registerIDLookup(system.DbEntity, sdb.systemIDs) {
			sdb.systemsByID[system.ID] = system
			return nil
		}
		err = fmt.Errorf("%w: system name", ErrDuplicateEntity)
	} else {
		err = fmt.Errorf("%w: system id", ErrDuplicateEntity)
	}

	return fmt.Errorf("%s (#%d): %w", system.DbName, system.ID, err)
}

func (sdb *SystemDatabase) GetSystemByID(id EntityID) (system *System) {
	if system, exists := sdb.systemsByID[id]; exists {
		return system
	}
	return nil
}

func (sdb *SystemDatabase) GetSystem(name string) (system *System) {
	if id, exists := sdb.systemIDs[strings.ToLower(name)]; exists {
		system = sdb.systemsByID[id]
	}
	return
}

func (sdb *SystemDatabase) registerFacility(facility *Facility) error {
	var exists bool
	system := facility.System
	if system == nil {
		return fmt.Errorf("%s (#%d): attempted to register facility without a system", facility.DbName, facility.ID)
	}
	if _, exists = sdb.facilitiesByID[facility.ID]; exists != false {
		return fmt.Errorf("%s/%s (#%d): %w: facility id", system.DbName, facility.DbName, facility.ID, ErrDuplicateEntity)
	}

	for _, existing := range system.facilities {
		if strings.EqualFold(existing.DbName, facility.DbName) {
			return fmt.Errorf("%s/%s (#%d): %w: facility name in system", system.DbName, facility.DbName, facility.ID, ErrDuplicateEntity)
		}
	}

	system.facilities = append(system.facilities, facility)
	sdb.facilitiesByID[facility.ID] = facility

	return nil
}
