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
			6,
			[]int{700, 1412, 843, 823, 1560, 2080},
		},
		{
			opencv.LoadImageM("test2.jpg", 1),
			6,
			[]int{383, 1044, 1353, 1331, 1560, 2080},
		},
		{
			opencv.LoadImageM("test3.jpg", 1),
			6,
			[]int{1137, 1277, 825, 977, 1560, 2080},
		},
		{
			opencv.LoadImageM("test4.jpg", 1),
			6,
			[]int{1096, 1322, 1156, 1157, 1560, 2080},
		},
		{
			opencv.LoadImageM("test5.jpg", 1),
			6,
			[]int{799, 1871, 549, 570, 1560, 2080},
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
		t.Logf("got: %v", got)
	}
}
