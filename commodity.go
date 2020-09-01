package main

import (
	"fmt"
	gom "github.com/kfsone/gomenacing/pkg/gomschema"
)

// Commodity is a representation of a tradeable item.
type Commodity struct {
	DbEntity
	CategoryID      gom.Commodity_Category // Which category the commodity is part of.
	IsRare          bool                   // True for rare items.
	IsNonMarketable bool                   // Can we sell this?
	AverageCr       uint32                 // Average market price.
}

// NewCommodity constructs a Commodity instance based on the input parameters.
func NewCommodity(entity DbEntity, category gom.Commodity_Category, rare, nonMarketable bool, avgPrice uint32) (*Commodity, error) {
	return &Commodity{
		DbEntity:        entity,
		CategoryID:      category,
		IsRare:          rare,
		IsNonMarketable: nonMarketable,
		AverageCr:       avgPrice,
	}, nil
}

// Name is the user-facing name adjusted for level of detail.
func (c *Commodity) Name() string {
	return c.DbName
}

func (c *Commodity) GetDbId() string {
	return fmt.Sprintf("%04x", c.DbEntity.ID)
}
