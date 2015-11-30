package main

import (
	"errors"
	"math"
)

// VecAddition v1 + v2 = v3
func VecAddition(v Vector, u Vector) (Vector, error) {
	if len(v) != len(u) {
		return nil, errors.New("vector dimensions do not coincide")
	}
	var vPlusU Vector
	for i := 0; i < len(v); i++ {
		vPlusU = append(vPlusU, (v[i] + u[i])) //	add an element to the new vector which is the sum of element i from v1 and v2
	}
	return vPlusU, nil
}

// VecScalarMultiply scalar * v = [scalar*e1, scalar*e2, scalar*e3]
func VecScalarMultiply(v Vector, scalar float64) (Vector, error) {
	if len(v) == 0 {
		return nil, errors.New("v is an empty vector")
	}
	var sm Vector
	for i := range v {
		sm = append(sm, (v[i] * scalar))
	}
	return sm, nil
}

// DotProduct returns the sum of the product of elements of two i-dimension vectors:
// v₁u₁ + v₂u₂ + ⠂⠂⠂ + vᵢ₋₁+ vᵢ
func DotProduct(v Vector, u Vector) (float64, error) {
	if len(v) != len(u) {
		return 0, errors.New("vector dimensions do not coincide")
	}
	var f float64
	for i := 0; i < len(v); i++ {
		f += (v[i] * u[i])
	}
	return f, nil
}

// CrossProduct produces a new Vector orthogonal to both v and u.
// Only supported for 3D vectors.
func CrossProduct(v Vector, u Vector) (Vector, error) {
	if len(v) != 3 || len(u) != 3 {
		return nil, errors.New("vector dimension != 3")
	}
	var cp Vector
	cp = append(cp, (v[y]*u[z] - v[z]*u[y]))
	cp = append(cp, (v[z]*u[x] - v[x]*u[z]))
	cp = append(cp, (v[x]*u[y] - v[y]*u[x]))
	return cp, nil
}

// AngleFromOrigin calculates the angle of a given vector from the origin
// relative to the x-axis of 𝐄 (the model environment)
func AngleFromOrigin(v Vector) float64 {
	return math.Atan2(v[x], v[y])
}

// RelativeAngle – does what it says on the box.
func RelativeAngle(v Vector, u Vector) (float64, error) {
	if len(v) == 0 || len(u) == 0 {
		return 0, errors.New("v or u is an empty vector")
	}
	det := (v[x] * u[y]) - (v[y] * u[x])
	dot, err := DotProduct(v, u)
	if err != nil {
		return 0, err
	}
	angle := math.Atan2(det, dot)
	return angle, nil
}

// UnitAngle will map any floating-point value to its angle on a unit circle.
func UnitAngle(angle float64) float64 {
	twoPi := math.Pi * 2
	return angle - (twoPi * math.Floor(angle/twoPi))
}

// TranslatePositionToSector2D : translates the co-ordinates of a 2D vector to sector indices location (2D Version)
func TranslatePositionToSector2D(ed float64, n int, v Vector) (int, int) {
	fn := float64(n)
	x := int((v[x] + ed) / (2 * ed) * fn)
	y := int((v[y] + ed) / (2 * ed) * fn)
	return x, y
}

// Magnitude does the classic calculation for length of a vector
// (or, distance from origin)
func Magnitude(v Vector) (float64, error) {
	if len(v) == 0 {
		return 0, errors.New("v is an empty vector")
	}
	var ǁvǁsq float64
	for i := 0; i < len(v); i++ {
		ǁvǁsq += v[i] * v[i]
	}
	return math.Sqrt(ǁvǁsq), nil
}

// VectorDistance calculates the distance between two positions
func VectorDistance(v Vector, u Vector) (float64, error) {
	if len(v) != len(u) {
		return 0, errors.New("vector dimensions do not coincide")
	}
	vd := Vector{}
	for i := 0; i < len(v); i++ {
		diff := (v[i] - u[i])
		vd = append(vd, diff)
	}
	return Magnitude(vd)
}

/*ColourDistance quantifies the value difference between two ColRGB structs,
returning a floating-point ratio from 0.0 to 1.0.
Multiply the returned value by100 for a percentage.
NOTE: this is a distinct concept from the distance between them as 3D vectors,
as there would be 2 other ColRGB for any ColRGB with an identical magnitude.
e.g. [255 0 0] [0 255 0] [0 0 255] will all have the same magnitude, but are
pure Red, pure Blue, pure Green respectively! */
func ColourDistance(c1 ColRGB, c2 ColRGB) float64 {
	redDiff := float64(c1.red-c2.red) / 255
	greenDiff := float64(c1.green-c2.green) / 255
	blueDiff := float64(c1.blue-c1.blue) / 255
	return (redDiff + greenDiff + blueDiff) / 3
	/*	will  */
}
