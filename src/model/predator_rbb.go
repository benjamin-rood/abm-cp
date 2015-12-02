package model

import "errors"

// VPbehaviour – set of actions only VisualPredator agents will perform.
type VPbehaviour interface {
	VisualSearch([]ColourPolymorhicPrey, float64) (*ColourPolymorhicPrey, error)
	// ColourImprinting updates VP colour / visual recognition bias
	ColourImprinting(ColRGB, float64) error
	VSRSectorSamples(float64, int) [4][2]int
}

// VSRSectorSamples checks which sectors the VP agent's
// Visual Search Radius intersects.
// This initial version samples from 4 points on the circumference
// of the circle with radius vp.visRange originating at the VP agent's position
// The four sample points are at 0, π/2, π, 3π/2 radians.
func (vp *VisualPredator) VSRSectorSamples(d float64, n int) [4][2]int {
	sectorSamples := [4][2]int{}

	sectorSamples[0][0], sectorSamples[0][1] = TranslatePositionToSector2D(d, n, Vector{(vp.pos[x] + vp.visRange), vp.pos[y]})

	sectorSamples[1][0], sectorSamples[1][1] = TranslatePositionToSector2D(d, n, Vector{vp.pos[x], (vp.pos[y] + vp.visRange)})

	sectorSamples[2][0], sectorSamples[2][1] = TranslatePositionToSector2D(d, n, Vector{vp.pos[x] - vp.visRange, vp.pos[y]})

	sectorSamples[3][0], sectorSamples[3][1] = TranslatePositionToSector2D(d, n, Vector{vp.pos[x], (vp.pos[y] - vp.visRange)})

	return sectorSamples
}

// VisualSearch tries to 'recognise' a nearby prey agent to attack.
func (vp *VisualPredator) VisualSearch(population []ColourPolymorhicPrey, vsrSearchChance float64) (*ColourPolymorhicPrey, error) {
	for i := range population {
		population[i].𝛘 = ColourDistance(vp.colImprint, population[i].colour)
	}

	population = VisualSort(population)

	for i := range population {
		distanceToTarget, err := VectorDistance(vp.pos, population[i].pos)
		if err != nil {
			return nil, err
		}
		if distanceToTarget > vp.visRange {
			return nil, errors.New("VisualSearch failed")
		}
		if (distanceToTarget * vp.visAcuity * population[i].𝛘) > vsrSearchChance {
			return &population[i], nil
		}
	}

	return nil, errors.New("VisualSearch failed")
}

// ColourImprinting updates VP colour / visual recognition bias
// Uses a bias / weighting value, 𝜎 (sigma) to control the degree of
// adaptation VP will make to differences in 'eaten' CPP colours.
func (vp *VisualPredator) ColourImprinting(target ColRGB, 𝜎 float64) error {
	𝚫red := byte(float64(vp.colImprint.red-target.red) * 𝜎)
	𝚫green := byte(float64(vp.colImprint.green-target.green) * 𝜎)
	𝚫blue := byte(float64(vp.colImprint.blue-target.blue) * 𝜎)
	vp.colImprint = ColRGB{
		vp.colImprint.red - 𝚫red,
		vp.colImprint.green - 𝚫green,
		vp.colImprint.blue - 𝚫blue}
	return nil
}
