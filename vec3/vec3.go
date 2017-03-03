package vec3

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

//3D vector struct.
type Vec3 struct{ X, Y, Z float64 }

//Creates a new 3D point given either none or 3 numbers.
//Default constructor is (0, 0, 0).
func New(coordinates ...float64) *Vec3 {
	if len(coordinates) >= 3 {
		return &Vec3{coordinates[0], coordinates[1], coordinates[2]}
	} else {
		return &Vec3{0, 0, 0}
	}
}

//Parses a string to find x, y and z values.
//Returns an error if the string is invalid.
func FromString(str string) (*Vec3, error) {
	re := regexp.MustCompile("\\-?\\d+(\\.\\d+)?")
	coords := re.FindAllString(str, 3)
	if len(coords) == 3 {
		x, _ := strconv.ParseFloat(coords[0], 64)
		y, _ := strconv.ParseFloat(coords[1], 64)
		z, _ := strconv.ParseFloat(coords[2], 64)
		return &Vec3{x, y, z}, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid input string '%v'", str))
	}
}

//Beautify string representation.
func (v *Vec3) String() string {
	return fmt.Sprintf("[%f, %f, %f]", v.X, v.Y, v.Z)
}

//The distance of the vector.
func (v *Vec3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

//Vector1 + Vector2
func (v *Vec3) Add(o *Vec3) *Vec3 {
	return &Vec3{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

//Vector1 - Vector2
func (v *Vec3) Substract(o *Vec3) *Vec3 {
	return &Vec3{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

//Vector1 * Vector2, Vector1 * number.
func (v *Vec3) Multiply(o interface{}) *Vec3 {
	switch vv := o.(type) {
	case Vec3:
		return &Vec3{v.X * vv.X, v.Y * vv.Y, v.Z * vv.Z}
	case float64:
		return &Vec3{v.X * vv, v.Y * vv, v.Z * vv}
	case int:
		vv2 := float64(vv)
		return &Vec3{v.X * vv2, v.Y * vv2, v.Z * vv2}
	default:
		return v
	}
}

//Vector1 / Vector2, Vector1 / number.
func (v *Vec3) Divide(o interface{}) *Vec3 {
	switch vv := o.(type) {
	case Vec3:
		return &Vec3{v.X / vv.X, v.Y / vv.Y, v.Z / vv.Z}
	case float64:
		return &Vec3{v.X / vv, v.Y / vv, v.Z / vv}
	case int:
		vv2 := float64(vv)
		return &Vec3{v.X / vv2, v.Y / vv2, v.Z / vv2}
	default:
		return v
	}
}

//Dot product of two vectors.
//The dot product is directly related to the cosine of the angle
//between two vectors in Euclidean space of any number of dimensions.
func (v *Vec3) Dot(o *Vec3) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

//Cross product of two vectors.
//The cross product a × b is defined as a vector c that is perpendicular
//to both a and b, with a direction given by the right-hand rule and a
//magnitude equal to the area of the parallelogram that the vectors span.
func (v *Vec3) Cross(o *Vec3) *Vec3 {
	return &Vec3{v.Y*o.Z - v.Z*o.Y, v.Z*o.X - v.X*o.Z, v.X*o.Y - v.Y*o.X}
}
