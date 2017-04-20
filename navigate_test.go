package main

import (
	"testing"
)

func TestRatio(t *testing.T) {
	tests := []struct {
		x, y int
		want float32
	}{
		{1, 2, 0.5},
		{2, 1, 2},
		{1920, 1080, 1.7777777777},
		{1, 0, 0.0},
	}

	for _, test := range tests {
		got, _ := Ratio(test.x, test.y)
		if got != test.want {
			t.Errorf("got: %v, want: %v", got, test.want)
		}
	}
}
