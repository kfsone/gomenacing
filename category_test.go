package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategory_Name(t *testing.T) {
	assert.Equal(t, "Hello World", Category{DbEntity:DbEntity{1, "Hello World"}}.Name())
}

func TestNewCategory(t *testing.T) {
	category, err := NewCategory(123, "hello")
	if assert.Nil(t, err) {
		assert.Equal(t, EntityID(123), category.Id)
		assert.Equal(t, "HELLO", category.DbName)
		assert.Empty(t, category.Commodities)
	}
}
