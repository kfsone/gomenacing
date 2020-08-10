package main

import "time"

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
	DatabaseEntity
	System         *System             // The system housing this facility.
	Features       FacilityFeatureMask // Features it has.
	LsFromStar     float64             // Distance from star.
	TypeId         int32               // Frontier facility type.
	GovernmentId   int32               // Government operating the facility.
	AllegianceId   int32               // Group to which the facility is allied.
	CommodityCount uint64              // How many commodities have been seen here.
	Updated        time.Time           // When the facility was last updated.
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
