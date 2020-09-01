package main

// Coordinate is a type representing a 3d coordinate within the galaxy.
type Coordinate struct {
	X, Y, Z float64
}

// SectorKey maps a Coordinate into a SectorKey.
func (c Coordinate) SectorKey() SectorKey {
	return SectorKey{int(c.X) >> 5, int(c.Y) >> 5, int(c.Z) >> 5}
}

// Distance calculates the distance^2 between two coordinates.
func (c Coordinate) Distance(rhs Coordinate) Square {
	x, y, z := c.X-rhs.X, c.Y-rhs.Y, c.Z-rhs.Z
	return Square(x*x + y*y + z*z)
}
