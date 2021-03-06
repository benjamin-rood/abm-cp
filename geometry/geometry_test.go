// geometry is for all the calculations related to agents and their interaction in the environment
// but not their behaviours.
package geometry

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/benjamin-rood/abm-cp/calc"
)

func Test(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	uval := float64(rand.Int())
	vval := float64(rand.Int())
	wval := float64(rand.Int())
	xval := float64(rand.Int())
	yval := float64(rand.Int())
	zval := float64(rand.Int())
	randU := rand.Float64()
	randV := rand.Float64()
	randW := rand.Float64()
	randX := rand.Float64()
	randY := rand.Float64()
	randZ := rand.Float64()

	var vecAdditionTests = []struct {
		v, u, want Vector
	}{
		{Vector{1.0, 0.0}, Vector{0.0, 1.0}, Vector{1.0, 1.0}},
		{Vector{xval, yval}, Vector{wval, zval}, Vector{xval + wval, yval + zval}},
		{Vector{randX, randY}, Vector{randW, randZ},
			Vector{randX + randW, randY + randZ}},
	}

	for _, c := range vecAdditionTests {
		got, _ := VecAddition(c.v, c.u)
		if got.Equal(c.want) == false {
			t.Errorf("VecAddition(%v, %v) == %v, want %v\n", c.v, c.u, got, c.want)
		}
	}

	var vecScalarMultiplicationTests = []struct {
		v    Vector
		s    float64
		want Vector
	}{
		{Vector{1.0, 0.0}, -0.1, Vector{1.0 * -0.1, 0.0 * -0.1}},
		{Vector{xval, yval}, randX, Vector{xval * randX, yval * randX}},
		{Vector{randX, randY}, 3.3, Vector{randX * 3.3, randY * 3.3}},
	}

	for _, d := range vecScalarMultiplicationTests {
		got, _ := VecScalarMultiply(d.v, d.s)
		if got.Equal(d.want) == false {
			t.Errorf("VecScalarMultiply(%q, %v) == %q, want %q\n", d.v, d.s, got, d.want)
		}
	}

	var dotProductTests = []struct {
		v, u Vector
		want float64
	}{
		{Vector{1.0, 0.0}, Vector{0.0, 1.0}, 0.0},
		{Vector{xval, yval}, Vector{wval, zval}, (xval*wval + yval*zval)},
		{Vector{randX, randY}, Vector{randW, randZ}, (randX*randW + randY*randZ)},
	}

	for _, e := range dotProductTests {
		got, _ := DotProduct(e.v, e.u)
		if got != e.want {
			t.Errorf("DotProduct(%v, %v) == %v, want %v\n", e.v, e.u, got, e.want)
		}
	}
	/*
	  y*w - z*v
	  z*u - x*w
	  x*v - y*u
	*/
	var crossProductTests = []struct {
		v, u, want Vector
	}{
		{Vector{1.0, 2.0, 3.0}, Vector{3.0, 4.0, 5.0}, Vector{-2.0, 4.0, -2.0}},
		{Vector{xval, yval, zval}, Vector{uval, vval, wval},
			Vector{(yval*wval - zval*vval), (zval*uval - xval*wval), (xval*vval - yval*uval)}},
		{Vector{randX, randY, randZ}, Vector{randU, randV, randW},
			Vector{(randY*randW - randZ*randV), (randZ*randU - randX*randW), (randX*randV - randY*randU)}},
	}

	for _, f := range crossProductTests {
		got, _ := CrossProduct(f.v, f.u)
		if got.Equal(f.want) == false {
			t.Errorf("VecScalarMultiply(%v, %v) == %v, want %q\n", f.v, f.u, got, f.want)
		}
	}
	/*
	  // AngleFromOrigin calculates the angle of a given vector from the origin
	  // relative to the x-axis of 𝐄 (the model environment)
	  func AngleFromOrigin(v Vector) float64 {
	  	return math.Atan2(v[x], v[y])
	  }
	*/
	var angleFromOriginTests = []struct {
		v    Vector
		want float64
	}{
		{Vector{1.0, 2.0}, calc.ToFixed(math.Atan2(1.0, 2.0), 5)},
		{Vector{xval, yval}, calc.ToFixed(math.Atan2(xval, yval), 5)},
		{Vector{randW, randZ}, calc.ToFixed(math.Atan2(randW, randZ), 5)},
	}

	for _, g := range angleFromOriginTests {
		got, _ := AngleFromOrigin(g.v)
		if got != g.want {
			t.Errorf("AngleFromOrigin(%v), == %v, want %v\n", g.v, got, g.want)
		}
	}

	var relativeAngleTests = []struct {
		v, u Vector
		want float64
	}{
		{Vector{0, 0}, Vector{0, 1}, calc.ToFixed((math.Pi / 2), 5)},
		{Vector{1, 0}, Vector{-1, 0}, calc.ToFixed(math.Pi, 5)},
		{Vector{1, 0}, Vector{0, 0}, calc.ToFixed(math.Pi, 5)},
		{Vector{0, 0}, Vector{1, 0}, 0},
		{Vector{0, 1}, Vector{0, 0}, calc.ToFixed((-math.Pi / 2), 5)},
		{Vector{0, 0}, Vector{1, 1}, calc.ToFixed((math.Pi / 4), 5)},
		{Vector{0, 0}, Vector{-1, -1}, calc.ToFixed(((-3 * math.Pi) / 4), 5)},
		{Vector{0, 0}, Vector{-1, 1}, calc.ToFixed(((3 * math.Pi) / 4), 5)},
		{Vector{1, 0}, Vector{-1, 1}, 2.67795},
	}

	for _, r := range relativeAngleTests {
		got, _ := RelativeAngle(r.v, r.u)
		if got != r.want {
			t.Errorf("RelativeAngle(%v, %v) == %v, want %v\n", r.v, r.u, got, r.want)
		}
	}

	var unitAngleTests = []struct {
		a, want float64
	}{
		{2 * math.Pi, 0},
		{(2 * math.Pi) + randU, calc.ToFixed(randU, 5)},
		{(3 * math.Pi), calc.ToFixed(math.Pi, 5)},
	}

	for _, i := range unitAngleTests {
		got := UnitAngle(i.a)
		if got != i.want {
			t.Errorf("UnitAngle(%v), == %v, want %v\n", i.a, got, i.want)
		}
	}

	d0 := float64(4.35)
	r0 := float64(0.31)
	n0 := int((2.0 * d0) / (2.0 * r0)) // n0 = 14
	v0 := Vector{4.34, -4.34}
	d1 := float64(17.10189)
	r1 := float64(0.125)
	n1 := int((2 * d1) / (2 * r1)) // n1 = 136
	v1 := Vector{-3.056, 11.13}

	var translatePositionToSector2DTests = []struct {
		ed               float64
		n                int
		v                Vector
		rowWant, colWant int
	}{
		{1.0, int((2 * 1.0) / (2 * 0.2)), Vector{-0.18, 0.0}, 2, 2},
		{d0, n0, v0, 13, 13},
		{d1, n1, v1, 23, 55},
	}

	for _, j := range translatePositionToSector2DTests {
		rowGot, colGot := TranslatePositionToSector2D(j.ed, j.n, j.v)
		if rowGot != j.rowWant || colGot != j.colWant {
			t.Errorf("TranslatePositionToSector2D(%v, %v, %v) == (%v,%v), want (%v,%v)\n", j.ed, j.n, j.v, rowGot, colGot, j.rowWant, j.colWant)
		}
	}

	v2 := Vector{-0.00146, 9.13, 0.006}

	var magnitudeTests = []struct {
		v    Vector
		want float64
	}{
		{v0, 6.13769},
		{v1, 11.54193},
		{v2, 9.13000},
	}

	for _, k := range magnitudeTests {
		got, _ := Magnitude(k.v)
		if got != k.want {
			t.Errorf("Magnitude(%v) == %v, want %v\n", k.v, got, k.want)
		}
	}

	var normaliseTests = []struct {
		v    Vector
		want Vector
	}{
		{v0, Vector{0.70711, -0.70711}},
		{v1, Vector{-0.26477, 0.96431}},
		{v2, Vector{-0.00016, 1.0, 0.00066}},
	}

	for _, l := range normaliseTests {
		got, _ := Normalise(l.v)
		if got.Equal(l.want) == false {
			t.Errorf("Normalise(%v) == %v, want %v\n", l.v, got, l.want)
		}
	}
}
