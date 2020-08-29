package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSquare_Root(t *testing.T) {
	s := Square(9.0)
	assert.Equal(t, s.Root(), 3.0)
}
