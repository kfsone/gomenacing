package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinate_Distance(t *testing.T) {
	t.Run("Zero relative", func(t *testing.T) {
		zero := Coordinate{}
		assert.Equal(t, Square(0.00), zero.Distance(zero))
		assert.Equal(t, Square(3.00), zero.Distance(Coordinate{X: -1.0, Y: -1.0, Z: -1.0}))
		assert.Equal(t, Square(13.25), zero.Distance(Coordinate{X: 2.0, Y: -3.0, Z: 0.5}))
	})
	t.Run("Non-zero relative", func(t *testing.T) {
		lhs := Coordinate{X: -3.0, Y: 4.2, Z: 5.5}
		assert.Equal(t, Square(0.00), lhs.Distance(lhs))
		assert.Equal(t, Square(1.00), lhs.Distance(Coordinate{X: -2.0, Y: 4.2, Z: 5.5}))
		assert.Equal(t, Square(2.00), lhs.Distance(Coordinate{X: -2.0, Y: 5.2, Z: 5.5}))
		assert.Equal(t, Square(3.00), lhs.Distance(Coordinate{X: -2.0, Y: 5.2, Z: 4.5}))
		assert.Equal(t, Square(56.89), lhs.Distance(Coordinate{X: -0.0}))
	})
}

func TestCoordinate_SectorKey(t *testing.T) {
	assert.Equal(t, SectorKey{}, Coordinate{}.SectorKey())
	assert.Equal(t, SectorKey{X: -1}, Coordinate{X: -1, Y: 1, Z: 2}.SectorKey())
	assert.Equal(t, SectorKey{X: 1, Y: -2}, Coordinate{X: 1 << 5, Y: -2 << 5, Z: (1 << 5) - 1}.SectorKey())
}
