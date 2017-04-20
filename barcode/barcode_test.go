package barcode

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"testing"
)

func TestQRData(t *testing.T) {
	tests := []struct {
		img          string
		wantData     string
		wantBoundary []image.Point
	}{
		{"demo.jpg", "P.02", []image.Point{{1248, 1694}, {1261, 2246}, {1805, 2232}, {1792, 1683}}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Image:%v", test.img), func(t *testing.T) {
			img, err := ioutil.ReadFile(test.img)
			if err != nil {
				t.Fatalf("Failed reading file: %v", err)
			}

			symbols, err := QRData(bytes.NewReader(img))
			if err != nil {
				t.Fatalf("Failed getting QR data")
			}

			for _, s := range symbols {
				if s.Data != test.wantData {
					t.Errorf("Data mismatch: got: %v, want: %v", s.Data, test.wantData)
				}
				if !compareSlice(s.Boundary, test.wantBoundary) {
					t.Errorf("Boundary mismatch: got: %v, want: %v", s.Boundary, test.wantBoundary)
				}
				t.Logf("Data: %v, Boundary: %v", s.Data, s.Boundary)
			}
		})
	}
}

func compareSlice(s1, s2 []image.Point) bool {
	if len(s1) != len(s2) {
		return false
	}

	s2 = s2[:len(s2)] // eliminate bound check
	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}

	return true
}
