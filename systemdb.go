package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const SystemsJson = "systems_populated.jsonl"
const filename = "c:/users/oliver/data/eddb/systems_populated.jsonl"

type SystemDatabase struct {
	systemsById    map[EntityID]*System
	systemIds      map[string]EntityID
	facilitiesById map[EntityID]*Facility
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

func NewSystemDatabase() *SystemDatabase {
	return &SystemDatabase{make(map[EntityID]*System), make(map[string]EntityID), make(map[EntityID]*Facility)}
}

func ImportJsonFile(filename string, fields []string, callback func(JsonLine) error) (chan error, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// GoRoutine to consume json rows and pass them to the callback
	errorsCh := make(chan error, 4)
	go func() {
		defer file.Close()

		defer close(errorsCh)
		rows := GetJsonRowsFromFile(filename, file, fields, errorsCh)
		for row := range rows {
			if err := callback(row); err != nil {
				errorsCh <- err
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

func (sdb *SystemDatabase) importSystems() error {
	///TODO: Get name from env
	var system *System
	errorsCh, err := ImportJsonFile(filename, systemFields, func(json JsonLine) (err error) {
		if system, err = NewSystemFromJson(json.Results); err != nil {
			return fmt.Errorf("%s:%d: %w", filename, json.LineNo, err)
		}
		if err = sdb.registerSystem(system); err != nil {
			return fmt.Errorf("%s:%d: %w", filename, json.LineNo, err)
		}
		return nil
	})
	err = countErrors(errorsCh)
	if err == nil {
	}
	if err != nil {
		return err
	}

	fmt.Printf("Loaded %d systems.\n", len(sdb.systemIds))

	return nil
}

func (sdb *SystemDatabase) importFacilities() error {
	///TODO: Get name from env
	const filename = "c:/users/oliver/data/eddb/stations.jsonl"
	var facility *Facility
	errorsCh, err := ImportJsonFile(filename, facilityFields, func(json JsonLine) (err error) {
		if facility, err = NewFacilityFromJson(json.Results, sdb); err != nil {
			return fmt.Errorf("%s:%d: %w", filename, json.LineNo, err)
		}
		if err = sdb.registerFacility(facility); err != nil {
			return fmt.Errorf("%s:%d: %w", filename, json.LineNo, err)
		}
		return nil
	})
	if err == nil {
		err = countErrors(errorsCh)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Loaded %d facilities.\n", len(sdb.facilitiesById))

	return nil
}

func (sdb *SystemDatabase) registerSystem(system *System) error {
	if _, present := sdb.systemsById[system.Id]; present != false {
		return fmt.Errorf("%s (#%d): %w: system id", system.DbName, system.Id, ErrDuplicateEntity)
	}
	if _, present := sdb.systemIds[system.DbName]; present != false {
		return fmt.Errorf("%s (#%d): %w: system name", system.DbName, system.Id, ErrDuplicateEntity)
	}

	sdb.systemsById[system.Id] = system
	sdb.systemIds[system.DbName] = system.Id

	return nil
}

func (sdb *SystemDatabase) registerFacility(facility *Facility) error {
	if facility.System == nil {
		return errors.New("attempted to register facility with nil system")
	}

	if _, exists := sdb.facilitiesById[facility.Id]; exists != false {
		return fmt.Errorf("%s (#%d): %w: facility id", facility.Name(2), facility.Id, ErrDuplicateEntity)
	}

	for _, existing := range facility.System.Facilities {
		if strings.EqualFold(existing.DbName, facility.DbName) {
			return fmt.Errorf("%s (#%d): %w: facility name in system", facility.Name(2), facility.Id, ErrDuplicateEntity)
		}
	}

	facility.System.Facilities = append(facility.System.Facilities, facility)
	sdb.facilitiesById[facility.Id] = facility

	return nil
}
