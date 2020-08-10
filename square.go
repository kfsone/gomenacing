package main

import "math"

// Representation of the square of some value.
type Square float64

// Root returns the square root of a Square value.
func (s Square) Root() float64 {
	return math.Sqrt(float64(s))
}
