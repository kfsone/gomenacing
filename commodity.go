package main

import (
	"errors"
	"fmt"

	"github.com/tidwall/gjson"
)

// Commodity is a representation of a tradeable item.
type Commodity struct {
	DbEntity
	CategoryID EntityID // Which category the commodity is part of.
	AvgPrice   int64    // Average market price.
	FDevID     FDevID   // FrontierDev internal ID.
	Unsellable bool     // Can we sell this?
}

// ErrMissingCategory is an error raised when attempting to create a commodity without a category.
var ErrMissingCategory = errors.New("called NewCommodity with nil category")

// Fields used to describe commodity in eddb.
//var commodityFields = []string{"id", "name", "category.id", "is_non_marketable", "average_price", "ed_id"}

// NewCommodityFromJSONMap creates a Commodity instance from a json map.
func NewCommodityFromJSONMap(json gjson.Result) (*Commodity, error) {
	if !json.IsObject() {
		return nil, fmt.Errorf("unrecognized json value: %+v", json)
	}
	data := json.Map()
	id, name := data["id"].Int(), data["name"].String()
	entity, err := NewDbEntity(id, name)
	if err != nil {
		return nil, err
	}
	categoryID := data["category_id"].Int()
	nonMarketable, avgPrice, fdevID := data["is_non_marketable"].Bool(), data["average_price"].Int(), data["ed_id"].Int()
	item, err := NewCommodity(entity, categoryID, avgPrice, fdevID)
	if item != nil {
		item.Unsellable = nonMarketable
	}
	return item, err
}

// NewCommodity constructs a Commodity instance based on the input parameters.
func NewCommodity(entity DbEntity, categoryID int64, avgPrice int64, fdevID int64) (*Commodity, error) {
	if categoryID <= 0 || categoryID >= 1<<32 {
		return nil, fmt.Errorf("invalid commodity category id: %d", categoryID)
	}
	return &Commodity{
		DbEntity:   entity,
		CategoryID: EntityID(categoryID),
		AvgPrice:   avgPrice,
		FDevID:     FDevID(fdevID),
	}, nil
}

// Name is the user-facing name adjusted for level of detail.
func (i Commodity) Name() string {
	return i.DbName
}
