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
}

// MissingCategory is an error raised when attempting to create a commodity without a category.
var MissingCategory = errors.New("called NewCommodity with nil category")

// Fields used to describe commodity in eddb.
//var commodityFields = []string{"id", "name", "category.id", "category.name", "is_non_marketable", "average_price", "ed_id"}
func NewCommodityFromJson(json []gjson.Result) (*Commodity, error) {
	id, name, categoryId := json[0].Int(), json[1].String(), json[2].Int()
	nonMarketable, avgPrice, fdevId := json[3].Bool(), json[4].Int(), json[5].Int()
	if nonMarketable {
		return nil, fmt.Errorf("non-marketable item included: %s (#%d)", name, id)
	}
	item, err := NewCommodity(id, name, categoryId, avgPrice, fdevId)
	return item, err
}

// NewCommodity constructs a Commodity instance based on the input parameters.
func NewCommodity(id int64, dbName string, categoryId int64, avgPrice int64, fdevId int64) (*Commodity, error) {
	dbEntity, err := NewDbEntity(id, dbName)
	if err != nil {
		return nil, err
	}
	if categoryId <= 0 || categoryId >= 1 << 32 {
		return nil, fmt.Errorf("invalid commodity category id: %d", categoryId)
	}
	return &Commodity{
		DbEntity:   dbEntity,
		CategoryId: EntityID(categoryId),
		AvgPrice:   avgPrice,
		FDevId:     FDevID(fdevId),
	}, nil
}

// Name is the user-facing name adjusted for level of detail.
func (i Commodity) Name() string {
	return i.DbName
}
