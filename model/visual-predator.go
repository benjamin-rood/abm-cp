package model

import (
	"errors"
	"math"

	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/geometry"
)

// VPbehaviour – set of actions only VisualPredator agents will perform.
type VPbehaviour interface {
	VisualSearch([]ColourPolymorhicPrey, float64) (*ColourPolymorhicPrey, error)
	// ColourImprinting updates VP colour / visual recognition bias
	ColourImprinting(colour.RGB, float64) error
	VSRSectorSamples(float64, int) ([4][2]int, error)
	Turn(float64)
	Move()
}

// Turn updates dir𝚯 and dir vector to the new heading offset by 𝚯
func (vp *VisualPredator) Turn(𝚯 float64) {
	newHeading := geometry.UnitAngle(vp.dir𝚯 + 𝚯)
	vp.dir[x] = math.Cos(newHeading)
	vp.dir[y] = math.Sin(newHeading)
	vp.dir𝚯 = newHeading
}

// Move updates the agent's position if it doesn't encounter any errors.
func (vp *VisualPredator) Move() error {
	var posOffset, newPos geometry.Vector
	var err error
	posOffset, err = geometry.VecScalarMultiply(vp.dir, vp.movS*vp.movA)
	if err != nil {
		return errors.New("agent move failed: " + err.Error())
	}
	newPos, err = geometry.VecAddition(vp.pos, posOffset)
	if err != nil {
		return errors.New("agent move failed: " + err.Error())
	}
	vp.pos = newPos
	return nil
}

// VSRSectorSamples checks which sectors the VP agent's
// Visual Search Radius intersects.
// This initial version samples from 4 points on the circumference
// of the circle with radius vp.visRange originating at the VP agent's position
// The four sample points on the circumference at 45°, 135°, 225°, 315°
// or π/4, 3π/4, 5π/4, 7π/4 radians,
// or NE, NW, SW, SE on a compass, if you want to think of it that way :-)
func (vp *VisualPredator) VSRSectorSamples(d float64, n int) ([4][2]int, error) {
	sectorSamples := [4][2]int{}

	x45 := vp.pos[x] + (vp.vsr * (math.Cos(math.Pi / 4)))
	y45 := vp.pos[y] + (vp.vsr * (math.Sin(math.Pi / 4)))

	x135 := vp.pos[x] + (vp.vsr * (math.Cos(3 * math.Pi / 4)))
	y135 := vp.pos[y] + (vp.vsr * (math.Sin(3 * math.Pi / 4)))

	x225 := vp.pos[x] + (vp.vsr * (math.Cos(5 * math.Pi / 4)))
	y225 := vp.pos[y] + (vp.vsr * (math.Sin(5 * math.Pi / 4)))

	x315 := vp.pos[x] + (vp.vsr * (math.Cos(7 * math.Pi / 4)))
	y315 := vp.pos[y] + (vp.vsr * (math.Sin(7 * math.Pi / 4)))

	sectorSamples[0][0], sectorSamples[0][1] = geometry.TranslatePositionToSector2D(d, n, geometry.Vector{x45, y45})

	sectorSamples[1][0], sectorSamples[1][1] = geometry.TranslatePositionToSector2D(d, n, geometry.Vector{x135, y135})

	sectorSamples[2][0], sectorSamples[2][1] = geometry.TranslatePositionToSector2D(d, n, geometry.Vector{x225, y225})

	sectorSamples[3][0], sectorSamples[3][1] = geometry.TranslatePositionToSector2D(d, n, geometry.Vector{x315, y315})

	return sectorSamples, nil
}

// VisualSearch tries to 'recognise' a nearby prey agent to attack.
func (vp *VisualPredator) VisualSearch(population []ColourPolymorhicPrey, vsrSearchChance float64) (*ColourPolymorhicPrey, error) {
	for i := range population {
		population[i].𝛘 = colour.RGBDistance(vp.colImprint, population[i].colouration)
	}

	population = VisualSort(population)

	for i := range population {
		distanceToTarget, err := geometry.VectorDistance(vp.pos, population[i].pos)
		if err != nil {
			return nil, err
		}
		if distanceToTarget > vp.vsr {
			return nil, errors.New("VisualSearch failed")
		}
		if (distanceToTarget * vp.γ * population[i].𝛘) > vsrSearchChance {
			return &population[i], nil
		}
	}

	return nil, errors.New("VisualSearch failed")
}

// ColourImprinting updates VP colour / visual recognition bias
// Uses a bias / weighting value, 𝜎 (sigma) to control the degree of
// adaptation VP will make to differences in 'eaten' CPP colours.
func (vp *VisualPredator) ColourImprinting(target colour.RGB, 𝜎 float64) error {
	𝚫red := byte(float64(vp.colImprint.Red-target.Red) * 𝜎)
	𝚫green := byte(float64(vp.colImprint.Green-target.Green) * 𝜎)
	𝚫blue := byte(float64(vp.colImprint.Blue-target.Blue) * 𝜎)
	vp.colImprint.Red = vp.colImprint.Red - 𝚫red
	vp.colImprint.Green = vp.colImprint.Green - 𝚫green
	vp.colImprint.Blue = vp.colImprint.Blue - 𝚫blue
	return nil
}
