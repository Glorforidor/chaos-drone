package oor

import (
	"io/ioutil"
	"math"
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

func Test(t *testing.T) {
	o := New()
	defer o.Free()
	dirPics := "drone_pics/"
	files, _ := ioutil.ReadDir(dirPics)
	image := opencv.LoadImageM(dirPics+files[2].Name(), opencv.CV_LOAD_IMAGE_COLOR)
	i, err := o.DetectEllipses(image)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(i)
	t.Log(Center(i[0], i[1], i[2], i[3]))
}
func Center(x, y, width, height int) Point {
	x = round(float64((x + x + width)) * 0.5)
	y = round(float64((y + y + height)) * 0.5)
	return Point{x, y}
}

// Point contains a X and Y coordinate.
type Point struct {
	X, Y int
}

func round(x float64) int {
	if x < 0 {
		return int(math.Ceil(x - 0.5))
	}
	return int(math.Floor(x + 0.5))
}
