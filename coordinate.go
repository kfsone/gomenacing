package main

// Coordinate is a type representing a 3d coordinate within the galaxy.
type Coordinate struct {
	X, Y, Z float64
}

func (c *Coordinate) Coordinate() *Coordinate {
	return c
}

// SectorKey maps a Coordinate into a SectorKey.
func (c Coordinate) SectorKey() SectorKey {
	return SectorKey{int(c.X) >> 5, int(c.Y) >> 5, int(c.Z) >> 5}
}

type Positioned interface {
	Coordinate() *Coordinate
}

// Distance calculates the distance^2 between two positioned objects.
func Distance(l Positioned, r Positioned) Square {
	lhs, rhs := l.Coordinate(), r.Coordinate()
	xDelta, yDelta, zDelta := lhs.X-rhs.X, lhs.Y-rhs.Y, lhs.Z-rhs.Z
	return Square(xDelta*xDelta + yDelta*yDelta + zDelta*zDelta)
}
