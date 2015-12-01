package geometry

import (
	"math/rand"
	"testing"
	"time"
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
}
