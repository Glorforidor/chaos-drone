package barcode

import (
	"fmt"
	"image"
	_ "image/jpeg" //For JPEG image support
	_ "image/png"  //For PNG image support
	"io"
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

// QRData returns QR data from a image. An error is returned if image can't be
// decoded or there is a problem with scanning the image.
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

var scanner *barcode.ImageScanner
var ellipseYOffset float64 = 3.4 // how many QR code heights do we need to offset our center point to get to the ellipse ring?

// Init initalizes the barcode resources
func Init() {
	scanner = barcode.NewScanner().SetEnabledAll(true)
}

// QRData returns QR data from a image. An error is returned if image can't be
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
func GetEllipseOverQR(camImg *cv.IplImage, qrText string) ([]int, error) {
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
			cv.Circle(camImg, cv.Point{X: mx, Y: int(float64(ly) + float64(uy-ly)*0.5)}, 7, cv.NewScalar(255, 255, 0, 0), 4, 8, 0)
			//var upLen = float64(uy-ly) * float64(ellipseYOffset)
			//var tilt = math.Atan2(float64(s.Boundary[0].X)-float64(mx)*0.25, float64(s.Boundary[0].Y)-float64(ly)+float64(uy-ly)*0.5) + math.Pi*(135.0/180.0)
			return []int{mx, int(float64(ly) + float64(uy-ly)*0.5)}, nil //[]int{int(float64(mx)*0.25 - math.Cos(tilt)*upLen), int(float64(ly) + float64(uy-ly)*0.5 - math.Sin(tilt)*upLen)}, nil
		}
	}

	return nil, nil
}

// QRScan draws a rectangle around the QR code and the data on the image. An
// error is returned if the image could not be scanned.
func QRScan(camImg *cv.IplImage) ([]string, error) {
	img := barcode.NewImage(camImg.ToImage())
	scanner := barcode.NewScanner().SetEnabledAll(true)

	ellipseYOffset := 3.4 // how many QR code heights do we need to offset our center point to get to the ellipse ring?

	symbols, err := scanner.ScanImage(img)
	if err != nil {
		return nil, errors.Wrap(err, "could not scan barcode in image")
	}
	qrtext := make([]string, len(symbols))
	for k, s := range symbols {
		qrtext[k] = s.Data
		// Debug purpose
		// fmt.Println(s.Type.Name(), s.Data, s.Quality, s.Boundary)
		var mx, ly, uy int
		for i := 0; i < 4; i++ {
			var pt1 = cv.Point{
				X: s.Boundary[i].X,
				Y: s.Boundary[i].Y,
			}
			var pt2 = cv.Point{
				X: s.Boundary[nextIndex(i, 3)].X,
				Y: s.Boundary[nextIndex(i, 3)].Y,
			}
			mx += pt1.X
			if i == 0 || ly > pt1.Y {
				ly = pt1.Y
			}
			if i == 0 || uy < pt1.Y {
				uy = pt1.Y
			}
			if i == 0 {
				cv.Circle(camImg, pt1, 7, cv.NewScalar(0, 100, 255, 0), 4, 8, 0)
			}
			cv.Line(camImg, pt1, pt2, cv.NewScalar(0, 255, 0, 0), 2, 8, 0)
		}
		var upLen = float64(uy-ly) * float64(ellipseYOffset)
		var tilt = math.Atan2(float64(s.Boundary[0].X)-float64(mx)*0.25, float64(s.Boundary[0].Y)-float64(ly)+float64(uy-ly)*0.5) + math.Pi*(135.0/180.0)
		var pt3 = cv.Point{
			X: int(float64(mx)*0.25 - math.Cos(tilt)*upLen),
			Y: int(float64(ly) + float64(uy-ly)*0.5 - math.Sin(tilt)*upLen),
		}
		fmt.Printf("Tilt: %v", math.Atan2(float64(s.Boundary[0].X)-float64(pt3.X), float64(s.Boundary[0].Y)-float64(ly)+float64(uy-ly)*0.5)/math.Pi*180)
		cv.Circle(camImg, pt3, 7, cv.NewScalar(255, 100, 0, 0), 4, 8, 0)
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

	return qrtext, nil
}
