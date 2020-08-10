package main

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"os"
	"strings"
)

type SystemDatabase struct {
	systemsById    map[EntityID]*System
	systemIds      map[string]EntityID
	facilitiesById map[EntityID]*Facility
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

func (sdb *SystemDatabase) importSystems(env *Env) error {
	///TODO: Get name from env
	const filename = "c:/users/oliver/data/eddb/systems_populated.jsonl"
	if file, err := os.Open(filename); err != nil {
		return err
	} else {
		defer file.Close()
		err = IterateLinesInFile(filename, file, func(json string) error {
			if _, err := sdb.addSystemFromJson(json); err != nil {
				return env.FilterError(err)
			}
			return nil
		})
		fmt.Printf("Loaded %d Systems.\n", len(sdb.systemsById))

		return err
	}
}

func (sdb *SystemDatabase) importFacilities(env *Env) error {
	///TODO: Get name from env
	const filename = "c:/users/oliver/data/eddb/stations.jsonl"
	if file, err := os.Open(filename); err != nil {
		return err
	} else {
		defer file.Close()
		err = IterateLinesInFile(filename, file, func(json string) error {
			if _, err := sdb.addFacilityFromJson(json); err != nil {
				return env.FilterError(err)
			}
			return nil
		})
		fmt.Printf("Loaded %d Facilities.\n", len(sdb.facilitiesById))

		return err
	}
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
		if existing.Id == facility.Id {
			return fmt.Errorf("%s (#%d): %w: facility id in system", facility.Name(2), facility.Id, ErrDuplicateEntity)
		}
		if strings.EqualFold(existing.DbName, facility.DbName) {
			return fmt.Errorf("%s (#%d): %w: facility name in system", facility.Name(2), facility.Id, ErrDuplicateEntity)
		}
	}

	facility.System.Facilities = append(facility.System.Facilities, facility)
	sdb.facilitiesById[facility.Id] = facility

	return nil
}

func (sdb *SystemDatabase) addSystemFromJson(json string) (*System, error) {
	if !gjson.Valid(json) {
		return nil, errors.New("invalid json: " + json)
	}
	results := gjson.GetMany(json, "id", "name", "x", "y", "z", "needs_permit")
	if len(results) != 6 {
		return nil, errors.New("malformed system entry: " + json)
	}
	position := Coordinate{results[2].Float(), results[3].Float(), results[4].Float()}
	system, err := NewSystem(results[0].Int(), results[1].String(), position, results[5].Bool())
	if err != nil {
		return nil, err
	}
	if err = sdb.registerSystem(system); err != nil {
		return nil, err
	}

	return system, nil
}

func (sdb *SystemDatabase) addFacilityFromJson(json string) (*Facility, error) {
	if !gjson.Valid(json) {
		return nil, errors.New("invalid json: " + json)
	}
	results := gjson.GetMany(json, facilityFields...)
	if len(results) != len(facilityFields) {
		return nil, errors.New("malformed facility entry: " + json)
	}
	facilityId, facilityName, systemId := results[0].Int(), results[1].String(), EntityID(results[2].Int())
	system, ok := sdb.systemsById[systemId]
	if !ok {
		return nil, fmt.Errorf("%s (#%d): %w: system id #%d", facilityName, facilityId, ErrUnknownEntity, systemId)
	}
	var featureMask = stringToFeaturePad(results[3].String())
	for i, mask := range featureMasks {
		if results[8+i].Bool() {
			featureMask |= mask
		}
	}
	facility, err := system.NewFacility(facilityId, facilityName, featureMask)
	if err != nil {
		return nil, err
	}
	facility.LsFromStar = results[4].Float()
	facility.TypeId = int32(results[5].Int())
	facility.GovernmentId = int32(results[6].Int())
	facility.AllegianceId = int32(results[7].Int())

	if err = sdb.registerFacility(facility); err != nil {
		return nil, err
	}

	return facility, nil
}
