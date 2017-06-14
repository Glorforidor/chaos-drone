package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	cv "github.com/lazywei/go-opencv/opencv"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/audio"
	"gobot.io/x/gobot/platforms/opencv"
	"gobot.io/x/gobot/platforms/parrot/ardrone"

	"github.com/Glorforidor/chaos-drone/barcode"
	"github.com/Glorforidor/chaos-drone/navigation"
	"github.com/Glorforidor/chaos-drone/oor"
)

const (
	moveSpeed   = 0.025
	rotateSpeed = 0.005
	detectDelay = 2
)

var (
	barcodeIndex = 0
	barcodes     = []string{"P.04"}
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	barcode.Init()

	window := opencv.NewWindowDriver()
	camera := opencv.NewCameraDriver("tcp://192.168.1.1:5555")
	ardroneAdaptor := ardrone.NewAdaptor("192.168.1.1") // ardrone2_117047
	drone := ardrone.NewDriver(ardroneAdaptor)
	audioDriver := audio.NewAdaptor()
	audioControl := make([]audio.Driver, 3)
	audioControl[0] = *audio.NewDriver(audioDriver, "./navigation/audio1.mp3")
	audioControl[1] = *audio.NewDriver(audioDriver, "./navigation/audio2.mp3")
	audioControl[2] = *audio.NewDriver(audioDriver, "./navigation/audio3.mp3")

	goOOR := oor.New()
	defer goOOR.Free()

	// killThisProgram AwesomeAs' killing machine!!!... will stop the program;)
	killThisProgram := false // Turn on to make the drone land
	onlyCameraFeed := false  // Turn on to prevent flying, so we can collect data.

	//ringBuffer := [4]cv.Rect{}

	// appendToRingBuffer appends a ring image to the ringBuffer. Every append
	// shuffle the array left removing the first.
	/*appendToRingBuffer := func(bounds cv.Rect) {
		for i := 1; i < 4; i++ {
			ringBuffer[i-1] = ringBuffer[i]
		}
		ringBuffer[3] = bounds
	}*/

	/*getMedianOfRingBuffers := func() [4]int {
		sum := [4]int{}
		for i := 0; i < 4; i++ {
			if ringBuffer[i].Width() > sum[2] && ringBuffer[i].Height() > sum[3] {
				sum[0] = ringBuffer[i].X()
				sum[1] = ringBuffer[i].Y()
				sum[2] = ringBuffer[i].Width()
				sum[3] = ringBuffer[i].Height()
			}
		}
		/ *for i := 0; i < 4; i++ {
			sum[i] = int(float64(sum[i]) / 5.0)
		}* /
		return sum
	}*/

	if killThisProgram {
		fmt.Println("KILLTHISPROGRAM IS ACTIVE! SHUTTING DOWN DRONE!")
	}

	// kill everything after main returns.
	defer func() {
		drone.Land()
		ardroneAdaptor.Finalize()
		audioDriver.Finalize()
		camera.Connection().Finalize()
	}()

	work := func() {
		// again kill if necessary.
		if killThisProgram {
			drone.Land()
			ardroneAdaptor.Finalize()
			audioDriver.Finalize()
			camera.Connection().Finalize()
			return
		} else if !onlyCameraFeed {
			drone.TakeOff()
		}

		// image is the drone image which will be used for detecting rings and
		// barcodes.
		var image *cv.IplImage

		// detect determines when a image is detected by other functions.
		detect := false

		// turn on the camera driver.
		camera.On(opencv.Frame, func(data interface{}) {
			// Type assert the raw camera data to opencv image format
			image = data.(*cv.IplImage)

			// If not detected by other functions let it show here. This is just
			// to give feedback immediately
			if !detect {
				window.ShowImage(image)
			}
		})
		flyingFunc := func(data interface{}) {
			if !onlyCameraFeed && !killThisProgram {
				gobot.After(1*time.Second, func() {
					drone.Up(0.05)
					time.Sleep(2500 * time.Millisecond)
					drone.Hover()
				})
				gobot.After(detectDelay*100.0*time.Millisecond, func() {
					drone.Hover()
					errs := audioControl[rand.Intn(3)].Play()
					for _, err := range errs {
						fmt.Printf("An error occoured with audio: %v\n", err)
					}
				})
			} else {
				drone.Land()
			}
			if onlyCameraFeed {
				gobot.After(detectDelay*100.0*time.Millisecond, func() {
					errs := audioControl[rand.Intn(3)].Play()
					for _, err := range errs {
						fmt.Printf("An error occoured with audio: %v\n", err)
					}
				})
			}
			gobot.After(detectDelay*100.0*time.Millisecond, func() {
				log.Println("Detect enabled.")
				detect = true

				// perhaps kill the program one more time?
				if !onlyCameraFeed {
					gobot.Every(300*time.Millisecond, func() {
						if killThisProgram {
							drone.Land()
							ardroneAdaptor.Finalize()
							audioDriver.Finalize()
							camera.Connection().Finalize()
						}
					})
				}

				// qrPoint holds the x and y coordinates of the position over
				// the barcode. These coordinates should be the place where the
				// drone fly through.
				var qrPoints []cv.Point
				// qrPointSet is to determine if a set has been found.
				var qrPointSet bool
				gobot.Every(300*time.Millisecond, func() {
					if image == nil {
						log.Printf("image not captured: %v\n", image)
						return
					}
					img := image.Clone()
					defer img.Release()
					qrPointSet = false
					qrPointsTmp, err := barcode.GetEllipseOverQR(img, "P.05")
					//qrText, qrErr := barcode.QRScan(i2)
					if err != nil {
						fmt.Printf("An error occoured with QR scanning: %v\n", err)
					} else if qrPointsTmp != nil {
						fmt.Printf("QR code found, position: %d, %d\n", qrPointsTmp[0].X, qrPointsTmp[0].Y)
						//navigation.IsLocked = true
						// drone.Up(0.01)
						//time.Sleep(500 * time.Millisecond)
						// drone.Forward(0.05)
						//time.Sleep(2 * time.Second)
						// drone.Hover()
						//navigation.IsLocked = false
						qrPoints = qrPointsTmp
						qrPointSet = true
						//cv.Circle(i2, cv.Point{X: ellipsePoint[0], Y: ellipsePoint[1]}, 8, cv.NewScalar(0, 0, 255, 0), 4, 8, 0)
					} else {
						fmt.Println("No QR codes detected.")
					}
				})

				gobot.Every(200*time.Millisecond, func() {
					if image == nil {
						log.Printf("image not captured: %v\n", image)
						return
					}
					// clone image so we don't work directly on the stream.
					img := image.Clone()
					defer img.Release()

					if qrPointSet {
						// draw red circle where the drone should fly through.
						cv.Circle(img, qrPoints[0], 8, cv.NewScalar(0, 255, 0, 0), 4, 8, 0)
						cv.Circle(img, qrPoints[1], 8, cv.NewScalar(0, 0, 255, 0), 4, 8, 0)
					}

					// appendToRingBuffer(rect)
					// scan the image for ellipse and get information where it
					// is.
					ellipseData, err := goOOR.DetectEllipses(img.GetMat())
					if err != nil {
						log.Printf("could not detect ellipse: %v\n", err)
						return
					}

					x := ellipseData[0]  // Rectangle left
					y := ellipseData[1]  // Rectangle top
					w := ellipseData[2]  // Rectangle right
					h := ellipseData[3]  // Rectangle bottom
					cx := ellipseData[4] // Image center X
					cy := ellipseData[5] // Image center Y

					// medBounds := getMedianOfRingBuffers()
					// x = medBounds[0]
					// y = medBounds[1]
					// w = medBounds[2]
					// h = medBounds[3]

					center := navigation.Center(x, y, w, h)
					log.Println("Center of the ring:", center)
					dp := navigation.Placement(cx, cy, center.X, center.Y)
					log.Println("Drones placement of the ring:", dp)
					var cp int
					if qrPointSet {
						cp = navigation.Placement(cx, cy, qrPoints[0].X, qrPoints[0].Y)
						log.Println("QR placement of the ring:", cp)
					}

					// construct a rectangle from the ellipse data.
					rect := cv.NewRect(x, y, w, h)

					//appendToRingBuffer(rect)

					//fmt.Printf("Rectangle: (x = %d, y = %d), w = %d, h = %d\n", x, y, w, h)

					// draw the rectangle on the image.
					opencv.DrawRectangles(img, []*cv.Rect{&rect}, 0, 255, 0, 5)

					// show the image on screen.
					window.ShowImage(img)

					if !onlyCameraFeed {
						if qrPointSet {
							if onTarget := navigation.Move(drone, cp); onTarget {
								navigation.FlyThroughRing(drone, 200, qrPoints[0].X-cx)
							}
						} else {
							if onTarget := navigation.Move(drone, dp); onTarget {
								navigation.FlyThroughRing(drone, 200, center.X-cx)
							}
						}
					}

					// barcodeIndex++
					// if barcodeIndex > 6 {
					// drone.Land()
					// }
				})
				if onlyCameraFeed {
					gobot.After(50*time.Second, func() {
						drone.Land()
						ardroneAdaptor.Finalize()
						audioDriver.Finalize()
						camera.Connection().Finalize()
					})
				} else {
					gobot.After(20*time.Second, func() {
						drone.Land()
						ardroneAdaptor.Finalize()
						audioDriver.Finalize()
						camera.Connection().Finalize()
					})
				}
			})
		}

		if onlyCameraFeed {
			flyingFunc(nil)
		} else {
			if err := drone.On(ardrone.Flying, flyingFunc); err != nil {
				log.Printf("the flying failed: %v\n", err)
				return
			}
		}
	}

	robot := gobot.NewRobot("Ardrone",
		[]gobot.Connection{ardroneAdaptor, audioDriver},
		[]gobot.Device{window, camera, drone, &audioControl[0], &audioControl[1], &audioControl[2]},
		work,
	)

	robot.Start()
}

/*opencv.DrawRectangles(
i,
[]*cv.Rect{cv.Rect(
	i,
	lt,
	br,
	cv.NewScalar(0, 0, 0, 0),
	1,
	1,
	0)},
0, 255, 0, 5)*/

/*faces := opencv.DetectFaces(cascade, i)
biggest := 0
var face *cv.Rect
for _, f := range faces {
	if f.Width() > biggest {
		biggest = f.Width()
		face = f
	}
}
if face != nil {
	opencv.DrawRectangles(i, []*cv.Rect{face}, 0, 255, 0, 5)
	hystX := 0.2
	hystY := 0.2
	centerX := float64(image.Width()) * 0.5
	centerY := float64(image.Height()) * 0.5
	turnX := -(float64(face.X()) - centerX) / centerX
	turnY := -(float64(face.Y()) - centerY) / centerY
	//Find object horizontal
	if turnX < -hystX {
		fmt.Println("turning ClockWise:", turnX)
		//drone.Clockwise(math.Abs((turn + 0.2) * 0.2))
		drone.Clockwise(0.01)
	} else if turnX > hystX {
		fmt.Println("turning CounterClockWise:", turnX)
		//drone.CounterClockwise(math.Abs((turn - 0.2) * 0.2))
		drone.CounterClockwise(0.01)
	} else if turnY > hystY { //Find object vertital
		fmt.Println("turning Up:", turnY)
		//drone.Clockwise(math.Abs((turn + 0.2) * 0.2))
		drone.Up(0.01)
	} else if turnY < -hystY {
		fmt.Println("turning Down:", turnY)
		//drone.CounterClockwise(math.Abs((turn - 0.2) * 0.2))
		drone.Down(0.01)
	} else { //if turnX < hystX && turnX > -hystX && turnY < hystY && turnY > -hystY {

		fmt.Println("Forward:")
		hover = false
		drone.Forward(0.1)
		gobot.After(1*time.Second, func() { hover = true })
	}
}*/
