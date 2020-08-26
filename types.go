package main

// EntityID is the underlying type used to store identifiers in the database.
type EntityID uint32

// FDevID captures an frontier development id.
type FDevID uint32

// SectorKey divides the galaxy into grids 32lyx32lyx32ly, which allows us to
// quickly find starts within small distances of each other.
type SectorKey struct {
	X, Y, Z int
}

// TestDir is a test helper: Create a temporary directory to operate in.
type TestDir string

func (td TestDir) String() string {
	return string(td)
}

// Path returns the string path of the test directory.
func (td TestDir) Path() string {
	return string(td)
}
