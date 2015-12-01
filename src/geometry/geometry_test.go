package geometry

import (
	"math/rand"
	"testing"
	"time"
)

func Test(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	xval := float64(-0.5)
	yval := float64(0.5)
	wval := float64(1.1)
	zval := float64(-9.9)
	randX := rand.Float64()
	randY := rand.Float64()
	randW := rand.Float64()
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
			t.Errorf("VecAddition(%q, %q) == %q, want %q\n", c.v, c.u, got, c.want)
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

}
