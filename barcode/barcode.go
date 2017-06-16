package barcode

import (
	"image"
	_ "image/jpeg" //For JPEG image support
	_ "image/png"  //For PNG image support
	"math"

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

var scanner *barcode.ImageScanner

const ellipseYOffset = 3.8 // how many QR code heights do we need to offset our center point to get to the ellipse ring?

// Init initalizes the barcode resources
func Init() {
	scanner = barcode.NewScanner().SetEnabledAll(true)
}

// QRRawData returns QR data from a image. An error is returned if image can't be
// decoded or there is a problem with scanning the image.
func QRRawData(m image.Image) ([]*barcode.Symbol, error) {
	i := barcode.NewImage(m)

	symbols, err := scanner.ScanImage(i)
	if err != nil {
		return nil, errors.Wrap(err, "could not scan image")
	}

	return symbols, nil
}

// GetEllipseOverQR returns an XY coordinate located over the given QR code,
// if any. If the QR code is not found, nil is returned instead.
func GetEllipseOverQR(camImg *cv.IplImage, qrText string) ([]cv.Point, error) {
	camImg.ToImage()
	// We need to make the ToImage() call twice, due to shenenigans
	symbols, err := QRRawData(camImg.ToImage())
	if err != nil {
		return nil, err
	}

	for _, s := range symbols {
		if s.Data == qrText {
			var mx, ly, uy int
			for i := 0; i < 4; i++ {
				x := s.Boundary[i].X
				y := s.Boundary[i].Y
				mx += x
				if i == 0 || ly > y {
					ly = y
				}
				if i == 0 || uy < y {
					uy = y
				}
			}
			var upLen = float64(uy-ly) * float64(ellipseYOffset)
			var tilt = math.Atan2(
				float64(s.Boundary[0].X)-float64(mx)*0.25,
				float64(s.Boundary[0].Y)-float64(ly)+float64(uy-ly)*0.5,
			) + math.Pi*(135.0/180.0)
			return []cv.Point{
				cv.Point{X: int(float64(mx) * 0.25), Y: int(float64(ly) + float64(uy-ly)*0.5)},
				cv.Point{X: int(float64(mx)*0.25 - math.Cos(tilt)*upLen), Y: int(float64(ly) + float64(uy-ly)*0.5 - math.Sin(tilt)*upLen)},
			}, nil
		}
	}

	return nil, nil
}
