package main

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

// Commodity is a representation of a tradeable item.
type Commodity struct {
	DbEntity
	CategoryId EntityID // Which category the commodity is part of.
	AvgPrice   int64    // Average market price.
	FDevId     FDevID   // FrontierDev internal ID.
	Unsellable bool     // Can we sell this?
}

// MissingCategory is an error raised when attempting to create a commodity without a category.
var MissingCategory = errors.New("called NewCommodity with nil category")

// Fields used to describe commodity in eddb.
//var commodityFields = []string{"id", "name", "category.id", "is_non_marketable", "average_price", "ed_id"}
func NewCommodityFromJsonMap(json gjson.Result) (*Commodity, error) {
	if !json.IsObject() {
		return nil, fmt.Errorf("unrecognized json value: %+v", json)
	}
	data := json.Map()
	id, name := data["id"].Int(), data["name"].String()
	entity, err := NewDbEntity(id, name)
	if err != nil {
		return nil, err
	}
	categoryId := data["category_id"].Int()
	nonMarketable, avgPrice, fdevId := data["is_non_marketable"].Bool(), data["average_price"].Int(), data["ed_id"].Int()
	item, err := NewCommodity(entity, categoryId, avgPrice, fdevId)
	if item != nil {
		item.Unsellable = nonMarketable
	}
	return item, err
}

// NewCommodity constructs a Commodity instance based on the input parameters.
func NewCommodity(entity DbEntity, categoryId int64, avgPrice int64, fdevId int64) (*Commodity, error) {
	if categoryId <= 0 || categoryId >= 1<<32 {
		return nil, fmt.Errorf("invalid commodity category id: %d", categoryId)
	}
	return &Commodity{
		DbEntity:   entity,
		CategoryId: EntityID(categoryId),
		AvgPrice:   avgPrice,
		FDevId:     FDevID(fdevId),
	}, nil
}

// Name is the user-facing name adjusted for level of detail.
func (i Commodity) Name() string {
	return i.DbName
}
