package main

import (
	"fmt"
	"testing"
)

func TestCenter(t *testing.T) {
	tests := []struct {
		x, y, width, height int
		want                Point
	}{
		{700, 1412, 843, 823, Point{1122, 1824}},
		{383, 1044, 1353, 1331, Point{1060, 1710}},
		{1137, 1277, 825, 977, Point{1550, 1766}},
		{1096, 1322, 1156, 1157, Point{1674, 1901}},
		{799, 1871, 549, 570, Point{1074, 2156}},
	}

	for _, test := range tests {
		got := Center(test.x, test.y, test.width, test.height)
		if got != test.want {
			t.Errorf("Center failed: got: %v, test.want: %v", got, test.want)
		}
	}

}

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

func TestDirection(t *testing.T) {
	tests := []struct {
		ratio         float32
		wantDirection int
		wantBadness   float32
	}{
		{1920.0 / 1080.0, VIRTICAL, 0.777777777},
		{1080.0 / 1920.0, HORIZONTAL, 0.4375},
		{1920.0 / 1840.0, STAY, 0.0},
		{1920.0 / 2000.0, STAY, 0.0},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test:%v", i), func(t *testing.T) {
			d, b := Direction(test.ratio)
			if d != test.wantDirection {
				t.Fatalf("Direciton failed: got: %v, want: %v", d, test.wantDirection)
			}
			if b != test.wantBadness {
				t.Errorf("Badness failed: got: %v, want: %v", b, test.wantBadness)
			}
		})
	}
}
