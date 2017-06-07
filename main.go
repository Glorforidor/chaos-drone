package main

import (
	"fmt"
	"runtime"
	"time"

	cv "github.com/lazywei/go-opencv/opencv"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/opencv"
	"gobot.io/x/gobot/platforms/parrot/ardrone"

	"github.com/Glorforidor/chaos-drone/barcode"
	"github.com/Glorforidor/chaos-drone/navigation"
	"github.com/Glorforidor/chaos-drone/oor"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//_, currentfile, _, _ := runtime.Caller(0)
	//cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")
	window := opencv.NewWindowDriver()
	camera := opencv.NewCameraDriver("tcp://192.168.1.1:5555")
	ardroneAdaptor := ardrone.NewAdaptor("192.168.1.1") // ardrone2_117047
	drone := ardrone.NewDriver(ardroneAdaptor)

	goOOR := oor.New()
	defer goOOR.Free()

	killThisProgram := false // Turn on to make the drone land
	onlyCameraFeed := false  // Turn on to prevent flying, so we can collect data.

	const moveSpeed = 0.025
	const rotateSpeed = 0.005
	const detectDelay = 5

	ringBuffer := [4]cv.Rect{}

	appendToRingBuffer := func(bounds cv.Rect) {
		for i := 1; i < 4; i++ {
			ringBuffer[i-1] = ringBuffer[i]
		}
		ringBuffer[3] = bounds
	}

	getMedianOfRingBuffers := func() [4]int {
		sum := [4]int{}
		for i := 0; i < 4; i++ {
			if ringBuffer[i].Width() > sum[2] && ringBuffer[i].Height() > sum[3] {
				sum[0] = ringBuffer[i].X()
				sum[1] = ringBuffer[i].Y()
				sum[2] = ringBuffer[i].Width()
				sum[3] = ringBuffer[i].Height()
			}
		}
		/*for i := 0; i < 4; i++ {
			sum[i] = int(float64(sum[i]) / 5.0)
		}*/
		return sum
	}

	if killThisProgram {
		fmt.Println("KILLTHISPROGRAM IS ACTIVE! SHUTTING DOWN DRONE!")
	}

	defer (func() {
		drone.Land()
		ardroneAdaptor.Finalize()
		camera.Connection().Finalize()
	})()

	barcode.Init()

	work := func() {
		detect := false
		if killThisProgram {
			drone.Land()
			ardroneAdaptor.Finalize()
			camera.Connection().Finalize()
		} else if !onlyCameraFeed {
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
		flyingFunc := func(data interface{}) {
			if !onlyCameraFeed && !killThisProgram {
				gobot.After(1*time.Second, func() { drone.Up(0.9) })
				gobot.After(detectDelay*time.Second, func() { drone.Hover() /*navigation.FlyThroughRing(drone)*/ })
			} else {
				drone.Land()
			}
			gobot.After(detectDelay*time.Second, func() {
				detect = true
				if !onlyCameraFeed {
					gobot.Every(300*time.Millisecond, func() {
						if hover {
							if killThisProgram {
								drone.Land()
								ardroneAdaptor.Finalize()
								camera.Connection().Finalize()
							}
						}
					})
				}
				var qrPoint cv.Point
				var qrPointSet bool
				gobot.Every(300*time.Millisecond, func() {
					qrPointSet = false
					if image != nil {
						ellipsePoint, err := barcode.GetEllipseOverQR(image, "W02.00")
						//qrText, qrErr := barcode.QRScan(i2)
						if err != nil {
							fmt.Printf("An error occoured with QR scanning: %v\n", err)
						} else if ellipsePoint != nil {
							//fmt.Printf("Amount of QR codes: %d, Data: %v\n", len(qrText), qrText)
							qrPoint = cv.Point{X: ellipsePoint[0], Y: ellipsePoint[1]}
							qrPointSet = true
							//cv.Circle(i2, cv.Point{X: ellipsePoint[0], Y: ellipsePoint[1]}, 8, cv.NewScalar(0, 0, 255, 0), 4, 8, 0)
						}
					}
				})

				gobot.Every(300*time.Millisecond, func() {
					i := image

					var i2 *cv.IplImage

					if i != nil {
						i2 = i.Clone()

						if qrPointSet {
							cv.Circle(i2, qrPoint, 8, cv.NewScalar(0, 0, 255, 0), 4, 8, 0)
						}

						ellipseData, err := goOOR.DetectEllipses(i2.GetMat())
						if err == nil {
							var x, y, w, h, cx, cy int
							//x = ellipseData[0]  // Rectangle left
							//y = ellipseData[1]  // Rectangle top
							//w = ellipseData[2]  // Rectangle right
							//h = ellipseData[3]  // Rectangle bottom
							cx = ellipseData[4] // Image center X
							cy = ellipseData[5] // Image center Y

							medBounds := getMedianOfRingBuffers()
							x = medBounds[0]
							y = medBounds[1]
							w = medBounds[2]
							h = medBounds[3]

							if !navigation.IsLocked {

								ratio, err := navigation.Ratio(w, h)
								if !onlyCameraFeed && err == nil {
									dir, badness := navigation.Direction(ratio)
									center := navigation.Center(x, y, w, h)
									move := navigation.Placement(cx, cy, center.X, center.Y)
									/*if w > 60 && h > 100 && w*h < 160*160 {
										navigation.IsLocked = true
										drone.Up(moveSpeed * 0.5)
										time.Sleep(140 * time.Millisecond)
										drone.Forward(moveSpeed * 0.05)
										time.Sleep(1000 * time.Millisecond)
										drone.Hover()
										time.Sleep(1000 * time.Millisecond)
										navigation.IsLocked = false
									} else {
										drone.Hover()
									}*/
									if badness > 0 {
										switch dir {
										case navigation.Horizontal:
											switch move {
											case navigation.Left:
												drone.CounterClockwise(rotateSpeed * badness)
												fmt.Println("Flying counter clockwise")
											case navigation.Right:
												drone.Clockwise(rotateSpeed * badness)
												fmt.Println("Flying clockwise")
											}
										case navigation.Vertical:
											switch move {
											case navigation.Down:
												drone.Down(moveSpeed * badness)
												fmt.Println("Flying down 1")
											case navigation.Up:
												//drone.Up(moveSpeed * badness)
												fmt.Println("Flying up 1")
											}
										}
									} else {
										switch move {
										case navigation.Down:
											//drone.Down(moveSpeed * badness)
											fmt.Println("Flying down 2")
										case navigation.Up:
											//drone.Up(moveSpeed * badness)
											fmt.Println("Flying up 2")
										case navigation.Left:
											//drone.Left(moveSpeed * badness)
											fmt.Println("Flying left")
										case navigation.Right:
											//drone.Right(moveSpeed * badness)
											fmt.Println("Flying right")
										case navigation.OnTarget:
											// Lock on
											//drone.Hover()
											fmt.Println("HECK YEAH!")
										}
									}
								}
							}

							rect := cv.Rect{}
							rect.Init(ellipseData[0], ellipseData[1], ellipseData[2], ellipseData[3])

							appendToRingBuffer(rect)

							//fmt.Printf("Rectangle: (x = %d, y = %d), w = %d, h = %d\n", x, y, w, h)

							opencv.DrawRectangles(i2, []*cv.Rect{&rect}, 0, 255, 0, 5)

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
					} else {
						fmt.Println("Image is nil!")
					}
					window.ShowImage(i2)
				})
				if onlyCameraFeed {
					gobot.After(15*time.Second, func() {
						drone.Land()
						ardroneAdaptor.Finalize()
						camera.Connection().Finalize()
					})
				} else {
					gobot.After(10*time.Second, func() {
						hover = false
						drone.Land()
						ardroneAdaptor.Finalize()
						camera.Connection().Finalize()
					})
				}
			})
		}

		if onlyCameraFeed {
			flyingFunc(nil)
		} else {
			drone.On(ardrone.Flying, flyingFunc)
		}
	}

	robot := gobot.NewRobot("face",
		[]gobot.Connection{ardroneAdaptor},
		[]gobot.Device{window, camera, drone},
		work,
	)

	robot.Start()
}

//Pcall acts as a protected call, returning wether the call went through successfully, and its return value.
func Pcall(f func([]interface{}) []interface{}, params []interface{}) (success bool, result []interface{}) {
	defer func() {
		if r := recover(); r != nil {
			success = false
			result = make([]interface{}, 1)
			result[0] = r
			fmt.Printf("An error occoured in Pcall: %v\n", r)
		}
	}()
	return true, f(params)
}
