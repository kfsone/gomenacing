package main

// EntityID is the underlying type used to store identifiers in the database.
type EntityID uint32

// FDevID captures an frontier development id.
type FDevID uint32

// Sectors divide the galaxy into grids 32lyx32lyx32ly, which allows us to
// quickly find starts within small distances of each other.
type SectorKey struct {
	X, Y, Z int
}
