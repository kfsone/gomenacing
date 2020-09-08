package main

import (
	gom "github.com/kfsone/gomenacing/pkg/gomschema"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFacility(t *testing.T) {
	system := NewSystem(DbEntity{ID: 100, DbName: "system"}, Coordinate{})
	assert.Empty(t, system.facilities)

	t.Run("Reject bad values", func(t *testing.T) {
		entity := DbEntity{ID: 1, DbName: "first"}
		facility, err := NewFacility(entity, nil, 0, 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, "nil system for facility", err.Error())
		}
	})

	t.Run("Create genuine facility", func(t *testing.T) {
		facility, err := NewFacility(DbEntity{ID: 111, DbName: "First"}, system, gom.FacilityType_FTAsteroidBase, FeatLargePad)
		assert.Nil(t, err)
		assert.NotNil(t, facility)
		assert.Equal(t, DbEntity{ID: 111, DbName: "First"}, facility.DbEntity)
		assert.Equal(t, system, facility.System)
		assert.Equal(t, gom.FacilityType_FTAsteroidBase, facility.FacilityType)
		assert.Equal(t, FeatLargePad, facility.Features)
		assert.Nil(t, facility.listings)
		assert.Empty(t, system.facilities)
	})

	t.Run("Create commodity facility", func(t *testing.T) {
		facility, err := NewFacility(DbEntity{ID: 9999, DbName: "Second"}, system, gom.FacilityType_FTCoriolisStarport, FeatSmallPad|FeatCommodities)
		assert.Nil(t, err)
		assert.NotNil(t, facility)
		assert.Equal(t, DbEntity{ID: 9999, DbName: "Second"}, facility.DbEntity)
		assert.Equal(t, system, facility.System)
		assert.Equal(t, gom.FacilityType_FTCoriolisStarport, facility.FacilityType)
		assert.Equal(t, FeatSmallPad|FeatCommodities, facility.Features)
		if assert.NotNil(t, facility.listings) {
			assert.Len(t, facility.listings, 0)
		}
		assert.Empty(t, system.facilities)
	})
}

func TestFacility_GetDbId(t *testing.T) {
	facility := Facility{}
	assert.Equal(t, "000000", facility.GetDbId())
	facility.DbEntity = DbEntity{ID: 0x987ace, DbName: "Monkey Island"}
	assert.Equal(t, "987ace", facility.GetDbId())
}

func TestFacility_HasFeatures(t *testing.T) {
	// Test against a facility with no features.
	facility := Facility{}
	assert.True(t, facility.HasFeatures(FacilityFeatureMask(0)))

	facility.Features = FeatMarket | FeatFleet
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
}

func TestFacility_IsTrading(t *testing.T) {
	var facility Facility
	listing := Listing{EntityID(1), 0, 0, 0, 0, time.Now()}
	listings := map[EntityID]*Listing{ listing.CommodityID: &listing }
	facility = Facility{}
	assert.False(t, facility.IsTrading())
	facility = Facility{Features: FeatMarket}
	assert.False(t, facility.IsTrading())
	facility = Facility{Features: FeatCommodities}
	assert.True(t, facility.IsTrading())
	facility = Facility{listings: listings}
	assert.True(t, facility.IsTrading())
	facility = Facility{Features: FeatBlackMarket | FeatMarket | FeatMediumPad | FeatPlanetary}
	assert.False(t, facility.IsTrading())
	facility = Facility{Features: FeatBlackMarket | FeatCommodities}
	assert.True(t, facility.IsTrading())
	facility = Facility{Features: FeatBlackMarket, listings: listings}
	assert.True(t, facility.IsTrading())
}

func TestFacility_Name(t *testing.T) {
	system := System{DbEntity: DbEntity{DbName: "SystemName"}}
	facility := Facility{DbEntity: DbEntity{DbName: "StationName"}, System: &system}
	expectedName := "SystemName/StationName"
	assert.Equal(t, expectedName, facility.Name())
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
