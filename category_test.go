package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategory_Name(t *testing.T) {
	c := NewCategory(1112, "Hello world")
	assert.Equal(t, c.Name(0), "HELLO WORLD")

	c = NewCategory(1234, "pizza PIE")
	assert.Equal(t, c.Name(100), "PIZZA PIE")
}

func TestNewCategory(t *testing.T) {
	category := NewCategory(123, "hello")
	assert.Equal(t, EntityID(123), category.Id)
	assert.Equal(t, "hello", category.DbName)
	assert.Empty(t, category.Commodities)
}
