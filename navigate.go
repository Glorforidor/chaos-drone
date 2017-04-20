package main

import "errors"

var ErrDivByZero = errors.New("divide by zero")

// Ratio return the ratio between x and y. If the value is above 1 then x is
// greater than y and below 0 y is greater than x.
// An error of type ErrDivByZero is returned if y is zero or less.
func Ratio(x, y int) (float32, error) {
	if y <= 0 {
		return 0.0, ErrDivByZero
	}

	return float32(x) / float32(y), nil
}
