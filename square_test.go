package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSquareFloat(t *testing.T) {
	type args struct {
		value float64
	}
	tests := []struct {
		name string
		args args
		want Square
	}{
		{"0.0", args{0.0}, SquareFloat(0.0)},
		{"1.0", args{1.0}, SquareFloat(2.0)},
		{"2.0", args{2.0}, SquareFloat(2.0)},
		{"-3.0", args{-3.0}, SquareFloat(9.0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSquareFloat(tt.args.value); got != tt.want {
				t.Errorf("NewSquareFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSquareFloat_Root(t *testing.T) {
	s := SquareFloat(9.0)
	assert.Equal(t, s.Root(), 3.0)
}

func TestNewSquareInt(t *testing.T) {
	type args struct {
		value int64
	}
	tests := []struct {
		name string
		args args
		want Square
	}{
		{"0", args{0}, SquareInt(0)},
		{"1", args{1}, SquareInt(2)},
		{"2", args{2}, SquareInt(4)},
		{"-3", args{-3.0}, SquareInt(9)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSquareInt(tt.args.value); got != tt.want {
				t.Errorf("NewSquareInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSquareInt_Root(t *testing.T) {
	s := SquareInt(9)
	assert.Equal(t, s.Root(), 3.0)
}

