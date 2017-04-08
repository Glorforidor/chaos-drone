package oor

import (
	"testing"

	"github.com/lazywei/go-opencv/opencv"
)

func TestDetectEllipses(t *testing.T) {
	tests := []struct {
		image  *opencv.Mat
		length int
		want   []int
	}{
		{
			opencv.LoadImageM("test0.png", 1),
			0,
			nil,
		},
		{
			opencv.LoadImageM("test1.jpg", 1),
			4,
			[]int{1122, 1824, 1560, 2080},
		},
		{
			opencv.LoadImageM("test2.jpg", 1),
			4,
			[]int{1060, 1710, 1560, 2080},
		},
		{
			opencv.LoadImageM("test3.jpg", 1),
			4,
			[]int{1550, 1766, 1560, 2080},
		},
		{
			opencv.LoadImageM("test4.jpg", 1),
			4,
			[]int{1674, 1900, 1560, 2080},
		},
		{
			opencv.LoadImageM("test5.jpg", 1),
			4,
			[]int{1074, 2156, 1560, 2080},
		},
	}

	o := New()
	defer o.Free()
	for _, test := range tests {
		got, _ := o.DetectEllipses(test.image)
		if len(got) != test.length {
			t.Errorf("Length mismatch: got: %v, want: %v", len(got), len(test.want))
		}
		test.want = test.want[:len(got)] // eliminate bound checks
		for i, v := range got {
			if v != test.want[i] {
				t.Errorf("Value mismatch: got: %v, want: %v", v, test.want[i])
			}
		}
	}
}
