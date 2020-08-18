package main

import "strings"

// Category is a grouping of Commodities, e.g. Agricultural Products.
type Category struct {
	DbEntity
	// Commodities is the list of items within this group.
	Commodities []*Commodity
}

// NewCategory creates a new category entry.
func NewCategory(id int64, name string) (category Category, err error) {
	entity, err := NewDbEntity(id, strings.ToUpper(name))
	if err == nil {
		category.DbEntity = entity
	}
	return
}

// Name returns the user-facing name of this Category.
func (c Category) Name() string {
	return c.DbName
}
