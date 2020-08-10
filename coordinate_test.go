package main

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinate_Distance(t *testing.T) {
	t.Run("Zero relative", func(t *testing.T) {
		zero := Coordinate{0.0, 0.0, 0.0}
		assert.Equal(t, Square(0.00), zero.Distance(zero))
		assert.Equal(t, Square(3.00), zero.Distance(Coordinate{-1.0, -1.0, -1.0}))
		assert.Equal(t, Square(13.25), zero.Distance(Coordinate{2.0, -3.0, 0.5}))
	})
	t.Run("Non-zero relative", func(t *testing.T) {
		lhs := Coordinate{-3.0, 4.2, 5.5}
		assert.Equal(t, Square(0.00), lhs.Distance(lhs))
		assert.Equal(t, Square(1.00), lhs.Distance(Coordinate{-2.0, 4.2, 5.5}))
		assert.Equal(t, Square(2.00), lhs.Distance(Coordinate{-2.0, 5.2, 5.5}))
		assert.Equal(t, Square(3.00), lhs.Distance(Coordinate{-2.0, 5.2, 4.5}))
		assert.Equal(t, Square(56.89), lhs.Distance(Coordinate{-0.0, 0.0, 0.0}))
	})
}

func TestCoordinate_SectorKey(t *testing.T) {
	assert.Equal(t, SectorKey{}, Coordinate{}.SectorKey())
	assert.Equal(t, SectorKey{-1, 0, 0}, Coordinate{-1, 1, 2}.SectorKey())
	assert.Equal(t, SectorKey{1, -2, 0}, Coordinate{1 << 5, -2 << 5, (1 << 5) - 1}.SectorKey())
}

func TestNewCoordinateFromStrings(t *testing.T) {
	var (
		c   Coordinate
		err error
	)
	c, err = NewCoordinateFromStrings("0", "0", "0")
	require.Nil(t, err)
	assert.Equal(t, Coordinate{}, c)

	c, err = NewCoordinateFromStrings("0.0", "-0.0", "0.0")
	require.Nil(t, err)
	assert.Equal(t, Coordinate{}, c)

	c, err = NewCoordinateFromStrings("-00001.23", "+123.023", "-999999.999999")
	require.Nil(t, err)
	assert.Equal(t, Coordinate{-1.23, 123.023, -999999.999999}, c)

	// Now check errors
	c, err = NewCoordinateFromStrings("", "", "")
	assert.NotNil(t, err)

	c, err = NewCoordinateFromStrings("1", "", "")
	assert.NotNil(t, err)

	c, err = NewCoordinateFromStrings("1", "2", "")
	assert.NotNil(t, err)
}
