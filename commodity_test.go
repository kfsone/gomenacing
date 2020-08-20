package main

import (
	"github.com/tidwall/gjson"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommodity_Name(t *testing.T) {
	commodity, err := NewCommodity(DbEntity{2, "a thing"}, 3, 10, 111)
	assert.Nil(t, err)
	assert.NotNil(t, commodity)
	assert.Equal(t, commodity.DbName, commodity.Name())
}

func TestNewCommodity(t *testing.T) {
	t.Run("NewCommodity", func(t *testing.T) {
		commodity, err := NewCommodity(DbEntity{2222, "A Thing"}, 7, 10, 123)
		assert.Nil(t, err)
		assert.NotNil(t, commodity)
		assert.Equal(t, EntityID(2222), commodity.Id)
		assert.Equal(t, "A Thing", commodity.DbName)
		assert.Equal(t, EntityID(7), commodity.CategoryId)
		assert.Equal(t, int64(10), commodity.AvgPrice)
		assert.Equal(t, FDevID(123), commodity.FDevId)
	})

	t.Run("NewCommodity rejects nil category", func(t *testing.T) {
		commodity, err := NewCommodity(DbEntity{2222, "A Thing"}, 0, 10, 123)
		if assert.Nil(t, commodity) && assert.Error(t, err) {
			assert.Equal(t, "invalid commodity category id: 0", err.Error())
		}
	})
}

func TestNewCommodityFromJsonMap(t *testing.T) {
	t.Run("Rejects non-map", func(t *testing.T) {
		json := gjson.Parse("[]")
		c, err := NewCommodityFromJsonMap(json)
		assert.Error(t, err)
		assert.Nil(t, c)
	})

	t.Run("Rejects invalid id", func(t *testing.T) {
		json := gjson.Parse(`{"id":0, "name":""}`)
		c, err := NewCommodityFromJsonMap(json)
		assert.Error(t, err)
		assert.Nil(t, c)
	})

	t.Run("Rejects invalid name", func(t *testing.T) {
		json := gjson.Parse(`{"id":1, "name":""}`)
		c, err := NewCommodityFromJsonMap(json)
		assert.Error(t, err)
		assert.Nil(t, c)
	})

	t.Run("Rejects a bad category id", func(t *testing.T) {
		json := gjson.Parse(`{"id":1, "name":"test", "category_id":0}`)
		c, err := NewCommodityFromJsonMap(json)
		assert.Error(t, err)
		assert.Nil(t, c)
	})

	t.Run("Good defaulting", func(t *testing.T) {
		json := gjson.Parse(`{"id":1, "name":"test", "category_id":3`)
		c, err := NewCommodityFromJsonMap(json)
		assert.Nil(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, Commodity{DbEntity{1, "test"}, 3, 0, 0, false}, *c)
	})

	t.Run("Accepts good data", func(t *testing.T) {
		json := gjson.Parse(`{"id":2, "name":"test2", "category_id":7, "average_price":64, "ed_id":234, "is_non_marketable":true}`)
		c, err := NewCommodityFromJsonMap(json)
		assert.Nil(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, Commodity{DbEntity{2, "test2"}, 7, 64, 234, true}, *c)
	})
}
