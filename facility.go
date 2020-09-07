package main

import (
	"errors"
	"fmt"
	gom "github.com/kfsone/gomenacing/pkg/gomschema"
	"sort"
	"time"
)

// FacilityFeatureMask holds a bit-mask of features/services of Facilities.
type FacilityFeatureMask uint

const (
	FeatMarket = FacilityFeatureMask(gom.FeatureMasks_FeatMarket)
	FeatBlackMarket = FacilityFeatureMask(gom.FeatureMasks_FeatBlackMarket)
	FeatCommodities = FacilityFeatureMask(gom.FeatureMasks_FeatCommodities)
	FeatDocking = FacilityFeatureMask(gom.FeatureMasks_FeatDocking)
	FeatFleet = FacilityFeatureMask(gom.FeatureMasks_FeatFleet)
	FeatLargePad = FacilityFeatureMask(gom.FeatureMasks_FeatLargePad)
	FeatMediumPad = FacilityFeatureMask(gom.FeatureMasks_FeatMarket)
	FeatOutfitting = FacilityFeatureMask(gom.FeatureMasks_FeatOutfitting)
	FeatPlanetary = FacilityFeatureMask(gom.FeatureMasks_FeatPlanetary)
	FeatRearm = FacilityFeatureMask(gom.FeatureMasks_FeatRearm)
	FeatRefuel = FacilityFeatureMask(gom.FeatureMasks_FeatRefuel)
	FeatRepair = FacilityFeatureMask(gom.FeatureMasks_FeatRepair)
	FeatShipyard = FacilityFeatureMask(gom.FeatureMasks_FeatShipyard)
	FeatSmallPad = FacilityFeatureMask(gom.FeatureMasks_FeatSmallPad)
)

// Facility represents any orbital or planetary facility, where trade could happen.
type Facility struct {
	DbEntity
	System       *System
	TimestampUtc time.Time           // When the facility was last updated.
	FacilityType gom.FacilityType    // Frontier facility type.
	Features     FacilityFeatureMask // Features it has.
	LsFromStar   uint32              // Distance from star.
	Government   gom.GovernmentType  // Government operating the facility.
	Allegiance   gom.AllegianceType  // Group to which the facility is allied.

	listings []Listing // List of items sold
}

// NewFacility constructs a minimally populated Facility entity. It does not allocate
// listings unless the FeatureMask indicates Commodities are available here.
func NewFacility(dbEntity DbEntity, system *System, facilityType gom.FacilityType, features FacilityFeatureMask) (*Facility, error) {
	if system == nil {
		return nil, errors.New("nil system for facility")
	}
	facility := Facility{DbEntity: dbEntity, System: system, FacilityType: facilityType, Features: features}
	if facility.HasFeatures(FeatCommodities) {
		facility.listings = make([]Listing, 0, 32)
	}
	return &facility, nil
}

func (f *Facility) GetDbId() string {
	return fmt.Sprintf("%06x", f.DbEntity.ID)
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

func (f *Facility) AddListing(listing Listing) {
	commodity := listing.CommodityID
	if f.listings == nil {
		f.listings = make([]Listing, 0, 32)
	}
	var insertIdx = 0
	if len(f.listings) > 0 {
		insertIdx = sort.Search(len(f.listings), func(i int) bool { return f.listings[i].CommodityID >= commodity })
	}
	if insertIdx >= len(f.listings) {
		f.listings = append(f.listings, listing)
		return
	}
	if f.listings[insertIdx].CommodityID != commodity {
		f.listings = append(f.listings, listing)
		copy(f.listings[insertIdx+1:], f.listings[insertIdx:])
	}
	f.listings[insertIdx] = listing
}
