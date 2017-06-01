package vec3

import (
	"fmt"
	"math"
)

// Vec3 is a 3d vector struct.
type Vec3 struct{ X, Y, Z float64 }

// New creates a new 3D point given either none or 3 numbers.
// Default constructor is (0, 0, 0).
func New(coordinates ...float64) *Vec3 {
	if len(coordinates) >= 3 {
		return &Vec3{coordinates[0], coordinates[1], coordinates[2]}
	}

	return nil, fmt.Errorf("Invalid input string '%v'", str)
}

//Beautify string representation.
func (v *Vec3) String() string {
	return fmt.Sprintf("[%f, %f, %f]", v.X, v.Y, v.Z)
}

// Magnitude is the distance of the vector.
func (v *Vec3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Add adds to vectors.
func (v *Vec3) Add(o *Vec3) *Vec3 {
	return &Vec3{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

// Substract subtracts to vectors.
func (v *Vec3) Substract(o *Vec3) *Vec3 {
	return &Vec3{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

// Multiply multiplies to vectors.
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

// Divide divides to vectors.
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

// Dot product of two vectors.
// The dot product is directly related to the cosine of the angle
// between two vectors in Euclidean space of any number of dimensions.
func (v *Vec3) Dot(o *Vec3) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

// Cross product of two vectors.
// The cross product a × b is defined as a vector c that is perpendicular
// to both a and b, with a direction given by the right-hand rule and a
// magnitude equal to the area of the parallelogram that the vectors span.
func (v *Vec3) Cross(o *Vec3) *Vec3 {
	return &Vec3{v.Y*o.Z - v.Z*o.Y, v.Z*o.X - v.X*o.Z, v.X*o.Y - v.Y*o.X}
}
