package main

import (
	gom "github.com/kfsone/gomenacing/pkg/gomschema"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommodity_Name(t *testing.T) {
	commodity, err := NewCommodity(DbEntity{ID: 2, DbName: "a thing"}, gom.Commodity_CatFoods, false, false, 10)
	if assert.Nil(t, err) {
		assert.NotNil(t, commodity)
		assert.Equal(t, "a thing", commodity.Name())
	}

	commodity, err = NewCommodity(DbEntity{ID: 2, DbName: "another thing"}, gom.Commodity_CatConsumerItems, true, true, 100)
	if assert.Nil(t, err) {
		assert.NotNil(t, commodity)
		assert.Equal(t, "another thing", commodity.Name())
	}
}

func TestNewCommodity(t *testing.T) {
	commodity, err := NewCommodity(DbEntity{ID: 2222, DbName: "A Thing"}, gom.Commodity_CatChemicals, true, false, 472)
	if assert.Nil(t, err) {
		assert.NotNil(t, commodity)
		assert.Equal(t, EntityID(2222), commodity.ID)
		assert.Equal(t, "A Thing", commodity.DbName)
		assert.True(t, commodity.IsRare)
		assert.False(t, commodity.IsNonMarketable)
		assert.Equal(t, gom.Commodity_CatChemicals, commodity.CategoryID)
		assert.Equal(t, uint32(472), commodity.AverageCr)
	}
}
