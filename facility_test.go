package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFacility(t *testing.T) {
	system, err := NewSystem(100, "system", Coordinate{}, false)
	require.Nil(t, err)
	assert.Empty(t, system.Facilities)

	t.Run("Reject bad values", func(t *testing.T) {
		facility, err := NewFacility(0, "", nil, 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, "invalid id: 0", err.Error())
		}

		facility, err = NewFacility(1<<32, "", nil, 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, fmt.Errorf("invalid id: %d", 1<<32), err)
		}

		facility, err = NewFacility(1<<32-1, "", nil, 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, "invalid/empty name: \"\"", err.Error())
		}

		facility, err = NewFacility(1, " \t ", nil, 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, "invalid/empty name: \" \t \"", err.Error())
		}

		facility, err = NewFacility(1, "first", nil, 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, "nil system", err.Error())
		}

		facility, err = NewFacility(1, "first", -1, 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, "invalid value for system id: -1", err.Error())
		}

		facility, err = NewFacility(1, "first", int64(-1), 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, "invalid value for system id: -1", err.Error())
		}

		facility, err = NewFacility(1, "first", EntityID(0), 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, "invalid value for system id: 0", err.Error())
		}

		facility, err = NewFacility(1, "first", struct{}{}, 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, "invalid parameter for system passed to NewFacility: struct {}{}", err.Error())
		}
	})

	t.Run("Create genuine facility", func(t *testing.T) {
		features := FeatBlackMarket | FeatSmallPad
		facility, err := NewFacility(111, "first", system, features)
		assert.Nil(t, err)
		assert.NotNil(t, facility)
		assert.Equal(t, EntityID(111), facility.Id)
		assert.Equal(t, "FIRST", facility.DbName)
		assert.Equal(t, system, facility.System)
		assert.Equal(t, features, facility.Features)

		assert.Empty(t, system.Facilities)
	})
}

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

func TestFacility_NewFacilityFromJson(t *testing.T) {
	json := `{
		"id": 3, "name": "Luna","system_id": "1",
		"max_landing_pad_size": "M",
		"type_id": 8,
		"government_id": 13,
		"allegiance_id": 27,
		"distance_to_star": 1.204,
		"has_blackmarket": true,
		"has_commodities": false,
		"has_docking": true,
		"has_market": false,
		"has_outfitting": true,
		"has_rearm": false,
		"has_refuel": true,
		"has_repair": false,
		"has_shipyard": true,
		"is_planetary": false
	}`
	results := gjson.GetMany(json, facilityFields...)
	require.Len(t, facilityFields, len(results))

	facility, err := NewFacilityFromJson(results)
	require.Nil(t, err)
	require.NotNil(t, facility)
	assert.Equal(t, EntityID(3), facility.Id)
	assert.Equal(t, "LUNA", facility.DbName)
	assert.EqualValues(t, 1, facility.SystemId)
	assert.Nil(t, facility.System)
	assert.EqualValues(t, 8, facility.TypeId)
	assert.EqualValues(t, 13, facility.GovernmentId)
	assert.EqualValues(t, 27, facility.AllegianceId)
	assert.Equal(t, 1.204, facility.LsFromStar)
	assert.Equal(t, FeatMediumPad|FeatBlackMarket|FeatDocking|FeatOutfitting|FeatRefuel|FeatShipyard, facility.Features)
}
