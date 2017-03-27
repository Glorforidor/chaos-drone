package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"

	cv "github.com/lazywei/go-opencv/opencv"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/opencv"
	"gobot.io/x/gobot/platforms/parrot/ardrone"

	"gopkg.in/bieber/barcode.v0"
)

var font = cv.InitFont(cv.CV_FONT_HERSHEY_SIMPLEX, float32(0.65), float32(0.65), float32(1.0), 2, 8)
var fontColor = cv.NewScalar(0, 0, 255, 0)
var fontColorWhite = cv.NewScalar(255, 255, 255, 0)

func nextIndex(i int, max int) int {
	if i < max {
		return i + 1
	}
	return 0
}

func qrScan(camImg *cv.IplImage) {
	img := barcode.NewImage(camImg.ToImage())
	scanner := barcode.NewScanner().SetEnabledAll(true)

	symbols, _ := scanner.ScanImage(img)
	for _, s := range symbols {
		fmt.Println(s.Type.Name(), s.Data, s.Quality, s.Boundary)
		for i := 0; i < 4; i++ {
			var pt1 = cv.Point{X: s.Boundary[i].X, Y: s.Boundary[i].Y}
			var pt2 = cv.Point{X: s.Boundary[nextIndex(i, 3)].X, Y: s.Boundary[nextIndex(i, 3)].Y}
			cv.Line(camImg, pt1, pt2, cv.NewScalar(0, 255, 0, 0), 2, 8, 0)
		}
		for x := -1; x < 2; x++ {
			for y := -6; y < -3; y++ {
				font.PutText(camImg, s.Data, cv.Point{X: s.Boundary[1].X + 4 + x, Y: s.Boundary[1].Y + y}, fontColorWhite)
			}
		}
		font.PutText(camImg, s.Data, cv.Point{X: s.Boundary[1].X, Y: s.Boundary[1].Y - 5}, fontColor)
	}
}

func main() {
	camera := opencv.NewCameraDriver("tcp://192.168.1.1:5555")
	ardroneAdaptor := ardrone.NewAdaptor()
	drone := ardrone.NewDriver(ardroneAdaptor)

	win := opencv.NewWindowDriver()
	var image *cv.IplImage
	detect := false

	work := func() {
		camera.On(opencv.Frame, func(data interface{}) {
			image = data.(*cv.IplImage)
			if !detect {
				win.ShowImage(image)
			}
		})

		gobot.After(2*time.Second, func() {
			detect = true
			for i := 0; i < 50; i++ {
				qrScan(image)
				win.ShowImage(image)
				time.Sleep(300 * time.Millisecond)
			}
			os.Exit(0)
		})
	}

	robot := gobot.NewRobot("qrcode",
		[]gobot.Connection{ardroneAdaptor},
		[]gobot.Device{win, camera, drone},
		work,
	)

	robot.Start()
}
