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
	EddbSystems    string = "systems_populated.jsonl"
	EddbFacilities string = "stations.jsonl"
	//EddbCommodities string = "commodities.json"
)

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

func ImportJsonFile(filename string, fields []string, callback func(*JsonLine) error) (chan error, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// GoRoutine to consume json rows and pass them to the callback
	errorsCh := make(chan error, 4)
	go func() {
		defer func () { failOnError(file.Close()) }()
		defer close(errorsCh)
		rows := GetJsonRowsFromFile(filename, file, fields, errorsCh)
		for row := range rows {
			if err := callback(&row); err != nil {
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
