package main

import "strconv"

// Coordinate is a type representing a 3d coordinate within the galaxy.
type Coordinate struct {
	X, Y, Z float64
}

// NewCoordinateFromStrings returns a coordinate constructed from string representations of
// the x, y and z values.
func NewCoordinateFromStrings(x, y, z string) (Coordinate, error) {
	coordinate := Coordinate{}
	var err error
	if coordinate.X, err = strconv.ParseFloat(x, 64); err == nil {
		if coordinate.Y, err = strconv.ParseFloat(y, 64); err == nil {
			if coordinate.Z, err = strconv.ParseFloat(z, 64); err == nil {
				return coordinate, nil
			}
		}
	}
	return Coordinate{}, err

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
