package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommodity_Name(t *testing.T) {
	commodity, err := NewCommodity(2, "a thing", 3, 10, 111)
	assert.Nil(t, err)
	assert.NotNil(t, commodity)
	assert.Equal(t, commodity.DbName, commodity.Name())
}

func TestNewCommodity(t *testing.T) {
	t.Run("NewCommodity", func(t *testing.T) {
		commodity, err := NewCommodity(2222, "A Thing", 7, 10, 123)
		assert.Nil(t, err)
		assert.NotNil(t, commodity)
		assert.Equal(t, EntityID(2222), commodity.Id)
		assert.Equal(t, "A Thing", commodity.DbName)
		assert.Equal(t, EntityID(7), commodity.CategoryId)
		assert.Equal(t, int64(10), commodity.AvgPrice)
		assert.Equal(t, FDevID(123), commodity.FDevId)
	})

	t.Run("NewCommodity rejects nil category", func(t *testing.T) {
		commodity, err := NewCommodity(2222, "A Thing", 0, 10, 123)
		if assert.Nil(t, commodity) && assert.Error(t, err) {
			assert.Equal(t, "invalid commodity category id: 0", err.Error())
		}
	})
}
