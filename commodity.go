package main

import "errors"

// Commodity is a representation of a tradable item.
type Commodity struct {
	DbEntity
	Category *Category // Which category the commodity is part of.
	FullName string    // The UI name of the product.
	AvgPrice int       // Average market price.
	FDevId   FDevID    // FrontierDev internal ID.
}

// MissingCategory is an error raised when attempting to create a commodity without a category.
var MissingCategory = errors.New("called NewCommodity with nil category")

// NewCommodity constructs a Commodity instance based on the input parameters.
func NewCommodity(id EntityID, dbName string, category *Category, fullName string, avgPrice int, fdevId FDevID) (*Commodity, error) {
	if category == nil {
		return nil, MissingCategory
	}
	return &Commodity{
		DbEntity: DbEntity{id, dbName},
		Category: category,
		FullName: fullName,
		AvgPrice: avgPrice,
		FDevId:   fdevId,
	}, nil
}

// Name is the user-facing name adjusted for level of detail.
func (i Commodity) Name(detail int) string {
	if detail > 0 {
		return i.FullName
	} else {
		return i.DbName
	}
}
