package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFacility_HasFeatures(t *testing.T) {
	facility := Facility{
		Features: FeatMarket | FeatFleet,
	}
	assert.True(t, facility.HasFeatures(FeatMarket))
	assert.True(t, facility.HasFeatures(FeatFleet))
	assert.True(t, facility.HasFeatures(FeatMarket|FeatMarket))
	assert.False(t, facility.HasFeatures(FeatBlackMarket|FeatMediumPad))
	assert.False(t, facility.HasFeatures(FeatMarket|FeatMarket|FeatBlackMarket))

	// If you ask if a facility has no features, it must have *no* features.
	assert.False(t, facility.HasFeatures(FacilityFeatureMask(0)))

	facility.Features =
		FeatBlackMarket |
			FeatCommodities |
			FeatDocking |
			FeatFleet |
			FeatMarket |
			FeatOutfitting |
			FeatPlanetary |
			FeatRearm |
			FeatRefuel |
			FeatRepair |
			FeatShipyard |
			0
	assert.True(t, facility.HasFeatures(facility.Features))
	assert.False(t, facility.HasFeatures(FeatSmallPad))
	assert.True(t, facility.HasFeatures(FeatFleet|FeatDocking))

	// Test against a facility with no features.
	assert.True(t, Facility{}.HasFeatures(FacilityFeatureMask(0)))
}

func TestFacility_IsTrading(t *testing.T) {
	assert.False(t, Facility{}.IsTrading())
	assert.True(t, Facility{Features: FeatMarket}.IsTrading())
	assert.True(t, Facility{CommodityCount: 1}.IsTrading())
	assert.False(t, Facility{Features: FeatBlackMarket | FeatMediumPad | FeatPlanetary}.IsTrading())
	assert.True(t, Facility{Features: FeatBlackMarket | FeatMarket}.IsTrading())
	assert.True(t, Facility{Features: FeatBlackMarket, CommodityCount: 2}.IsTrading())
}

func TestFacility_Name(t *testing.T) {
	system := System{DatabaseEntity: DatabaseEntity{DbName: "SystemName"}}
	facility := Facility{DatabaseEntity: DatabaseEntity{DbName: "StationName"}, System: &system}
	expectedName := "SystemName/StationName"
	assert.Equal(t, expectedName, facility.Name(-1))
	assert.Equal(t, expectedName, facility.Name(0))
	assert.Equal(t, expectedName, facility.Name(1))
	assert.Equal(t, expectedName, facility.Name(3))
}

func TestFacility_SupportsPadSize(t *testing.T) {
	facility := Facility{Features: FacilityFeatureMask(0)}
	assert.False(t, facility.SupportsPadSize(FeatSmallPad))
	assert.False(t, facility.SupportsPadSize(FeatMediumPad))
	assert.False(t, facility.SupportsPadSize(FeatLargePad))
	assert.False(t, facility.SupportsPadSize(FacilityFeatureMask(0)))
	assert.False(t, facility.SupportsPadSize(FeatShipyard))

	facility.Features = FeatSmallPad
	assert.True(t, facility.SupportsPadSize(FeatSmallPad))
	assert.False(t, facility.SupportsPadSize(FeatMediumPad))
	assert.False(t, facility.SupportsPadSize(FeatLargePad))
	assert.False(t, facility.SupportsPadSize(FacilityFeatureMask(0)))
	assert.False(t, facility.SupportsPadSize(FeatOutfitting))

	facility.Features = FeatMediumPad
	assert.True(t, facility.SupportsPadSize(FeatSmallPad))
	assert.True(t, facility.SupportsPadSize(FeatMediumPad))
	assert.False(t, facility.SupportsPadSize(FeatLargePad))
	assert.False(t, facility.SupportsPadSize(FacilityFeatureMask(0)))
	assert.False(t, facility.SupportsPadSize(FeatRearm))

	facility.Features = FeatLargePad
	assert.True(t, facility.SupportsPadSize(FeatSmallPad))
	assert.True(t, facility.SupportsPadSize(FeatMediumPad))
	assert.True(t, facility.SupportsPadSize(FeatLargePad))
	assert.False(t, facility.SupportsPadSize(FacilityFeatureMask(0)))
	assert.False(t, facility.SupportsPadSize(FeatRefuel))
}
