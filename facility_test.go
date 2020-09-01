package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//func TestNewFacility(t *testing.T) {
//	system, err := NewSystem(DbEntity{ID: 100, DbName: "system"}, Coordinate{}, false)
//	require.Nil(t, err)
//	assert.Empty(t, system.Facilities)
//
//	t.Run("Reject bad values", func(t *testing.T) {
//		entity := DbEntity{ID: 1, DbName: "first"}
//		facility, err := NewFacility(entity, nil, 0)
//		assert.Nil(t, facility)
//		if assert.Error(t, err) {
//			assert.Equal(t, "nil system", err.Error())
//		}
//
//		facility, err = NewFacility(entity, -1, 0)
//		assert.Nil(t, facility)
//		if assert.Error(t, err) {
//			assert.Equal(t, "invalid value for system id: -1", err.Error())
//		}
//
//		facility, err = NewFacility(entity, int64(-1), 0)
//		assert.Nil(t, facility)
//		if assert.Error(t, err) {
//			assert.Equal(t, "invalid value for system id: -1", err.Error())
//		}
//
//		facility, err = NewFacility(entity, EntityID(0), 0)
//		assert.Nil(t, facility)
//		if assert.Error(t, err) {
//			assert.Equal(t, "invalid value for system id: 0", err.Error())
//		}
//
//		facility, err = NewFacility(entity, struct{}{}, 0)
//		assert.Nil(t, facility)
//		if assert.Error(t, err) {
//			assert.Equal(t, "invalid parameter for system passed to NewFacility: struct {}{}", err.Error())
//		}
//	})
//
//	t.Run("Create genuine facility", func(t *testing.T) {
//		features := FeatBlackMarket | FeatSmallPad
//		facility, err := NewFacility(DbEntity{ID: 111, DbName: "First"}, system, features)
//		assert.Nil(t, err)
//		assert.NotNil(t, facility)
//		assert.Equal(t, DbEntity{ID: 111, DbName: "First"}, facility.DbEntity)
//		assert.Equal(t, system, facility.System)
//		assert.Equal(t, features, facility.Features)
//
//		assert.Empty(t, system.Facilities)
//	})
//}
//
//func TestFacility_HasFeatures(t *testing.T) {
//	facility := Facility{
//		Features: FeatMarket | FeatFleet,
//	}
//	assert.True(t, facility.HasFeatures(FeatMarket))
//	assert.True(t, facility.HasFeatures(FeatFleet))
//	assert.True(t, facility.HasFeatures(FeatMarket|FeatMarket))
//	assert.False(t, facility.HasFeatures(FeatBlackMarket|FeatMediumPad))
//	assert.False(t, facility.HasFeatures(FeatMarket|FeatMarket|FeatBlackMarket))
//
//	// If you ask if a facility has no features, it must have *no* features.
//	assert.False(t, facility.HasFeatures(FacilityFeatureMask(0)))
//
//	facility.Features =
//		FeatBlackMarket |
//			FeatCommodities |
//			FeatDocking |
//			FeatFleet |
//			FeatMarket |
//			FeatOutfitting |
//			FeatPlanetary |
//			FeatRearm |
//			FeatRefuel |
//			FeatRepair |
//			FeatShipyard |
//			0
//	assert.True(t, facility.HasFeatures(facility.Features))
//	assert.False(t, facility.HasFeatures(FeatSmallPad))
//	assert.True(t, facility.HasFeatures(FeatFleet|FeatDocking))
//
//	// Test against a facility with no features.
//	assert.True(t, Facility{}.HasFeatures(FacilityFeatureMask(0)))
//}

func TestFacility_IsTrading(t *testing.T) {
	var facility Facility
	listings := []Listing{{}, {}}
	facility = Facility{}
	assert.False(t, facility.IsTrading())
	facility = Facility{Features: FeatMarket}
	assert.True(t, facility.IsTrading())
	facility = Facility{listings: listings}
	assert.True(t, facility.IsTrading())
	facility = Facility{Features: FeatBlackMarket | FeatMediumPad | FeatPlanetary}
	assert.False(t, facility.IsTrading())
	facility = Facility{Features: FeatBlackMarket | FeatMarket}
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

func TestFacility_AddListing(t *testing.T) {
	listings := []Listing{
		{CommodityID: 30},
		{CommodityID: 7},
		{CommodityID: 19},
		{CommodityID: 20},
		{CommodityID: 2},
		{CommodityID: 3},
		{CommodityID: 4},
		{CommodityID: 128},
		{CommodityID: 129},
		{CommodityID: 127},
	}
	expectedOrder := []uint32{2, 3, 4, 7, 19, 20, 30, 127, 128, 129}

	facility := Facility{}
	assert.Nil(t, facility.listings)

	facility.AddListing(listings[0])
	require.NotNil(t, facility.listings)
	if assert.Len(t, facility.listings, 1) {
		assert.Equal(t, listings[0], facility.listings[0])
	}

	listings[0].TimestampUtc = time.Now().Add(-time.Second)
	facility.AddListing(listings[0])
	require.Len(t, facility.listings, 1)
	assert.Equal(t, listings[0], facility.listings[0])

	facility.AddListing(listings[1])
	if assert.Len(t, facility.listings, 2) {
		assert.EqualValues(t, []Listing{ listings[1], listings[0] }, facility.listings)
	}

	// Check we don't mess up with the order re-adding.
	listings[1].TimestampUtc = time.Now().Add(-time.Second)
	facility.AddListing(listings[1])
	if assert.Len(t, facility.listings, 2) {
		assert.EqualValues(t, []Listing{ listings[1], listings[0] }, facility.listings)
	}

	// Add the third value.
	facility.AddListing(listings[2])
	if assert.Len(t, facility.listings, 3) {
		assert.EqualValues(t, listings[2], facility.listings[1])
	}

	// Now insert all of the listings, and check that we build an in-order list.
	for _, listing := range listings {
		facility.AddListing(listing)
	}
	if assert.Len(t, facility.listings, len(listings)) {
		var lastID uint32
		for idx, listing := range facility.listings {
			t.Run(fmt.Sprintf("Check listing %d", idx), func(t *testing.T) {
				id := uint32(listing.CommodityID)
				if assert.Greater(t, id, lastID) {
					assert.Equal(t, expectedOrder[idx], id)
				}
				lastID = uint32(listing.CommodityID)
			})
		}
	}
}
