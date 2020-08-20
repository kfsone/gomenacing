package main

import (
	"time"
)

type Listing struct {
	CommodityId EntityID  `json:"id"`
	Supply      int64     `json:"sup"`
	Demand      int64     `json:"dem"`
	BuyingAt    int64     `json:"buy"`
	SellingAt   int64     `json:"sell"`
	Recorded    time.Time `json:"time"`
}

const (
	EddbListings = "listings.csv"
)

//func NewListingFromArray(array []string) (*Listing, error) {
//	if len(array) != 6 {
//		return nil, errors.New("invalid listing array")
//	}
//	commodityId, err := strconv.ParseInt(array[0], 10, 64)
//	if err != nil {
//		return nil, err
//	}
//	if commodityId <= 0 || commodityId >= 1<<32 {
//		return nil, errors.New("invalid commodity id")
//	}
//	return &Listing{EntityID(commodityId),
//		array[1], array[2], array[3], array[4],
//		time.Unix(array[5], 0)}, nil
//}
