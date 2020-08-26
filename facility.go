package main

import (
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

// FacilityFeatureMask holds a bit-mask of features/services of Facilities.
type FacilityFeatureMask uint

const (
	// FeatMarket indicates a market.
	FeatMarket FacilityFeatureMask = 1 << iota
	// FeatBlackMarket indicates a black market is present.
	FeatBlackMarket
	// FeatCommodities indicates there's an actual commodity list (isn't this Market?)
	FeatCommodities
	FeatDocking    // FeatDocking indicates docking (why is this here?)
	FeatFleet      // Is a fleet.
	FeatLargePad   // Has a large pad.
	FeatMediumPad  // Has a medium pad.
	FeatOutfitting // Has an Outfitting service.
	FeatPlanetary  // Is on a planet.
	FeatRearm      // Has a Rearming service.
	FeatRefuel     // Provides refuelling.
	FeatRepair     // Provides repair service.
	FeatShipyard   // Sells ships.
	FeatSmallPad   // Has a small pad.
)

// Facility represents any orbital or planetary facility, where trade could happen.
type Facility struct {
	DbEntity
	SystemID     EntityID            `json:"system_id"` // The system housing this facility.
	System       *System             `json:"-"`
	Features     FacilityFeatureMask `json:"features"` // Features it has.
	LsFromStar   float64             `json:"ls"`       // Distance from star.
	TypeID       int32               `json:"type"`     // Frontier facility type.
	GovernmentID int32               `json:"govt"`     // Government operating the facility.
	AllegianceID int32               `json:"alleg"`    // Group to which the facility is allied.
	Updated      time.Time           `json:"updated"`  // When the facility was last updated.
	Listings     []Listing           `json:"-"`        // List of items sold
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
	return f.HasFeatures(FeatMarket) || len(f.Listings) > 0
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

func checkSystemID(systemID int64) (entityID EntityID, err error) {
	if systemID <= 0 || systemID >= 1<<32 {
		return EntityID(0), fmt.Errorf("invalid value for system id: %d", systemID)
	}
	return EntityID(systemID), nil
}

func NewFacility(entity DbEntity, system interface{}, features FacilityFeatureMask) (facility *Facility, err error) {
	var systemID EntityID
	var systemPtr *System

	switch typed := system.(type) {
	case nil:
		return nil, fmt.Errorf("nil system")
	case int64:
		if systemID, err = checkSystemID(typed); err != nil {
			return nil, err
		}
	case EntityID:
		if systemID, err = checkSystemID(int64(typed)); err != nil {
			return nil, err
		}
	case int:
		if systemID, err = checkSystemID(int64(typed)); err != nil {
			return nil, err
		}
	case *System:
		systemPtr = typed
		systemID = typed.ID
	default:
		return nil, fmt.Errorf("invalid parameter for system passed to NewFacility: %#v", system)
	}

	facility = &Facility{
		DbEntity: entity,
		SystemID: systemID,
		System:   systemPtr,
		Features: features,
	}

	return
}

func NewFacilityFromJSON(json []gjson.Result) (*Facility, error) {
	entity, err := NewDbEntityFromJSON(json)
	if err != nil {
		return nil, err
	}
	systemID := json[2].Int()
	facility, err := NewFacility(entity, systemID, 0)
	if err != nil {
		return nil, err
	}
	var featureMask = stringToFeaturePad(json[3].String())
	for i, mask := range featureMasks {
		if json[8+i].Bool() {
			featureMask |= mask
		}
	}
	facility.Features = featureMask
	facility.LsFromStar = json[4].Float()
	facility.TypeID = int32(json[5].Int())
	facility.GovernmentID = int32(json[6].Int())
	facility.AllegianceID = int32(json[7].Int())
	return facility, err
}
