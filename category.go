package main

import "strings"

// Category is a grouping of Commodities, e.g. Agricultural Products.
type Category struct {
	DbEntity
	// Commodities is the list of items within this group.
	Commodities []*Commodity
}

// NewCategory creates a new category entry.
func NewCategory(ID EntityID, DbName string) Category {
	return Category{DbEntity: DbEntity{ID, DbName}}
}

// Name returns the user-facing name of this Category.
func (c Category) Name(_ int) string {
	return strings.ToUpper(c.DbName)
}
