package oor

// #cgo pkg-config: opencv
// #include "oor.h"
import "C"
import (
	"errors"
	"unsafe"

	"github.com/lazywei/go-opencv/opencv"
)

// GoOOR Wrapper for C++ code in Go for object detection.
type GoOOR struct {
	coor C.OOR
}

// New wraps our C++ binding to a Go struct.
func New() GoOOR {
	return GoOOR{coor: C.OORInit()}
}

// Free deletes the GoOOR object.
func (g GoOOR) Free() {
	C.OORFree(g.coor)
}

// DetectEllipses finds ellipses in the current image from a video stream.
// The return slice first four values are the rectangle information which
// surround the larges Elipses and the last two are the center of the image.
func (g GoOOR) DetectEllipses(img *opencv.Mat) ([]int, error) {
	var c *C.int = C.DetectEllipses(g.coor, unsafe.Pointer(img))
	// Store the int* which is an array
	pc := unsafe.Pointer(c)

	if pc == nil {
		return nil, errors.New("c code return null")
	}
	defer C.free(pc) // free memory

	// The length of the array
	length := 6
	// Convert c to pointer to an array, and then slice it.
	cSlice := (*[1 << 4]C.int)(pc)[:length:length]

	// Make an empty slice and add elements from cSlice
	s := make([]int, 0, length)
	for _, v := range cSlice {
		s = append(s, int(v))
	}

	if s[0] == 0 && s[1] == 0 && s[2] == 0 && s[3] == 0 {
		return nil, errors.New("no ellipses detected")
	}

	return s, nil
}
