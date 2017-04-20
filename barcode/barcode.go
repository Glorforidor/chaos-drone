package barcode

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	cv "github.com/lazywei/go-opencv/opencv"
	"github.com/pkg/errors"

	"gopkg.in/bieber/barcode.v0"
)

var (
	font = cv.InitFont(
		cv.CV_FONT_HERSHEY_SIMPLEX, float32(0.65), float32(0.65),
		float32(1.0), 2, 8,
	)
	fontColor      = cv.NewScalar(0, 0, 255, 0)
	fontColorWhite = cv.NewScalar(255, 255, 255, 0)
)

func nextIndex(i int, max int) int {
	if i < max {
		return i + 1
	}
	return 0
}

func QRData(img io.Reader) ([]*barcode.Symbol, error) {
	m, _, err := image.Decode(img)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode image")
	}

	i := barcode.NewImage(m)
	scan := barcode.NewScanner().SetEnabledAll(true)

	symbols, err := scan.ScanImage(i)
	if err != nil {
		return nil, errors.Wrap(err, "could not scan image")
	}

	return symbols, nil
}

func QRScan(camImg *cv.IplImage) error {
	img := barcode.NewImage(camImg.ToImage())
	scanner := barcode.NewScanner().SetEnabledAll(true)

	symbols, err := scanner.ScanImage(img)
	if err != nil {
		return errors.Wrap(err, "could not scan barcode in image")
	}
	for _, s := range symbols {
		// Debug purpose
		// fmt.Println(s.Type.Name(), s.Data, s.Quality, s.Boundary)
		for i := 0; i < 4; i++ {
			var pt1 = cv.Point{
				X: s.Boundary[i].X,
				Y: s.Boundary[i].Y,
			}
			var pt2 = cv.Point{
				X: s.Boundary[nextIndex(i, 3)].X,
				Y: s.Boundary[nextIndex(i, 3)].Y,
			}
			cv.Line(camImg, pt1, pt2, cv.NewScalar(0, 255, 0, 0), 2, 8, 0)
		}
		for x := -1; x < 2; x++ {
			for y := -6; y < -3; y++ {
				font.PutText(camImg, s.Data,
					cv.Point{
						X: s.Boundary[1].X + 4 + x,
						Y: s.Boundary[1].Y + y,
					}, fontColorWhite)
			}
		}
		font.PutText(camImg, s.Data,
			cv.Point{
				X: s.Boundary[1].X,
				Y: s.Boundary[1].Y - 5,
			}, fontColor)
	}

	return nil
}
