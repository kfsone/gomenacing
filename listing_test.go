package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListing_GetDbId(t *testing.T) {
	facility := Facility{DbEntity: DbEntity{ID: 0x018019, DbName: "SOL"}}
	listing := Listing{CommodityID: 0x1234}

	assert.Equal(t, "0180191234", listing.GetDbId(&facility))
}

func TestListing_GetId(t *testing.T) {
	listing := Listing{CommodityID: 456}
	assert.Equal(t, uint32(456), listing.GetId())
}

func TestListing_GetTimestampUtc(t *testing.T) {
	listing := Listing{TimestampUtc: 4770}
	assert.Equal(t, uint64(4770), listing.GetTimestampUtc())
}
