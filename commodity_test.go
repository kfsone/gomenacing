package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommodity_Name(t *testing.T) {
	category := NewCategory(1001, "Things")
	commodity, err := NewCommodity(2, "a thing", &category, "a genuine thing", 10, 111)
	assert.Nil(t, err)
	assert.NotNil(t, commodity)
	assert.Equal(t, commodity.DbName, commodity.Name(-1))
	assert.Equal(t, commodity.DbName, commodity.Name(0))
	assert.Equal(t, commodity.FullName, commodity.Name(1))
	assert.Equal(t, commodity.FullName, commodity.Name(2))
	assert.Equal(t, commodity.FullName, commodity.Name(999))
}

func TestNewCommodity(t *testing.T) {
	t.Run("NewCommodity", func(t *testing.T) {
		category := NewCategory(1010, "Things")
		commodity, err := NewCommodity(2222, "A Thing", &category, "A Genuine Thing", 10, 123)
		assert.Nil(t, err)
		assert.NotNil(t, commodity)
		assert.Equal(t, EntityID(2222), commodity.Id)
		assert.Equal(t, "A Thing", commodity.DbName)
		assert.Equal(t, &category, commodity.Category)
		assert.Equal(t, "A Genuine Thing", commodity.FullName)
		assert.Equal(t, 10, commodity.AvgPrice)
		assert.Equal(t, FDevID(123), commodity.FDevId)
	})

	t.Run("NewCommodity rejects nil category", func(t *testing.T) {
		commodity, err := NewCommodity(2222, "A Thing", nil, "A Genuine Thing", 10, 123)
		assert.Nil(t, commodity)
		assert.Equal(t, MissingCategory, err)
	})
}
