package calc

import (
	"math"
	"testing"
)

func Test(t *testing.T) {
	var roundTests = []struct {
		x    float64
		want int
	}{
		{0.3, 0},
		{0.9, 1},
		{0.1, 0},
		{0.5, 1},
	}

	for _, rt := range roundTests {
		got := round(rt.x)
		if got != rt.want {
			t.Errorf("round(%v) == %q, want %q\n", rt.x, got, rt.want)
		}
	}

	var toFixedTests = []struct {
		n    float64
		p    int
		want float64
	}{
		{math.Pi, 2, 3.14},
		{0.699304, 1, 0.7},
		{0.0566564, 3, 0.057},
		{0.973645, 4, 0.9736},
		{(1.76024973645) * math.Pow(10, 4), 4, 17602.4974},
	}

	for _, tf := range toFixedTests {
		got := ToFixed(tf.n, tf.p)
		if got != tf.want {
			t.Errorf("ToFixed(%v, %q) == %q, want %v\n", tf.n, tf.p, got, tf.want)
		}
	}

	var randFloatInTests = []struct {
		min, max float64
	}{
		{0.0, 1.0},
		{-1.0, 1.0},
		{1.76024973645, 33.99},
	}

	for _, rfi := range randFloatInTests {
		got := RandFloatIn(rfi.min, rfi.max)
		if (got < rfi.min) || (got >= rfi.max) {
			t.Errorf("RandFloatIn(%v, %v) == %v\n", rfi.min, rfi.max, got)
		}
	}

	var randIntInTests = []struct {
		min, max int
	}{
		{0, 10},
		{-10, 10},
		{-176024, 3399},
	}

	for _, rii := range randIntInTests {
		got := RandIntIn(rii.min, rii.max)
		if (got < rii.min) || (got >= rii.max) {
			t.Errorf("RandFloatIn(%v, %v) == %v\n", rii.min, rii.max, got)
		}
	}

	var clampFloatInTests = []struct {
		f, min, max, want float64
	}{
		{5.6751, 0.0, 1.0, 1.0},
		{5.6751, 1.0, 0.0, 5.6751},
		{0.6751, 0.0, 1.0, 0.6751},
		{-0.6751, 0.0, 1.0, 0.0},
		{-0.6751, -1.0, 1.0, -0.6751},
	}

	for _, cf := range clampFloatInTests {
		got := ClampFloatIn(cf.f, cf.min, cf.max)
		if got != cf.want {
			t.Errorf("ClampFloatIn(%v, %v, %v) == %v, want %v\n", cf.f, cf.min, cf.max, got, cf.want)
		}
	}
}
