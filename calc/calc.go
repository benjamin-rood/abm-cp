package calc

import (
	"math"
	"math/rand"
)

/* Thanks to David Calhoun for providing round and toFixed!
http://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision-in-golang
*/
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// ToFixed will give a rounded-up version of "num" to "precision" decimal places.
func ToFixed(num float64, precision int) (output float64) {
	output = math.Pow(10, float64(precision))
	output = float64(round(num*output)) / output
	return
}

// RandFloatIn will give a random value in [min, max)
func RandFloatIn(min float64, max float64) float64 {
	return (rand.Float64() * (max - min)) + min
}

// RandIntIn will give a random value in [min, max)
func RandIntIn(min int, max int) int {
	return rand.Intn(max-min) + min
}

// ClampFloatIn will ensure that a floating point value is within range [min, max]. Dependant on min < max
func ClampFloatIn(f float64, min float64, max float64) float64 {
	if min >= max {
		return f
	}
	if f < min {
		return min
	}
	if f > max {
		return max
	}
	return f
}
