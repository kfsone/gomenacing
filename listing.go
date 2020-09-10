package main

import (
	"fmt"
)

type Listing struct {
	CommodityID  EntityID
	Supply       uint32
	StationPays  uint32
	Demand       uint32
	StationAsks  uint32
	TimestampUtc uint64
}

func (l *Listing) GetId() uint32 {
	return uint32(l.CommodityID)
}

func (l *Listing) GetDbId(f *Facility) string {
	return fmt.Sprintf("%06x%4x", f.GetId(), l.GetId())
}

func (l *Listing) GetTimestampUtc() uint64 {
	return l.TimestampUtc
}
