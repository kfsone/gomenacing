package main

const SectorShift = 7
const SectorWidth = 1 << SectorShift

// Coordinate is a type representing a 3d coordinate within the galaxy.
type Coordinate struct {
	X, Y, Z float64
}

func (c Coordinate) Coordinate() *Coordinate {
	return &c
}

// SectorKey maps a Coordinate into a SectorKey.
func (c Coordinate) SectorKey() SectorKey {
	return SectorKey{int(c.X) >> SectorShift, int(c.Y) >> SectorShift, int(c.Z) >> SectorShift}
}

type Positioned interface {
	Coordinate() *Coordinate
}

// Distance calculates the distance^2 between two positioned objects.
func Distance(l Positioned, r Positioned) SquareFloat {
	lhs, rhs := l.Coordinate(), r.Coordinate()
	xDelta, yDelta, zDelta := NewSquareFloat(lhs.X-rhs.X), NewSquareFloat(lhs.Y-rhs.Y), NewSquareFloat(lhs.Z-rhs.Z)
	return xDelta + yDelta + zDelta
}
