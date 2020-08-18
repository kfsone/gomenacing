package main

import (
	"fmt"
	"github.com/tidwall/gjson"
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

func (f Facility) Name() string {
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

func checkSystemID(systemId int64) (entityId EntityID, err error) {
	if systemId <= 0 || systemId >= 1<<32 {
		return EntityID(0), fmt.Errorf("invalid value for system id: %d", systemId)
	}
	return EntityID(systemId), nil
}

func NewFacility(id int64, dbName string, system interface{}, features FacilityFeatureMask) (facility *Facility, err error) {
	entity, err := NewDbEntity(id, strings.ToUpper(dbName))
	if err != nil {
		return nil, err
	}

	var systemId EntityID
	var systemPtr *System

	switch typed := system.(type) {
	case nil:
		return nil, fmt.Errorf("nil system")
	case int64:
		if systemId, err = checkSystemID(typed); err != nil {
			return nil, err
		}
	case EntityID:
		if systemId, err = checkSystemID(int64(typed)); err != nil {
			return nil, err
		}
	case int:
		if systemId, err = checkSystemID(int64(typed)); err != nil {
			return nil, err
		}
	case *System:
		systemPtr = typed
		systemId = typed.Id
	default:
		return nil, fmt.Errorf("invalid parameter for system passed to NewFacility: %#v", system)
	}

	facility = &Facility{
		DbEntity: entity,
		SystemId: systemId,
		System:   systemPtr,
		Features: features,
	}

	return
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
