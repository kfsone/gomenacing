package main

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"strings"
	"time"
)

// FacilityFeatureMask holds a bit-mask of features/services of Facilities.
type FacilityFeatureMask uint

const (
	FeatMarket      FacilityFeatureMask = 1 << iota // Has a market.
	FeatBlackMarket                                 // Has a black market.
	FeatCommodities                                 // Has a commodity list (isn't this Market?)
	FeatDocking                                     // Provides Docking (why is this here?)
	FeatFleet                                       // Is a fleet.
	FeatLargePad                                    // Has a large pad.
	FeatMediumPad                                   // Has a medium pad.
	FeatOutfitting                                  // Has an Outfitting service.
	FeatPlanetary                                   // Is on a planet.
	FeatRearm                                       // Has a Rearming service.
	FeatRefuel                                      // Provides refuelling.
	FeatRepair                                      // Provides repair service.
	FeatShipyard                                    // Sells ships.
	FeatSmallPad                                    // Has a small pad.
)

// Facility represents any orbital or planetary facility, where trade could happen.
type Facility struct {
	DbEntity
	SystemId       EntityID            `json:"system_id"` // The system housing this facility.
	System         *System             `json:"-"`
	Features       FacilityFeatureMask `json:"features"` // Features it has.
	LsFromStar     float64             `json:"ls"`       // Distance from star.
	TypeId         int32               `json:"type"`     // Frontier facility type.
	GovernmentId   int32               `json:"govt"`     // Government operating the facility.
	AllegianceId   int32               `json:"alleg"`    // Group to which the facility is allied.
	CommodityCount uint64              `json:"-"`        // How many commodities have been seen here.
	Updated        time.Time           `json:"updated"`  // When the facility was last updated.
}

func (f Facility) Name(_ int) string {
	return f.System.DbName + "/" + f.DbName
}

// HasFeatures returns true if the facility has a matching set of features.
// If more than one feature is specified, all of  the features must be available
// at the facility to return true.
func (f Facility) HasFeatures(featureMask FacilityFeatureMask) bool {
	if featureMask == FacilityFeatureMask(0) {
		return f.Features == 0
	}
	return f.Features&featureMask == featureMask
}

func (f Facility) IsTrading() bool {
	return f.HasFeatures(FeatMarket) || f.CommodityCount > 0
}

// SupportsPadSize returns true if the Facility has a pad of known size >= size.
func (f Facility) SupportsPadSize(size FacilityFeatureMask) bool {
	switch size {
	case FeatLargePad:
		return f.Features&FeatLargePad != 0
	case FeatMediumPad:
		return f.Features&(FeatMediumPad|FeatLargePad) != 0
	case FeatSmallPad:
		return f.Features&(FeatSmallPad|FeatMediumPad|FeatLargePad) != 0
	default:
		return false
	}
}

func NewFacility(id int64, dbName string, system interface{}, features FacilityFeatureMask) (*Facility, error) {
	if id <= 0 || id >= 1<<32 {
		return nil, errors.New(fmt.Sprintf("invalid facility id: %d", id))
	}
	dbName = strings.TrimSpace(dbName)
	if len(dbName) == 0 {
		return nil, errors.New("invalid (empty) facility name")
	}

	var systemId EntityID
	var systemPtr *System

	switch typed := system.(type) {
	case int64:
		systemId = EntityID(typed)
	case *System:
		systemPtr = typed
		systemId = typed.Id
	default:
		log.Fatalf("invalid parameter for system passed to NewFacility: %#v", system)
	}

	facility := &Facility{
		DbEntity: DbEntity{
			Id:     EntityID(id),
			DbName: strings.ToUpper(dbName),
		},
		SystemId: systemId,
		System:   systemPtr,
		Features: features,
	}

	return facility, nil
}

func NewFacilityFromJson(json []gjson.Result) (*Facility, error) {
	facilityId, facilityName, systemId := json[0].Int(), json[1].String(), json[2].Int()
	var featureMask = stringToFeaturePad(json[3].String())
	for i, mask := range featureMasks {
		if json[8+i].Bool() {
			featureMask |= mask
		}
	}
	facility, err := NewFacility(facilityId, facilityName, systemId, featureMask)
	if err == nil {
		facility.LsFromStar = json[4].Float()
		facility.TypeId = int32(json[5].Int())
		facility.GovernmentId = int32(json[6].Int())
		facility.AllegianceId = int32(json[7].Int())
	}

	return facility, nil
}
