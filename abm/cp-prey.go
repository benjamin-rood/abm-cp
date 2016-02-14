package abm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/benjamin-rood/abm-cp/calc"
	"github.com/benjamin-rood/abm-cp/colour"
	"github.com/benjamin-rood/abm-cp/geometry"
	"github.com/benjamin-rood/abm-cp/render"
)

// ColourPolymorphicPrey – Prey agent type for Predator-Prey ABM
type ColourPolymorphicPrey struct {
	uuid        string //	do not export this field
	description AgentDescription
	pos         geometry.Vector //	position in the environment
	movS        float64         //	speed
	movA        float64         //	acceleration
	𝚯           float64         //	 heading angle
	dir         geometry.Vector //	must be implemented as a unit vector
	tr          float64         // turn rate / range (in radians)
	sr          float64         //	search range
	lifespan    int
	hunger      int        //	counter for interval between needing food
	fertility   int        //	counter for interval between birth and sex
	gravid      bool       //	i.e. pregnant
	colouration colour.RGB //	colour
}

// UUID is just a getter method for the unexported uuid field, which absolutely must not change after agent creation.
func (c *ColourPolymorphicPrey) UUID() string {
	return c.uuid
}

// MarshalJSON implements json.Marshaler interface on a CPP object
func (c ColourPolymorphicPrey) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"description":  c.description,
		"pos":          c.pos,
		"speed":        c.movS,
		"heading":      c.𝚯,
		"turn-rate":    c.tr,
		"search-range": c.sr,
		"lifespan":     c.lifespan,
		"hunger":       c.hunger,
		"fertility":    c.fertility,
		"colouration":  c.colouration,
	})
}

// GetDrawInfo exports the data set needed for agent visualisation.
func (c ColourPolymorphicPrey) GetDrawInfo() (ar render.AgentRender) {
	ar.Type = "cpp"
	ar.X = c.pos[x]
	ar.Y = c.pos[y]
	ar.Colour = c.colouration.To256()
	return
}

// GeneratePopulationCPP will create `size` number of Colour Polymorphic Prey agents
func GeneratePopulationCPP(size int, start int, mt int, context Context, timestamp string) []ColourPolymorphicPrey {
	pop := []ColourPolymorphicPrey{}
	for i := 0; i < size; i++ {
		agent := ColourPolymorphicPrey{}
		agent.uuid = uuid()
		agent.description = AgentDescription{AgentType: "cpp", AgentNum: start + i, ParentUUID: "", CreatedMT: mt, CreatedAT: timestamp}
		agent.pos = geometry.RandVector(context.Bounds)
		if context.CppAgeing {
			if context.RandomAges {
				agent.lifespan = calc.RandIntIn(int(float64(context.CppLifespan)*0.7), int(float64(context.CppLifespan)*1.3))
			} else {
				agent.lifespan = context.CppLifespan
			}
		} else {
			agent.lifespan = 99999
		}
		agent.movS = context.CppS
		agent.movA = context.CppA
		agent.𝚯 = rand.Float64() * (2 * math.Pi)
		agent.dir = geometry.UnitVector(agent.𝚯)
		agent.tr = context.CppTurn
		agent.sr = context.CppSr
		agent.hunger = 0
		agent.fertility = 1
		agent.gravid = false
		agent.colouration = colour.RandRGB()
		pop = append(pop, agent)
	}
	return pop
}

func cppSpawn(size int, parent ColourPolymorphicPrey, context Context, timestamp string) []ColourPolymorphicPrey {
	pop := []ColourPolymorphicPrey{}
	for i := 0; i < size; i++ {
		agent := parent
		agent.uuid = uuid()
		agent.pos = parent.pos
		if context.CppAgeing {
			if context.RandomAges {
				agent.lifespan = calc.RandIntIn(int(float64(context.CppLifespan)*0.7), int(float64(context.CppLifespan)*1.3))
			} else {
				agent.lifespan = context.CppLifespan
			}
		} else {
			agent.lifespan = 99999 //	i.e. Undead!
		}
		agent.movS = parent.movS
		agent.movA = parent.movA
		agent.𝚯 = rand.Float64() * (2 * math.Pi)
		agent.dir = geometry.UnitVector(agent.𝚯)
		agent.tr = parent.tr
		agent.sr = parent.sr
		agent.hunger = 0
		agent.fertility = 1
		agent.gravid = false
		agent.colouration = parent.colouration
		pop = append(pop, agent)
	}
	return pop
}

/*
The Colour Polymorphic Prey agent is currently defined by the following animalistic interfaces:
Mover
Breeder
Mortal
*/

// Mover interface:

// Turn implements agent Mover interface method for ColourPolymorphicPrey:
// updates 𝚯 and dir vector to the new heading offset by 𝚯
func (c *ColourPolymorphicPrey) Turn(𝚯 float64) {
	newHeading := geometry.UnitAngle(c.𝚯 + 𝚯)
	c.dir = geometry.UnitVector(newHeading)
	c.𝚯 = newHeading
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

// Reproduction implements Breeder interface method - ASEXUAL (self-reproduction) ColourPolymorphicPrey:
func (c *ColourPolymorphicPrey) Reproduction(chance float64, gestation int) bool {
	c.hunger++ //	energy cost
	ω := rand.Float64()
	if ω <= chance {
		c.gravid = true
		c.fertility = -gestation
		return true
	}
	c.fertility = 1
	return false
}

// Copulation implemets Breeder interface method for ColourPolymorphicPrey:
func (c *ColourPolymorphicPrey) copulation(mate *ColourPolymorphicPrey, chance float64, gestation int, sexualCost int) bool {
	if mate == nil {
		return false
	}
	if mate.fertility < sexualCost { //	mate must be sufficiently fertile also
		return false
	}
	ω := rand.Float64()
	mate.fertility = -sexualCost // it takes two to tango, buddy!
	if ω <= chance {
		c.gravid = true
		c.fertility = -gestation
		return true
	}
	c.fertility = 1
	return false
}

// Birth implemets Breeder interface method for ColourPolymorphicPrey:
func (c *ColourPolymorphicPrey) Birth(ctxt Context) []ColourPolymorphicPrey {
	n := 1
	if ctxt.CppSpawnSize > 1 {
		n = rand.Intn(ctxt.CppSpawnSize) + 1 //	i.e. range [1, b]
	}
	timestamp := fmt.Sprintf("%s", time.Now())
	progeny := cppSpawn(n, *c, ctxt, timestamp)
	for i := 0; i < len(progeny); i++ {
		progeny[i].mutation(ctxt.CppMutationFactor)
		progeny[i].pos, _ = geometry.FuzzifyVector(c.pos, c.movS)
	}
	c.hunger++ //	energy cost
	c.gravid = false
	return progeny
}

// For now, mutation only affects colouration, but could be extended to affect any other parameter.
func (c *ColourPolymorphicPrey) mutation(Mf float64) {
	c.colouration = colour.RandRGBClamped(c.colouration, Mf)
}

// set of methods implementing Mortal interface

// Age decrements the lifespan of an agent,
// and applies the effects of ageing (if any)
func (c *ColourPolymorphicPrey) Age(ctxt Context) (jump string) {
	c.hunger++
	c.fertility++
	if ctxt.CppAgeing {
		c.lifespan--
	}
	switch {
	case c.lifespan <= 0:
		jump = "DEATH"
	case c.fertility == 0:
		c.gravid = false
		jump = "SPAWN"
	case c.fertility >= ctxt.CppSexualCost: // period / sexual cost
		jump = "FERTILE"
	default:
		jump = "EXPLORE"
	}
	return
}

// String returns a clear textual presentation the internal values of the CPP agent
func (c *ColourPolymorphicPrey) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("pos=(%v,%v)\n", c.pos[x], c.pos[y]))
	buffer.WriteString(fmt.Sprintf("movS=%v\n", c.movS))
	buffer.WriteString(fmt.Sprintf("movA=%v\n", c.movA))
	buffer.WriteString(fmt.Sprintf("𝚯=%v\n", c.𝚯))
	buffer.WriteString(fmt.Sprintf("dir=(%v,%v)\n", c.dir[x], c.dir[y]))
	buffer.WriteString(fmt.Sprintf("tr=%v\n", c.tr))
	buffer.WriteString(fmt.Sprintf("sr=%v\n", c.sr))
	buffer.WriteString(fmt.Sprintf("lifespan=%v\n", c.lifespan))
	buffer.WriteString(fmt.Sprintf("hunger=%v\n", c.hunger))
	buffer.WriteString(fmt.Sprintf("fertility=%v\n", c.fertility))
	buffer.WriteString(fmt.Sprintf("gravid=%v\n", c.gravid))
	buffer.WriteString(fmt.Sprintf("colouration=%v\n", c.colouration))
	return buffer.String()
}
