package geometry

import (
	"math"
	"testing"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
)

func TestAngleToIntercept(t *testing.T) {
	// baseline initial test:
	samplePos := Vector{0, 0}
	sampleHeading := calc.ToFixed(3*math.Pi/4, 5) //  pointing towards (-1,1), halfway between 𝛑/2 and 𝛑 radians
	sampleTarget := Vector{1, -1}
	// result should be how MUCH we have to turn one way or another:
	Ψ, _ := AngleToIntercept(samplePos, sampleHeading, sampleTarget)
	// we disregard error value from AngleToIntercept as the only possible error would be from mismatched Vector lengths or Vectors of length != 2 which do not apply here.
	want := -3.14159
	if Ψ != want {
		t.Errorf("AngleToIntercept(%v, %v, %v) == %v, want: %v", samplePos, sampleHeading, sampleTarget, Ψ, want)
	}

	sampleTarget = Vector{-1, -1} // angle of -3𝛑/4 radians relative to samplePos
	Ψ, _ = AngleToIntercept(samplePos, sampleHeading, sampleTarget)
	want = 1.57081
	if Ψ != want {
		t.Errorf("AngleToIntercept(%v, %v, %v) == %v, want: %v", samplePos, sampleHeading, sampleTarget, Ψ, want)
	}
}
