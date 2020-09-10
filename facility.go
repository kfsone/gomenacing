package main

import (
	"errors"
	"fmt"
	gom "github.com/kfsone/gomenacing/pkg/gomschema"
)

// FacilityFeatureMask holds a bit-mask of features/services of Facilities.
type FacilityFeatureMask uint

const (
	FeatMarket      = FacilityFeatureMask(1 << gom.FeatureBit_Market)
	FeatBlackMarket = FacilityFeatureMask(1 << gom.FeatureBit_BlackMarket)
	FeatCommodities = FacilityFeatureMask(1 << gom.FeatureBit_Commodities)
	FeatDocking     = FacilityFeatureMask(1 << gom.FeatureBit_Docking)
	FeatFleet       = FacilityFeatureMask(1 << gom.FeatureBit_Fleet)
	FeatLargePad    = FacilityFeatureMask(1 << gom.FeatureBit_LargePad)
	FeatMediumPad   = FacilityFeatureMask(1 << gom.FeatureBit_Market)
	FeatOutfitting  = FacilityFeatureMask(1 << gom.FeatureBit_Outfitting)
	FeatPlanetary   = FacilityFeatureMask(1 << gom.FeatureBit_Planetary)
	FeatRearm       = FacilityFeatureMask(1 << gom.FeatureBit_Rearm)
	FeatRefuel      = FacilityFeatureMask(1 << gom.FeatureBit_Refuel)
	FeatRepair      = FacilityFeatureMask(1 << gom.FeatureBit_Repair)
	FeatShipyard    = FacilityFeatureMask(1 << gom.FeatureBit_Shipyard)
	FeatSmallPad    = FacilityFeatureMask(1 << gom.FeatureBit_SmallPad)
)

// Facility represents any orbital or planetary facility, where trade could happen.
type Facility struct {
	DbEntity
	System       *System
	TimestampUtc uint64              // When the facility was last updated.
	FacilityType gom.FacilityType    // Frontier facility type.
	Features     FacilityFeatureMask // Features it has.
	LsFromStar   uint32              // Distance from star.
	Government   gom.GovernmentType  // Government operating the facility.
	Allegiance   gom.AllegianceType  // Group to which the facility is allied.

	listings map[EntityID]*Listing // Table of sales/purchases
}

// NewFacility constructs a minimally populated Facility entity.
func NewFacility(dbEntity DbEntity, system *System, facilityType gom.FacilityType, features FacilityFeatureMask) (*Facility, error) {
	if system == nil {
		return nil, errors.New("nil system for facility")
	}
	facility := Facility{DbEntity: dbEntity, System: system, FacilityType: facilityType, Features: features}
	return &facility, nil
}

func (f *Facility) GetDbId() string {
	return fmt.Sprintf("%06x", f.DbEntity.ID)
}

func (f *Facility) GetTimestampUtc() uint64 {
	return uint64(f.TimestampUtc)
}

// HasFeatures returns true if the facility has a matching set of features.
// If more than one feature is specified, all of  the features must be available
// at the facility to return true.
func (f *Facility) HasFeatures(featureMask FacilityFeatureMask) bool {
	if featureMask == FacilityFeatureMask(0) {
		return f.Features == 0
	}
	return f.Features&featureMask == featureMask
}

func (f *Facility) IsTrading() bool {
	return f.HasFeatures(FeatCommodities) || len(f.listings) > 0
}

func (f *Facility) Name() string {
	return f.System.DbName + "/" + f.DbName
}

// SupportsPadSize returns true if the Facility has a pad of known size >= size.
func (f *Facility) SupportsPadSize(size FacilityFeatureMask) bool {
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
