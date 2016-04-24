package abm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"

	"github.com/benjamin-rood/abm-cp/colour"
	"github.com/benjamin-rood/abm-cp/geometry"
	"github.com/benjamin-rood/abm-cp/render"
)

// UUID is just a getter method for the unexported uuid field, which absolutely must not change after agent creation.
func (vp *VisualPredator) UUID() string {
	return vp.uuid
}

// MarshalJSON implements json.Marshaler interface for VisualPredator object
func (vp VisualPredator) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"description":             vp.description,
		"pos":                     vp.pos,
		"speed":                   vp.movS,
		"heading":                 vp.𝚯,
		"turn-rate":               vp.tr,
		"search-range":            vp.vsr,
		"lifespan":                vp.lifespan,
		"hunger":                  vp.hunger,
		"attack-success":          vp.attackSuccess,
		"fertility":               vp.fertility,
		"𝛄":                       vp.𝛄,
		"gravid":                  vp.gravid,
		"colour-target-value":     vp.τ,
		"colour-imprint-strength": vp.ετ,
	})
}

// GetDrawInfo exports the data set needed for agent visualisation.
func (vp *VisualPredator) GetDrawInfo() (ar render.AgentRender) {
	ar.Type = "vp"
	ar.X = vp.pos[x]
	ar.Y = vp.pos[y]
	ar.Heading = vp.𝚯
	if vp.attackSuccess {
		// inv := vp.τ.Invert()
		// ar.Colour = inv.To256()
		ar.Colour = colour.RGB256{Red: 0, Green: 0, Blue: 0} // blink black on successful attack!
	} else {
		ar.Colour = vp.τ.To256()
	}
	return
}

type proxVP struct {
	comp func(geometry.Vector) float64
	*VisualPredator
}

type byProximityVp []proxVP

func (p byProximityVp) Len() int      { return len(p) }
func (p byProximityVp) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p byProximityVp) Less(i, j int) bool {
	return (p[i].comp(p[i].pos) < p[j].comp(p[j].pos)) //
}

// String returns a clear textual presentation the internal values of the VP agent
func (vp *VisualPredator) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("pos=(%v,%v)\n", vp.pos[x], vp.pos[y]))
	buffer.WriteString(fmt.Sprintf("movS=%v\n", vp.movS))
	buffer.WriteString(fmt.Sprintf("movA=%v\n", vp.movA))
	buffer.WriteString(fmt.Sprintf("𝚯=%v\n", vp.𝚯))
	buffer.WriteString(fmt.Sprintf("dir=(%v,%v)\n", vp.dir[x], vp.dir[y]))
	buffer.WriteString(fmt.Sprintf("tr=%v\n", vp.tr))
	buffer.WriteString(fmt.Sprintf("VpVsr=%v\n", vp.vsr))
	buffer.WriteString(fmt.Sprintf("lifespan=%v\n", vp.lifespan))
	buffer.WriteString(fmt.Sprintf("hunger=%v\n", vp.hunger))
	buffer.WriteString(fmt.Sprintf("fertility=%v\n", vp.fertility))
	buffer.WriteString(fmt.Sprintf("gravid=%v\n", vp.gravid))
	buffer.WriteString(fmt.Sprintf("τ=%v\n", vp.τ))
	buffer.WriteString(fmt.Sprintf("ετ=%v\n", vp.ετ))
	buffer.WriteString(fmt.Sprintf("𝛄=%v\n", vp.𝛄))
	return buffer.String()
}

func vpTesterAgent(xPos float64, yPos float64) (tester VisualPredator) {
	tester = vpTestPop(1)[0]
	tester.pos[x] = xPos
	tester.pos[y] = yPos
	return
}

// colourImprinting updates VP colour / visual recognition bias
// Uses a bias / weighting value, 𝜎 (sigma) to control the degree of
// adaptation VP will make to differences in 'eaten' CP Prey  colours.
func (vp *VisualPredator) colourImprinting(target colour.RGB, 𝜎 float64) {
	𝚫red := (vp.τ.Red - target.Red) * 𝜎
	𝚫green := (vp.τ.Green - target.Green) * 𝜎
	𝚫blue := (vp.τ.Blue - target.Blue) * 𝜎
	vp.τ.Red = vp.τ.Red - 𝚫red
	vp.τ.Green = vp.τ.Green - 𝚫green
	vp.τ.Blue = vp.τ.Blue - 𝚫blue
}

func vpTestPop(size int) []VisualPredator {
	return GenerateVPredatorPopulation(size, 0, 0, TestConditionParams, testStamp)
}

// VSRSectorSampling checks which sectors the VP agent's
// Visual Search Radius intersects.
// This initial version samples from 4 points on the circumference
// of the circle with radius vp.visRange originating at the VP agent's position
// The four sample points on the circumference at 45°, 135°, 225°, 315°
// or π/4, 3π/4, 5π/4, 7π/4 radians,
// or NE, NW, SW, SE on a compass, if you want to think of it that way :-)
func (vp *VisualPredator) VSRSectorSampling(d float64, n int) ([4][2]int, error) {
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
