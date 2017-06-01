package main

import (
	"fmt"
	"runtime"
	"time"

	cv "github.com/lazywei/go-opencv/opencv"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/opencv"
	"gobot.io/x/gobot/platforms/parrot/ardrone"

	"github.com/Glorforidor/Chaos-Drone/oor"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//_, currentfile, _, _ := runtime.Caller(0)
	//cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")
	window := opencv.NewWindowDriver()
	camera := opencv.NewCameraDriver("tcp://192.168.1.1:5555")
	ardroneAdaptor := ardrone.NewAdaptor("192.168.1.1")
	drone := ardrone.NewDriver(ardroneAdaptor)

	goOOR := oor.New()

	killThisProgram := false // Turn on to make the drone land

	defer (func() {
		drone.Land()
		ardroneAdaptor.Finalize()
		window.Connection().Finalize()
		camera.Connection().Finalize()
	})()

	work := func() {
		detect := false
		if killThisProgram {
			drone.Land()
			ardroneAdaptor.Finalize()
			window.Connection().Finalize()
			camera.Connection().Finalize()
		} else {
			drone.TakeOff()
		}
		var image *cv.IplImage
		camera.On(opencv.Frame, func(data interface{}) {
			image = data.(*cv.IplImage)
			if !detect {
				window.ShowImage(image)
			}
		})
		hover := false
		drone.On(ardrone.Flying, func(data interface{}) {
			gobot.After(1*time.Second, func() { drone.Up(0.2) })
			gobot.After(2*time.Second, func() { drone.Hover() })
			gobot.After(5*time.Second, func() {
				detect = true
				gobot.Every(300*time.Millisecond, func() {
					if hover {
						if killThisProgram {
							drone.Land()
							ardroneAdaptor.Finalize()
							window.Connection().Finalize()
							camera.Connection().Finalize()
						} else {
							drone.Hover()
						}
					}
				})
				gobot.Every(1000*time.Millisecond, func() {
					i := image

					fmt.Println(i)

					if i != nil {
						ellipseData, err := goOOR.DetectEllipses(i.GetMat())
						if err != nil {
							var lx, ly, ux, uy int
							lx = ellipseData[0] // Rectangle left
							ly = ellipseData[1] // Rectangle top
							ux = ellipseData[2] // Rectangle right
							uy = ellipseData[3] // Rectangle bottom
							//cx = ellipseData[4] // Image center X
							//cy = ellipseData[5] // Image center Y

							lt := cv.Point{}
							br := cv.Point{}
							lt.X = lx
							lt.Y = ly
							br.X = ux
							br.Y = uy

							rect := cv.Rect{}
							rect.Init(lx, ly, ux-lx, uy-ly)

							fmt.Printf("Rectangle: (lx = %d, ly = %d), (ux = %d, ly = %d), w = %d, h = %d\n", lx, ly, ux, uy, ux-lx, uy-uy)

							opencv.DrawRectangles(i, []*cv.Rect{&rect}, 0, 255, 0, 5)

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
						}

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
					}
					window.ShowImage(i)
				})
				gobot.After(10*time.Second, func() {
					hover = false
					drone.Land()
					ardroneAdaptor.Finalize()
					window.Connection().Finalize()
					camera.Connection().Finalize()
				})
			})
		})
	}

	robot := gobot.NewRobot("face",
		[]gobot.Connection{ardroneAdaptor},
		[]gobot.Device{window, camera, drone},
		work,
	)

	robot.Start()
}
