package main

import "math"

// Square is a representation of the square of some value.
type Square interface {
	Root() float64
}

type SquareFloat float64
type SquareInt int64

func NewSquareFloat(value float64) SquareFloat {
	return SquareFloat(value * value)
}

func NewSquareInt(value int64) SquareInt {
	return SquareInt(value * value)
}

// Root returns the square root of a Square value.
func (s SquareFloat) Root() float64 {
	return math.Sqrt(float64(s))
}

func (s SquareInt) Root() float64 {
	return math.Sqrt(float64(s))
}
