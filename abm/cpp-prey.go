package abm

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/geometry"
	"github.com/benjamin-rood/abm-colour-polymorphism/render"
)

// ColourPolymorphicPrey – Prey agent type for Predator-Prey ABM
type ColourPolymorphicPrey struct {
	popIndex    int             //	index to the master population array.
	pos         geometry.Vector //	position in the environment
	movS        float64         //	speed
	movA        float64         //	acceleration
	dir𝚯        float64         //	 heading angle
	dir         geometry.Vector //	must be implemented as a unit vector
	Rτ          float64         // turn rate / range (in radians)
	sr          float64         //	search range
	lifespan    int
	hunger      int        //	counter for interval between needing food
	fertility   int        //	counter for interval between birth and sex
	gravid      bool       //	i.e. pregnant
	colouration colour.RGB //	colour
	𝛘           float64    //	 colour sorting value - colour distance/difference between vp.imprimt and cpp.colouration
	ϸ           float64    //  position sorting value - vector distance between vp.pos and cpp.pos
}

// String returns a clear textual presentation the internal values of the CPP agent
func (c *ColourPolymorphicPrey) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("cpp hereditary id{%d}:\n", c.popIndex))
	buffer.WriteString(fmt.Sprintf("pos=(%v,%v)\n", c.pos[x], c.pos[y]))
	buffer.WriteString(fmt.Sprintf("movS=%v\n", c.movS))
	buffer.WriteString(fmt.Sprintf("movA=%v\n", c.movA))
	buffer.WriteString(fmt.Sprintf("dir𝚯=%v\n", c.dir𝚯))
	buffer.WriteString(fmt.Sprintf("dir=(%v,%v)\n", c.dir[x], c.dir[y]))
	buffer.WriteString(fmt.Sprintf("Rτ=%v\n", c.Rτ))
	buffer.WriteString(fmt.Sprintf("sr=%v\n", c.sr))
	buffer.WriteString(fmt.Sprintf("lifespan=%v\n", c.lifespan))
	buffer.WriteString(fmt.Sprintf("hunger=%v\n", c.hunger))
	buffer.WriteString(fmt.Sprintf("fertility=%v\n", c.fertility))
	buffer.WriteString(fmt.Sprintf("gravid=%v\n", c.gravid))
	buffer.WriteString(fmt.Sprintf("colouration=%v\n", c.colouration))
	return buffer.String()
}

// GetDrawInfo exports the data set needed for agent visualisation.
func (c *ColourPolymorphicPrey) GetDrawInfo() (ar render.AgentRender) {
	ar.Type = "cpp"
	ar.X = c.pos[x]
	ar.Y = c.pos[y]
	ar.Colour = c.colouration.To256()
	return
}

// GeneratePopulation will create `size` number of agents
func GeneratePopulation(size int, context Context) []ColourPolymorphicPrey {
	var pop = []ColourPolymorphicPrey{}
	for i := 0; i < size; i++ {
		agent := ColourPolymorphicPrey{}
		agent.popIndex = i
		agent.pos = geometry.RandVector(context.Bounds)
		if context.CppAgeing {
			if context.RandomAges && (context.CppLifespan > 100) {
				agent.lifespan = calc.RandIntIn(5, context.CppLifespan)
			} else {
				agent.lifespan = context.CppLifespan
			}
		} else {
			agent.lifespan = -1 //	i.e. Undead!
		}
		agent.movS = context.CppS
		agent.movA = context.CppA
		agent.dir𝚯 = rand.Float64() * (2 * math.Pi)
		agent.dir = geometry.UnitVector(agent.dir𝚯)
		agent.Rτ = context.Cτ
		agent.sr = context.CppSr
		agent.hunger = 0
		agent.fertility = 1
		agent.gravid = false
		agent.colouration = colour.RandRGB()
		agent.𝛘 = 0.0
		agent.ϸ = 0.0
		pop = append(pop, agent)
	}
	return pop
}

// ProximitySort implements sort.Interface for []ColourPolymorphicPrey
// based on δ field.
type ProximitySort []ColourPolymorphicPrey

func (ps ProximitySort) Len() int           { return len(ps) }
func (ps ProximitySort) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
func (ps ProximitySort) Less(i, j int) bool { return ps[i].ϸ < ps[j].ϸ }

// VisualSort implements sort.Interface for []ColourPolymorphicPrey
// based on 𝛘 field – to assert visual bias of a VisualPredator based on it's colour imprinting value.
type VisualSort []ColourPolymorphicPrey

func (vs VisualSort) Len() int           { return len(vs) }
func (vs VisualSort) Swap(i, j int)      { vs[i], vs[j] = vs[j], vs[i] }
func (vs VisualSort) Less(i, j int) bool { return vs[i].𝛘 < vs[j].𝛘 }

/*
The Colour Polymorphic Prey agent is currently defined by the following animalistic interfaces:
Mover
Breeder
Mortal
*/

// Mover interface:

// Turn implements agent Mover interface method for ColourPolymorphicPrey:
// updates dir𝚯 and dir vector to the new heading offset by 𝚯
func (c *ColourPolymorphicPrey) Turn(𝚯 float64) {
	newHeading := geometry.UnitAngle(c.dir𝚯 + 𝚯)
	c.dir = geometry.UnitVector(newHeading)
	c.dir𝚯 = newHeading
}

// Move implements agent Mover interface method for ColourPolymorphicPrey:
// updates the agent's position according to its direction (heading) and
// velocity (speed*acceleration) if it doesn't encounter any errors.
func (c *ColourPolymorphicPrey) Move() error {
	var posOffset, newPos geometry.Vector
	var err error
	posOffset, err = geometry.VecScalarMultiply(c.dir, c.movS*c.movA)
	if err != nil {
		return errors.New("agent move failed: " + err.Error())
	}
	newPos, err = geometry.VecAddition(c.pos, posOffset)
	if err != nil {
		return errors.New("agent move failed: " + err.Error())
	}
	newPos[x] = calc.WrapFloatIn(newPos[x], -1.0, 1.0)
	newPos[y] = calc.WrapFloatIn(newPos[y], -1.0, 1.0)
	c.pos = newPos
	return nil
}

// Breeder interface:

// Fertility implements the blessed phases of the moon
func (c *ColourPolymorphicPrey) Fertility(Cȣ int) (jump string) {
	c.fertility++
	switch {
	case c.fertility == 0:
		c.gravid = false
		jump = "SPAWN"
		return
	case c.fertility >= Cȣ: // period / sexual cost
		jump = "MATE SEARCH"
		return
	}
	jump = "EXPLORE"
	return
}

// MateSearch implements Breeder interface method for ColourPolymorphicPrey:
// NEEDS BETTER HANDLING THAN JUST PUSHING THE ERROR UP!
func (c *ColourPolymorphicPrey) MateSearch(pop []ColourPolymorphicPrey, skip int) (mate *ColourPolymorphicPrey, err error) {
	mate = nil
	err = nil
	dist := 0.0
	for i := 0; i < len(pop); i++ {
		if i == skip {
			continue
		}
		dist, err = geometry.VectorDistance(c.pos, pop[i].pos)
		if err != nil {
			return
		}
		if dist <= c.sr {
			mate = &pop[i]
			return
		}
	}
	return
}

// Copulation implemets Breeder interface method for ColourPolymorphicPrey:
func (c *ColourPolymorphicPrey) Copulation(mate *ColourPolymorphicPrey, chance float64, gestation int, sexualCost int) bool {
	if mate == nil {
		return false
	}
	ω := rand.Float64()
	mate.fertility = -sexualCost // it takes two to tango, buddy!
	if ω <= chance {
		c.gravid = true
		c.fertility = -gestation
		return true
	}
	c.fertility = -sexualCost
	return false
}

// Birth implemets Breeder interface method for ColourPolymorphicPrey:
func (c *ColourPolymorphicPrey) Birth(ctxt Context) []ColourPolymorphicPrey {
	n := 1
	if ctxt.Cβ > 1 {
		n = rand.Intn(ctxt.Cβ) + 1 //	i.e. range [1, b]
	}
	progeny := GeneratePopulation(n, ctxt)
	for _, child := range progeny {
		child.mutation(c.colouration, ctxt.Mf)
		child.pos, _ = geometry.FuzzifyVector(c.pos, 0.1)
	}
	c.gravid = false
	return progeny
}

// For now, mutation only affects colouration, but could be extended to affect any other parameter.
func (c *ColourPolymorphicPrey) mutation(parentColour colour.RGB, Mf float64) {
	c.colouration = colour.RandRGBClamped(parentColour, Mf)
}

// set of methods implementing Mortal interface

// Age decrements the lifespan of an agent,
// and applies the effects of ageing (if any)
func (c *ColourPolymorphicPrey) Age() (jump string) {
	c.lifespan--
	c.hunger++
	if c.lifespan <= 0 {
		jump = "DEATH"
		return
	}
	jump = "HEALTHY"
	return
}

// Death – no specific functionality required so far.
func (c *ColourPolymorphicPrey) Death() {

}
