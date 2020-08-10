package main

// TradeOutcome represents the purchase and up-sale of a commodity, intended
// to be itemized as part of a TradeHop (facility -> facility).
type TradeOutcome struct {
	// Commodity is the item to transact.
	Commodity *Commodity
	// CostCr is the initial purchase price.
	CostCr int64
	// GainCr is how many credits will be gained on selling the item.
	GainCr int64
	// Supply indicates how many units the selling facility is expected to have.
	Supply int
	// SupplyLevel indicates how actively the selling facility stocks this item.
	SupplyLevel int
	// Demand indicates how many units the purchasing facility is expected to want.
	Demand int
	// DemandLevel indicates how actively the purchasing facility is acquiring this item.
	DemandLevel int
	// SrcAge is how old in seconds the seller's data was when this outcome was calculated.
	SrcAge int
	// DstAge is how old in seconds the buyer's data was when this outcome was calculated.
	DstAge int
}
