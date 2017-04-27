package main

import (
	"errors"
	"math"
)

const (
	VIRTICAL = iota
	HORIZONTAL
	STAY
)

const (
	NORTH = iota
	WEST
	SOUTH
	EAST
	ON_TARGET
)

// ErrDivByZero as the name imply an divide by zero error.
var ErrDivByZero = errors.New("divide by zero")

type Point struct {
	X, Y int
}

// round rounds to nearest integer.
func round(x float64) int {
	if x < 0 {
		return int(math.Ceil(x - 0.5))
	}
	return int(math.Floor(x + 0.5))
}

// Center calculates the center rectangle.
func Center(x, y, width, height int) Point {
	x = round(float64((x + x + width)) * 0.5)
	y = round(float64((y + y + height)) * 0.5)
	return Point{x, y}
}

// Ratio return the ratio between x and y. If the value is above 1 then x is
// greater than y and below 0 y is greater than x.
// An error of type ErrDivByZero is returned if y is zero or less.
func Ratio(x, y int) (float64, error) {
	if y <= 0 {
		return 0.0, ErrDivByZero
	}

	return float64(x) / float64(y), nil
}

// Direction checks the ratio to detect the drones possition relative to the
// object.
func Direction(ratio float64) (direction int, badness float64) {
	move := 1 - ratio
	direction = STAY
	badness = math.Abs(math.Log10(ratio))
	if 1.0 < badness {
		badness = 1.0
	}
	if -0.05 < move && move < 0.05 {
		badness = 0.0
	} else if move < -0.05 {
		direction = VIRTICAL
	} else {
		direction = HORIZONTAL
	}

	return direction, badness
}

// Placement calculates the position of the drone relative to the center of an
// ring.
func Placement(x, y, rx, ry int) int {
	placement := ON_TARGET
	switch {
	case x < rx:
		placement = WEST
	case x > rx:
		placement = EAST
	case y < ry:
		placement = NORTH
	case y > ry:
		placement = SOUTH
	}

	return placement
}
