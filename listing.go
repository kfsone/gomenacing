package main

import (
	"fmt"
	"time"
)

type Listing struct {
	CommodityID  EntityID
	Supply       uint32
	StationPays  uint32
	Demand       uint32
	StationAsks  uint32
	TimestampUtc time.Time
}

func NewListing(commodityID EntityID, supply uint32, stationPays uint32, demand uint32, stationAsks uint32, timestampUtc time.Time) *Listing {
	return &Listing{CommodityID: commodityID, Supply: supply, StationPays: stationPays, Demand: demand, StationAsks: stationAsks, TimestampUtc: timestampUtc}
}

func (l *Listing) GetId() uint32 {
	return uint32(l.CommodityID)
}

func (l *Listing) GetDbId(f *Facility) string {
	return fmt.Sprintf("%06x%4x", f.GetId(), l.GetId())
}
