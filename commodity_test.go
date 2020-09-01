package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommodity_Name(t *testing.T) {
	commodity, err := NewCommodity(DbEntity{ID: 2, DbName: "a thing"}, 3, 10)
	assert.Nil(t, err)
	assert.NotNil(t, commodity)
	assert.Equal(t, commodity.DbName, commodity.Name())
}

func TestNewCommodity(t *testing.T) {
	t.Run("NewCommodity", func(t *testing.T) {
		commodity, err := NewCommodity(DbEntity{ID: 2222, DbName: "A Thing"}, 7, 10)
		assert.Nil(t, err)
		assert.NotNil(t, commodity)
		assert.Equal(t, EntityID(2222), commodity.ID)
		assert.Equal(t, "A Thing", commodity.DbName)
		assert.Equal(t, EntityID(7), commodity.CategoryID)
		assert.Equal(t, int64(10), commodity.AverageCr)
	})

	t.Run("NewCommodity rejects nil category", func(t *testing.T) {
		commodity, err := NewCommodity(DbEntity{ID: 2222, DbName: "A Thing"}, 0, 10)
		if assert.Nil(t, commodity) && assert.Error(t, err) {
			assert.Equal(t, "invalid commodity category id: 0", err.Error())
		}
	})
}
