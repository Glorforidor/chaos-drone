package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/parrot/ardrone"
)

func main() {
	ardroneAdapter := ardrone.NewAdaptor("192.168.1.1")
	drone := ardrone.NewDriver(ardroneAdapter)

	work := func() {
		drone.On(drone.Event(ardrone.Flying), func(data interface{}) {
			gobot.After(2*time.Second, func() {
				drone.Land()
			})
		})
		drone.TakeOff()
	}

	robot := gobot.NewRobot("drone",
		[]gobot.Connection{ardroneAdapter},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
