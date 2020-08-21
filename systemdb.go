package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

var EddbPath = flag.StringP("eddbdir", "e", "", "Path to EDDB json files to import.")

const (
	EddbSystems     string = "systems_populated.jsonl"
	EddbFacilities  string = "stations.jsonl"
	EddbCommodities string = "commodities.json"
)

type SystemDatabase struct {
	db *Database
	// Index of Systems by their database ids.
	systemsById map[EntityID]*System
	// Look-up a system's EntityID by it's name.
	systemIds map[string]EntityID
	// Index of Facilities by their database ids.
	facilitiesById map[EntityID]*Facility
	// Index of Commodities by their database ids.
	commoditiesById map[EntityID]*Commodity
	// Look-up a commodity's EntityID by it's name.
	commodityIds map[string]EntityID
}

var systemFields = []string{
	"id", "name", "x", "y", "z", "needs_permit",
}

var featureMasks = []FacilityFeatureMask{
	FeatBlackMarket, FeatCommodities, FeatDocking, FeatMarket, FeatOutfitting, FeatRearm, FeatRefuel, FeatRepair, FeatShipyard, FeatPlanetary,
}

var facilityFields = []string{
	"id",
	"name",
	"system_id",
	"max_landing_pad_size",
	"distance_to_star",
	"type_id",
	"government_id",
	"allegiance_id",
	"has_blackmarket",
	"has_commodities",
	"has_docking",
	"has_market",
	"has_outfitting",
	"has_rearm",
	"has_refuel",
	"has_repair",
	"has_shipyard",
	"is_planetary",
}

func NewSystemDatabase(db *Database) *SystemDatabase {
	return &SystemDatabase{
		db:              db,
		systemsById:     make(map[EntityID]*System),
		systemIds:       make(map[string]EntityID),
		facilitiesById:  make(map[EntityID]*Facility),
		commoditiesById: make(map[EntityID]*Commodity),
		commodityIds:    make(map[string]EntityID),
	}
}

func ImportJsonlFile(filename string, fields []string, callback func(*JsonEntry) error) (chan error, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// GoRoutine to consume json rows and pass them to the callback
	errorsCh := make(chan error, 4)
	go func() {
		defer func () { failOnError(file.Close()) }()
		defer close(errorsCh)
		rows := GetJsonItemsFromFile(filename, file, fields, errorsCh)
		for row := range rows {
			if err := callback(&row); err != nil {
				errorsCh <- fmt.Errorf("%s:%d: %w", filename, row.LineNo, err)
			}
		}
	}()

	return errorsCh, nil
}

// Consumes an error channel, counts and reports any errors that don't get filtered
// by the environment and either returns nil if there were no errors, or an error
// describing how many errors there were.
func countErrors(errorCh <-chan error) error {
	var errorCount int
	for err := range errorCh {
		if err = FilterError(err); err != nil {
			errorCount += 1
			log.Print(err.Error())
		}
	}
	if errorCount > 0 {
		return fmt.Errorf("failed because of %d error(s)", errorCount)
	}
	return nil
}

func registerIdLookup(entity DbEntity, ids map[string]EntityID) bool {
	name := strings.ToLower(entity.DbName)
	if _, present := ids[name]; present != false {
		return false
	}
	ids[name] = entity.Id
	return true
}

func (sdb *SystemDatabase) registerCommodity(commodity *Commodity) (err error) {
	if _, present := sdb.commoditiesById[commodity.Id]; present == false {
		if registerIdLookup(commodity.DbEntity, sdb.commodityIds) {
			sdb.commoditiesById[commodity.Id] = commodity
			return nil
		}
		err = fmt.Errorf("%w: item name", ErrDuplicateEntity)
	} else {
		err = fmt.Errorf("%w: item id", ErrDuplicateEntity)
	}
	return fmt.Errorf("%s (#%d): %w", commodity.DbName, commodity.Id, err)
}

func (sdb *SystemDatabase) registerSystem(system *System) (err error) {
	if _, present := sdb.systemsById[system.Id]; present == false {
		if registerIdLookup(system.DbEntity, sdb.systemIds) {
			sdb.systemsById[system.Id] = system
			return nil
		}
		err = fmt.Errorf("%w: system name", ErrDuplicateEntity)
	} else {
		err = fmt.Errorf("%w: system id", ErrDuplicateEntity)
	}

	return fmt.Errorf("%s (#%d): %w", system.DbName, system.Id, err)
}

func (sdb *SystemDatabase) GetSystem(name string) (system *System) {
	if id, exists := sdb.systemIds[strings.ToLower(name)]; exists {
		system = sdb.systemsById[id]
	}
	return
}

func (sdb *SystemDatabase) registerFacility(facility *Facility) error {
	var exists bool
	system, systemId := facility.System, facility.SystemId
	if system == nil {
		if systemId == 0 {
			return fmt.Errorf("%s (#%d): attempted to register facility without a system id", facility.DbName, facility.Id)
		}
		system, exists = sdb.systemsById[facility.SystemId]
		if exists == false {
			return fmt.Errorf("%s (#%d): system id: %w: %d", facility.DbName, facility.Id, ErrUnknownEntity, systemId)
		}
	} else {
		systemId = system.Id
	}
	if _, exists = sdb.facilitiesById[facility.Id]; exists != false {
		return fmt.Errorf("%s/%s (#%d): %w: facility id", system.DbName, facility.DbName, facility.Id, ErrDuplicateEntity)
	}

	for _, existing := range system.Facilities {
		if strings.EqualFold(existing.DbName, facility.DbName) {
			return fmt.Errorf("%s/%s (#%d): %w: facility name in system", system.DbName, facility.DbName, facility.Id, ErrDuplicateEntity)
		}
	}

	facility.System, facility.SystemId = system, systemId
	system.Facilities = append(system.Facilities, facility)
	sdb.facilitiesById[facility.Id] = facility

	return nil
}
