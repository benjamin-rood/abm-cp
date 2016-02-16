package abm

import (
	"math"

	"github.com/benjamin-rood/abm-cp/geometry"
)

func visualSignalStrength(c float64) func(float64) float64 {
	var 𝓵 float64
	if c >= 1 {
		𝓵 = 1 / c
	} else {
		𝓵 = 1 - c
	}
	return func(𝛘 float64) float64 {
		return c*math.Exp(-c*𝛘) + 𝓵
	}
}

func visualSignalStrength2(c float64) func(float64) float64 {
	return func(𝛘 float64) float64 {
		return c * math.Exp(-c*𝛘)
	}
}

/*
type cd struct {
	comp func(float64) float64
	*MyType
}

type Comparitor []*cd

func (c Comparitor) Len() int      { return len(c) }
func (c Comparitor) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c Comparitor) Less(i, j int) bool {
	return c[i].comp(c[i].𝛘) < c[j].comp(c[j].𝛘)
}
*/

// The following code blocks are different approaches for SORTING sets of ColourPolymorphicPrey agents using sort.Sort():

type compCPP struct {
	comp func(geometry.Vector) float64
	*ColourPolymorphicPrey
}

type byComparitor []compCPP

func (c byComparitor) Len() int      { return len(c) }
func (c byComparitor) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c byComparitor) Less(i, j int) bool {
	return (c[i].comp(c[i].pos) < c[j].comp(c[j].pos)) //
}

type visualRecognition struct {
	δ    float64 //  position sorting value - vector distance between vp.pos and cpp.pos
	𝛘    float64 //	colour sorting value - colour distance/difference between vp.imprimt and cpp.colouration
	comp func(float64) float64
	rat  float64 //	value to rationalise the return from comp with
	*ColourPolymorphicPrey
}

type byVisualSignalStrength []visualRecognition

func (vss byVisualSignalStrength) Len() int      { return len(vss) }
func (vss byVisualSignalStrength) Swap(i, j int) { vss[i], vss[j] = vss[j], vss[i] }
func (vss byVisualSignalStrength) Less(i, j int) bool {
	return !(vss[i].comp(vss[i].𝛘) < vss[j].comp(vss[j].𝛘)) // As we want to sort Higher -> Lower values
}

type byOptimalAttackVector []visualRecognition

func (opt byOptimalAttackVector) Len() int      { return len(opt) }
func (opt byOptimalAttackVector) Swap(i, j int) { opt[i], opt[j] = opt[j], opt[i] }
func (opt byOptimalAttackVector) Less(i, j int) bool {
	return !((opt[i].comp(opt[i].𝛘) - opt[i].δ) < (opt[j].comp(opt[j].𝛘) - opt[j].δ)) // As we want to sort Higher -> Lower values
}

// byProximity implements sort.Interface for slice of *ColourPolymorphicPrey
// based on δ field.
type byProximity []visualRecognition

func (px byProximity) Len() int           { return len(px) }
func (px byProximity) Swap(i, j int)      { px[i], px[j] = px[j], px[i] }
func (px byProximity) Less(i, j int) bool { return px[i].δ < px[j].δ }

// byColourDifferentiation implements sort.Sort Interface for a slice of *ColourPolymorphicPrey
// based on 𝛘 field – to assert visual bias of a VisualPredator based on it's colour imprinting value.
type byColourDifferentiation []visualRecognition

func (vx byColourDifferentiation) Len() int           { return len(vx) }
func (vx byColourDifferentiation) Swap(i, j int)      { vx[i], vx[j] = vx[j], vx[i] }
func (vx byColourDifferentiation) Less(i, j int) bool { return vx[i].𝛘 < vx[j].𝛘 }

// byVisualDifferentiation implements sort.Sort Interface for a slice of *ColourPolymorphicPrey
// based on the sum of 𝛘 and δ fields
type byVisualDifferentiation []visualRecognition

func (vx byVisualDifferentiation) Len() int      { return len(vx) }
func (vx byVisualDifferentiation) Swap(i, j int) { vx[i], vx[j] = vx[j], vx[i] }
func (vx byVisualDifferentiation) Less(i, j int) bool {
	return (vx[i].𝛘 + vx[i].δ) < (vx[j].𝛘 + vx[j].δ)
}
