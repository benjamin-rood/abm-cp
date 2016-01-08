package abm

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/geometry"
	"github.com/benjamin-rood/abm-colour-polymorphism/render"
)

var (
	cppLifespan = 20
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
			if context.RandomAges && (context.CppLifespan > 5) {
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
	fmt.Println("c.fertility =", c.fertility)
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
func (c *ColourPolymorphicPrey) MateSearch(pop []ColourPolymorphicPrey) (*ColourPolymorphicPrey, error) {
	for _, p := range pop {
		dist, err := geometry.VectorDistance(c.pos, p.pos)
		if err != nil {
			return nil, err
		}
		if dist <= c.sr {
			fmt.Println("found mate")
			return &p, nil
		}
	}
	return nil, nil
}

// Copulation implemets Breeder interface method for ColourPolymorphicPrey:
func (c *ColourPolymorphicPrey) Copulation(mate *ColourPolymorphicPrey, chance float64, gestation int, sexualCost int) bool {
	if mate == nil {
		return false
	}
	ω := rand.Float64()
	if ω <= chance {
		c.gravid = true
		c.fertility = -gestation
		return true
	}
	c.fertility = -sexualCost
	return false
}

// Birth implemets Breeder interface method for ColourPolymorphicPrey:
func (c *ColourPolymorphicPrey) Birth(b int, mf float64) []ColourPolymorphicPrey {
	n := 1
	if b > 1 {
		n = rand.Intn(b-1) + 1 //	i.e. range [1, b]
	}
	progeny := c.spawn(n)
	for i := 0; i < n; i++ {
		progeny[i].mutation(mf)
	}
	c.gravid = false
	return progeny
}

// set of actions only ColourPolymorphicPrey agents will perform
type cppBehaviour interface {
	newChild() ColourPolymorphicPrey
	mutation(float64) // variation at time of birth
	spawn(int) []ColourPolymorphicPrey
}

func (c *ColourPolymorphicPrey) newChild() ColourPolymorphicPrey {
	child := *c
	child.lifespan = cppLifespan
	child.pos, _ = geometry.FuzzifyVector(c.pos, 0.1)
	return child
}

func (c *ColourPolymorphicPrey) mutation(mf float64) {
	c.colouration = colour.RandRGBClamped(c.colouration, mf)
}

func (c *ColourPolymorphicPrey) spawn(n int) (progeny []ColourPolymorphicPrey) {
	for i := 0; i < n; i++ {
		child := c.newChild()
		progeny = append(progeny, child)
	}
	return
}

// Log prints all the internal values of the CPP agent
func (c *ColourPolymorphicPrey) Log() {
	log.Printf("cpp agent id{%d}:\n", c.popIndex)
	log.Printf("pos=(%v,%v)\n", c.pos[x], c.pos[y])
	log.Printf("movS=%v\n", c.movS)
	log.Printf("movA=%v\n", c.movA)
	log.Printf("dir𝚯=%v\n", c.dir𝚯)
	log.Printf("dir=(%v,%v)\n", c.dir[x], c.dir[y])
	log.Printf("Rτ=%v\n", c.Rτ)
	log.Printf("sr=%v\n", c.sr)
	log.Printf("lifespan=%v\n", c.lifespan)
	log.Printf("hunger=%v\n", c.hunger)
	log.Printf("fertility=%v\n", c.fertility)
	log.Printf("gravid=%v\n", c.gravid)
	log.Printf("colouration=%v\n", c.colouration)
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

// Death not implemented
func (c *ColourPolymorphicPrey) Death() {

}
