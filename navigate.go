package main

import (
	"errors"
	"math"
)

const (
	VIRTICAL   = iota
	HORIZONTAL = iota
	STAY       = iota
)

// ErrDivByZero as the name imply an divide by zero error.
var ErrDivByZero = errors.New("divide by zero")

type Point struct {
	X, Y int
}

func round(x float64) int {
	if x < 0 {
		return int(math.Ceil(x - 0.5))
	}
	return int(math.Floor(x + 0.5))
}

func Center(x, y, width, height int) Point {
	x = round(float64((x + x + width)) * 0.5)
	y = round(float64((y + y + height)) * 0.5)
	return Point{x, y}
}

// Ratio return the ratio between x and y. If the value is above 1 then x is
// greater than y and below 0 y is greater than x.
// An error of type ErrDivByZero is returned if y is zero or less.
func Ratio(x, y int) (float32, error) {
	if y <= 0 {
		return 0.0, ErrDivByZero
	}

	return float32(x) / float32(y), nil
}

// Direction checks the ratio to detect the drones possition relative to the
// object.
func Direction(ratio float32) (direction int, badness float32) {
	move := 1 - ratio
	badness = float32(math.Abs(float64(move)))
	if 1.0 < badness {
		badness = 1.0
	}
	if -0.05 < move && move < 0.05 {
		return STAY, 0.0
	} else if move < -0.05 {
		return VIRTICAL, badness
	} else {
		return HORIZONTAL, badness
	}
}

func Placement() {}
