package main

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"strings"
)

type SystemDatabase struct {
	systemsById    map[EntityID]*System
	systemIds      map[string]EntityID
	facilitiesById map[EntityID]*Facility
}

var systemFields = []string{
	"id", "name", "x", "y", "z", "needs_permit",
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

var featureMasks = []FacilityFeatureMask{
	FeatBlackMarket, FeatCommodities, FeatDocking, FeatMarket, FeatOutfitting, FeatRearm, FeatRefuel, FeatRepair, FeatShipyard, FeatPlanetary,
}

func NewSystemDatabase() SystemDatabase {
	return SystemDatabase{make(map[EntityID]*System), make(map[string]EntityID), make(map[EntityID]*Facility)}
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
func countErrors(env *Env, filename string, errorCh <-chan error) error {
	var errorCount int
	for err := range errorCh {
		if err = env.FilterError(err); err != nil {
			errorCount += 1
			log.Print(err.Error())
		}
	}
	if errorCount > 0 {
		return fmt.Errorf("failed because of %d error(s)", errorCount)
	}
	return nil
}

func (sdb *SystemDatabase) importSystems(env *Env) error {
	///TODO: Get name from env
	const filename = "c:/users/oliver/data/eddb/systems_populated.jsonl"
	var system *System
	errorsCh, err := ImportJsonFile(filename, systemFields, func(json JsonLine) (err error) {
		if system, err = sdb.makeSystemFromJson(json.Results); err != nil {
			return fmt.Errorf("%s:%d: %w", filename, json.LineNo, err)
		}
		if err = sdb.registerSystem(system); err != nil {
			return fmt.Errorf("%s:%d: %w", filename, json.LineNo, err)
		}
		return nil
	})
	if err == nil {
		err = countErrors(env, filename, errorsCh)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Loaded %d systems.\n", len(sdb.systemIds))

	return nil
}

func (sdb *SystemDatabase) importFacilities(env *Env) error {
	///TODO: Get name from env
	const filename = "c:/users/oliver/data/eddb/stations.jsonl"
	var facility *Facility
	errorsCh, err := ImportJsonFile(filename, facilityFields, func(json JsonLine) (err error) {
		if facility, err = sdb.makeFacilityFromJson(json.Results); err != nil {
			return fmt.Errorf("%s:%d: %w", filename, json.LineNo, err)
		}
		if err = sdb.registerFacility(facility); err != nil {
			return fmt.Errorf("%s:%d: %w", filename, json.LineNo, err)
		}
		return nil
	})
	if err == nil {
		err = countErrors(env, filename, errorsCh)
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

func (sdb *SystemDatabase) makeSystemFromJson(json []gjson.Result) (*System, error) {
	position := Coordinate{json[2].Float(), json[3].Float(), json[4].Float()}
	return NewSystem(json[0].Int(), json[1].String(), position, json[5].Bool())
}

func (sdb *SystemDatabase) makeFacilityFromJson(json []gjson.Result) (*Facility, error) {
	facilityId, facilityName, systemId := json[0].Int(), json[1].String(), EntityID(json[2].Int())
	system, ok := sdb.systemsById[systemId]
	if !ok {
		return nil, fmt.Errorf("%s (#%d): %w: system id #%d", facilityName, facilityId, ErrUnknownEntity, systemId)
	}
	var featureMask = stringToFeaturePad(json[3].String())
	for i, mask := range featureMasks {
		if json[8+i].Bool() {
			featureMask |= mask
		}
	}
	facility, err := system.NewFacility(facilityId, facilityName, featureMask)
	if err == nil {
		facility.LsFromStar = json[4].Float()
		facility.TypeId = int32(json[5].Int())
		facility.GovernmentId = int32(json[6].Int())
		facility.AllegianceId = int32(json[7].Int())
	}

	return facility, nil
}
