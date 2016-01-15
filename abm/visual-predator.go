package abm

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/geometry"
	"github.com/benjamin-rood/abm-colour-polymorphism/render"
)

// VisualPredator - Predator agent type for Predator-Prey ABM
type VisualPredator struct {
	pos        geometry.Vector //	position in the environment
	movS       float64         //	speed
	movA       float64         //	acceleration
	Rτ         float64         // turn rate / range (in radians)
	dir        geometry.Vector //	must be implemented as a unit vector
	dir𝚯       float64         //	 heading angle
	lifespan   int
	hunger     int     //	counter for interval between needing food
	fertility  int     //	counter for interval between birth and sex
	gravid     bool    //	i.e. pregnant
	vsr        float64 //	visual search range
	γ          float64 //	visual acuity (initially, use 1.0)
	colImprint colour.RGB
}

// GeneratePopulationVP will create `size` number of Visual Predator agents
func GeneratePopulationVP(size int, context Context) (pop []VisualPredator) {
	for i := 0; i < size; i++ {
		agent := VisualPredator{}
		agent.pos = geometry.RandVector(context.Bounds)
		if context.VpAgeing {
			if context.RandomAges {
				agent.lifespan = calc.RandIntIn(int(float64(context.VpLifespan)*0.7), int(float64(context.VpLifespan)*1.3))
			} else {
				agent.lifespan = context.VpLifespan
			}
		} else {
			agent.lifespan = 99999
		}
		agent.movS = context.VS
		agent.movA = context.VA
		agent.dir𝚯 = rand.Float64() * (2 * math.Pi)
		agent.dir = geometry.UnitVector(agent.dir𝚯)
		agent.Rτ = context.VpTurn
		agent.vsr = context.Vsr
		agent.γ = context.Vγ //	visual acuity
		agent.hunger = 0
		agent.fertility = 1
		agent.gravid = false
		agent.colImprint = colour.RandRGB()
		pop = append(pop, agent)
	}
	return
}

// GetDrawInfo exports the data set needed for agent visualisation.
func (vp *VisualPredator) GetDrawInfo() (ar render.AgentRender) {
	ar.Type = "vp"
	ar.X = vp.pos[x]
	ar.Y = vp.pos[y]
	ar.Heading = vp.dir𝚯
	ar.Colour = vp.colImprint.To256()
	return
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
	newPos[x] = calc.WrapFloatIn(newPos[x], -1.0, 1.0)
	newPos[y] = calc.WrapFloatIn(newPos[y], -1.0, 1.0)
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

// PreySearch – uses Visual Search to try to 'recognise' a nearby prey agent within model Environment to attack.
func (vp *VisualPredator) PreySearch(population []ColourPolymorphicPrey, searchChance float64) (target *ColourPolymorphicPrey, err error) {

	searchChance = 1.0 - searchChance
	_ = "breakpoint" // godebug
	for i := range population {
		population[i].𝛘 = colour.RGBDistance(vp.colImprint, population[i].colouration)
	}

	population = VisualSort(population)
	_ = "breakpoint" // godebug

	for i := range population {
		var distanceToTarget float64
		distanceToTarget, err = geometry.VectorDistance(vp.pos, population[i].pos)
		if err != nil {
			return
		}
		if distanceToTarget > vp.vsr {
			return
		}
		fmt.Println("distanceToTarget * vp.γ * population[i].𝛘 =", (distanceToTarget * vp.γ * population[i].𝛘))
		if (distanceToTarget * vp.γ * population[i].𝛘) > searchChance {
			target = &population[i]
			fmt.Println("target found =", *target)
			time.Sleep(time.Second)
			return
		}
	}
	return
}

// Attack VP agent attempts to attack CP prey agent
func (vp *VisualPredator) Attack(prey *ColourPolymorphicPrey, vpAttackChance float64) bool {
	if prey == nil {
		return false
	}
	vpAttackChance = 1 - vpAttackChance
	// Predator moves to prey position...
	angle, err := geometry.RelativeAngle(vp.pos, prey.pos)
	if err != nil {
		log.Println("geometry.RelativeAngle fail:", err)
		return false
	}
	vp.dir𝚯 = angle
	vp.dir = geometry.UnitVector(angle)
	vp.pos = prey.pos
	// ...Predatory tries to eat prey!
	α := rand.Float64()
	if α > vpAttackChance {
		vp.colourImprinting(prey.colouration, 1.0)
		vp.hunger -= 5
		prey.lifespan = 0 //	i.e. prey agent flagged for removal at the next turn.
		time.Sleep(1 * time.Second)
		fmt.Println("eaten =", *prey)
		return true
	}
	return false
}

// colourImprinting updates VP colour / visual recognition bias
// Uses a bias / weighting value, 𝜎 (sigma) to control the degree of
// adaptation VP will make to differences in 'eaten' CPP colours.
func (vp *VisualPredator) colourImprinting(target colour.RGB, 𝜎 float64) error {
	𝚫red := (vp.colImprint.Red - target.Red) * 𝜎
	𝚫green := (vp.colImprint.Green - target.Green) * 𝜎
	𝚫blue := (vp.colImprint.Blue - target.Blue) * 𝜎
	vp.colImprint.Red = vp.colImprint.Red - 𝚫red
	vp.colImprint.Green = vp.colImprint.Green - 𝚫green
	vp.colImprint.Blue = vp.colImprint.Blue - 𝚫blue
	return nil
}

// animal-agent Mortal interface methods:

// Age the vp agent
func (vp *VisualPredator) Age(ctxt Context) (jump string) {
	vp.hunger++
	if ctxt.VpAgeing {
		vp.lifespan--
	}
	switch {
	case vp.lifespan <= 0:
		jump = "DEATH"
	case vp.hunger > 0:
		jump = "PREY SEARCH"
	default:
		jump = "PATROL"
	}
	return
}
