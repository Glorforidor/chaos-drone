package navigation

import (
	"errors"
	"log"
	"math"
	"time"

	"gobot.io/x/gobot/platforms/parrot/ardrone"
)

const (
	// Vertical movement.
	Vertical = iota
	// Horizontal movement.
	Horizontal
	// Stay means stay!
	Stay
)

const (
	// Up movement.
	Up = iota
	// Left movement.
	Left
	// Down movement.
	Down
	// Right movement.
	Right
	// OnTarget means On Target!
	OnTarget
)

// ErrDivByZero as the name imply a divide by zero error.
var ErrDivByZero = errors.New("divide by zero")

// IsLocked determines if the drone will react to commands
var IsLocked = false

// Point contains a X and Y coordinate.
type Point struct {
	X, Y int
}

// round rounds to nearest integer.
func round(x float64) int {
	if x < 0 {
		return int(math.Ceil(x - 0.5))
	}
	return int(math.Floor(x + 0.5))
}

// Center calculates the center rectangle.
func Center(x, y, width, height int) Point {
	x = round(float64((x + x + width)) * 0.5)
	y = round(float64((y + y + height)) * 0.5)
	return Point{x, y}
}

// Ratio return the ratio between x and y. If the value is above 1 then x is
// greater than y and below 0 y is greater than x.
// An error of type ErrDivByZero is returned if y is zero or less.
func Ratio(x, y int) (float64, error) {
	if y <= 0 {
		return 0.0, ErrDivByZero
	}
	return float64(x) / float64(y), nil
}

// Direction checks the ratio to detect the drones possition relative to the
// object.
func Direction(ratio float64) (direction int, badness float64) {
	move := 1 - ratio
	direction = Stay
	badness = math.Abs(math.Log10(ratio))
	if 1.0 < badness {
		badness = 1.0
	}
	if move < -0.05 {
		direction = Vertical
	} else if move > 0.05 {
		direction = Horizontal
	} else {
		badness = 0.0
	}
	return direction, badness
}

// Placement calculates the position of the drone relative to the center of an
// ring.
func Placement(x, y, rx, ry int) int {
	// calculate the distance from the object.
	c := math.Sqrt(math.Pow(float64(x-rx), 2) + math.Pow(float64(y-ry), 2))
	log.Println("Distance:", c)
	// relative close to the object, just say it is on target.
	if c < 60 {
		return OnTarget
	}

	placement := OnTarget
	switch {
	case x < rx:
		placement = Left
	case x > rx:
		placement = Right
	case y < ry:
		placement = Up
	case y > ry:
		placement = Down
	}
	return placement
}

const flyTime = 10

// FlyThroughRing is a command to make the drone lock on and fly straight for one second.
func FlyThroughRing(drone *ardrone.Driver, size int, xdiff int) {
	//if !IsLocked {
	log.Println("I am flying bitches!")
	IsLocked = true
	drone.Hover()
	time.Sleep(500 * time.Millisecond)
	if xdiff < -50 {
		drone.Right(0.05)
		time.Sleep(200 * time.Millisecond)
	}
	if xdiff > 50 {
		drone.Left(0.05)
		time.Sleep(200 * time.Millisecond)
	}
	drone.Up(0.125)
	time.Sleep(550 * time.Millisecond)
	drone.Hover()
	time.Sleep(2000 * time.Millisecond)
	drone.Forward(0.025)
	i := flyTime - math.Log10(float64(size))
	if i < 0.1 {
		i = 0.1
	}
	time.Sleep(time.Duration(i) * time.Second)
	drone.Hover()
	time.Sleep(200 * time.Millisecond)
	drone.Land()
	IsLocked = false
	log.Println("Im done yay!")
	//}
}

const (
	speed     = 0.025
	sleepTime = 50
)

// Move moves in the opposite direction of the drone placement.
func Move(drone *ardrone.Driver, placement int) bool {
	if !IsLocked {
		log.Println("Locking")
		IsLocked = true
		switch placement {
		case Up: // The drone is above the ring, fly down.
			log.Println("Going down")
			drone.Down(speed)
			time.Sleep(sleepTime * time.Millisecond)
			log.Println("Done down")
		case Down: // The drone is below the ring, fly up.
			log.Println("Going up")
			drone.Up(speed)
			time.Sleep(sleepTime * time.Millisecond)
			log.Println("Done up")
		case Left: // The drone is left of the ring, fly right.
			log.Println("Going right")
			drone.Right(speed)
			time.Sleep(sleepTime * time.Millisecond)
			log.Println("Done right")
		case Right: // The drone is right of the ring, fly left.
			log.Println("Going left")
			drone.Left(speed)
			time.Sleep(sleepTime * time.Millisecond)
			log.Println("Done left")
		case OnTarget: // The drone is in the center.
			log.Println("OnTarget")
			log.Println("Unlocking")
			IsLocked = false
			return true
		}
	}
	log.Println("Unlocking")
	IsLocked = false
	drone.Hover()
	return false
}
