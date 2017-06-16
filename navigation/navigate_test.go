package navigation

import (
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

func TestPlacement(t *testing.T) {
	// TODO: Take account for threshold.
	tests := []struct {
		x, y, rx, ry int
		want         int
	}{
		{1, 1, 2, 2, Left},
		{2, 1, 2, 2, Up},
		{2, 2, 2, 2, OnTarget},
		{2, 2, 1, 2, Right},
		{2, 2, 2, 1, Down},
	}

	for _, test := range tests {
		got := Placement(test.x, test.y, test.rx, test.ry)
		if got != test.want {
			t.Errorf("Wrong Placement: got: %v, want: %v", got, test.want)
		}
	}
}
