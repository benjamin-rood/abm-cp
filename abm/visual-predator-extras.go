package abm

import (
	"bytes"
	"fmt"
	"math"

	"github.com/benjamin-rood/abm-cp/geometry"
)

// String returns a clear textual presentation the internal values of the VP agent
func (vp *VisualPredator) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("pos=(%v,%v)\n", vp.pos[x], vp.pos[y]))
	buffer.WriteString(fmt.Sprintf("movS=%v\n", vp.movS))
	buffer.WriteString(fmt.Sprintf("movA=%v\n", vp.movA))
	buffer.WriteString(fmt.Sprintf("𝚯=%v\n", vp.𝚯))
	buffer.WriteString(fmt.Sprintf("dir=(%v,%v)\n", vp.dir[x], vp.dir[y]))
	buffer.WriteString(fmt.Sprintf("tr=%v\n", vp.tr))
	buffer.WriteString(fmt.Sprintf("Vsr=%v\n", vp.vsr))
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

func vpTestPop(size int) []VisualPredator {
	return GeneratePopulationVP(size, 0, 0, TestContext, testStamp)
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
