package main

import "time"

// Timestamped is an interface that provides the GetTimestampUtc() function.
type Timestamped interface {
	GetTimestampUtc() uint64
}

// getTimestamp converts a Timestamped value into a unix time.
func getTimestamp(t Timestamped) time.Time {
	return time.Unix(int64(t.GetTimestampUtc()), 0)
}
