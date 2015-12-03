package geometry

import (
	"errors"
	"math"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
)

// Vector : Any sized dimension representation of a point of vector space.
type Vector []float64

// VectorEquality – Trying to implement a quick version of checking for Vector equality
type VectorEquality interface {
	Equal(VectorEquality) bool
}

// Equal method implements an Equality comparison between vectors.
func (v Vector) Equal(u VectorEquality) bool {
	if len(v) != len(u.(Vector)) {
		return false
	}
	for i := 0; i < len(v); i++ {
		if v[i] != u.(Vector)[i] {
			return false
		}
	}
	return true
}

// just providing a conventional element naming system for
// 2D and 3D vectors.
// this 'aliasing' is just local to this file.
const (
	x = iota
	y
	z
)

// VecAddition – performs vector addition between two vectors, v and u
// s.t. v + u = [v₁+u₁ , v₂+u₂, ⠂⠂⠂ , vᵢ₋₁+uᵢ₋₁ , vᵢ+uᵢ ]
// on an i-th dimension vector.
func VecAddition(v Vector, u Vector) (Vector, error) {
	if len(v) != len(u) {
		return nil, errors.New("vector dimensions do not coincide")
	}
	var vPlusU Vector
	for i := 0; i < len(v); i++ {
		vPlusU = append(vPlusU, (v[i] + u[i])) //	add an element to the new vector which is the sum of element i from v and u
	}
	return vPlusU, nil
}

// VecScalarMultiply - performs scalar multiplication on a vector v,
// s.t. scalar * v = [scalar*e1, scalar*e2, scalar*e3]
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

// DotProduct returns the sum of the product of elements of
// two i-dimension vectors, u and v, as a scalar
// s.t. v•v = (v₁u₁ + v₂u₂ + ⠂⠂⠂ + vᵢ₋₁uᵢ₋₁ + vᵢuᵢ)
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
func AngleFromOrigin(v Vector) (float64, error) {
	if len(v) != 2 {
		return 0, errors.New("vector dimension != 2")
	}
	return calc.ToFixed(math.Atan2(v[x], v[y]), 5), nil
}

// RelativeAngle – does what it says on the box.
// Only implemented in 2D currently, or, as a comparitive rotation
// on the z-axis in 3D.
// (Not restricted, however, as one could have higher dimensionalites
// of vector and still test for relative angle, by passing a slice which indicates which axis of rotation to be performed. TO BE PROPERLY IMPLEMETED LATER!)
func RelativeAngle(v Vector, u Vector) (float64, error) {
	if len(v) == 0 || len(u) == 0 {
		return 0, errors.New("v or u is an empty vector")
	}
	if len(v) != len(u) {
		return 0, errors.New("vector dimensions do not coincide")
	}
	det := (v[x] * u[y]) - (v[y] * u[x])
	dot, err := DotProduct(v, u)
	if err != nil {
		return 0, err
	}
	angle := math.Atan2(det, dot)
	return calc.ToFixed(angle, 5), nil
}

// UnitAngle will map any floating-point value to its angle on a unit circle.
func UnitAngle(angle float64) float64 {
	twoPi := math.Pi * 2
	return calc.ToFixed((angle - (twoPi * math.Floor(angle/twoPi))), 5)
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
	return calc.ToFixed(math.Sqrt(ǁvǁsq), 5), nil
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

// Normalise returns the normalised Vector of v
// – it's what you might call a direction vector,
// as opposed to a position vector.
// NOTE: this is not the same as a 'norm'/'normal' which is the vector
// orthogonal to a plane or surface.
func Normalise(v Vector) (Vector, error) {
	if len(v) == 0 {
		return nil, errors.New("v is an empty vector")
	}
	var norm Vector
	var ǁvǁ float64
	var err error
	for i := 0; i < len(v); i++ {
		ǁvǁ, err = Magnitude(v)
		if err != nil {
			return nil, err
		}
		norm = append(norm, calc.ToFixed(v[i]/ǁvǁ, 5))
	}
	return norm, nil
}

// TranslatePositionToSector2D : translates the co-ordinates of a 2D vector to sector indices location (2D Version)
func TranslatePositionToSector2D(ed float64, n int, v Vector) (int, int) {
	fn := float64(n)
	col := int(((v[x] + ed) / (2 * ed)) * fn)
	row := int(((-1 * (v[y] - ed)) / (2 * ed)) * fn)
	return row, col
}
