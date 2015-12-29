package model

import (
	"errors"
	"math"

	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/geometry"
)

// ColourPolymorhicPrey – Prey agent type for Predator-Prey ABM
type ColourPolymorhicPrey struct {
	populationIndex uint            //	index to the master population array.
	pos             geometry.Vector //	position in the environment
	movS            float64         //	speed
	movA            float64         //	acceleration
	dir             geometry.Vector //	must be implemented as a unit vector
	dir𝚯            float64         //	 heading angle
	hunger          uint            //	counter for interval between needing food
	fertility       uint            //	counter for interval between birth and sex
	gravid          bool            //	i.e. pregnant
	colouration     colour.RGB      //	colour
	𝛘               float64         //	 colour sorting value - colour distance/difference between vp.imprimt and cpp.colouration
	ϸ               float64         //  position sorting value - vector distance between vp.pos and cpp.pos
}

// cppBehaviour – set of actions only VisualPredator agents will perform – unexported!
type cppBehaviour interface {
	mutation() // variation at time of birth
	spawn() []ColourPolymorhicPrey
}

// ProximitySort implements sort.Interface for []ColourPolymorhicPrey
// based on δ field.
type ProximitySort []ColourPolymorhicPrey

func (ps ProximitySort) Len() int           { return len(ps) }
func (ps ProximitySort) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
func (ps ProximitySort) Less(i, j int) bool { return ps[i].ϸ < ps[j].ϸ }

// VisualSort implements sort.Interface for []ColourPolymorhicPrey
// based on 𝛘 field – to assert visual bias of a VisualPredator based on it's colour imprinting value.
type VisualSort []ColourPolymorhicPrey

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

// Turn implements agent Mover interface method for ColourPolymorhicPrey:
// updates dir𝚯 and dir vector to the new heading offset by 𝚯
func (c *ColourPolymorhicPrey) Turn(𝚯 float64) {
	newHeading := geometry.UnitAngle(c.dir𝚯 + 𝚯)
	c.dir[x] = math.Cos(newHeading)
	c.dir[y] = math.Sin(newHeading)
	c.dir𝚯 = newHeading
}

// Move implements agent Mover interface method for ColourPolymorhicPrey:
// updates the agent's position according to its direction (heading) and
// velocity (speed*acceleration) if it doesn't encounter any errors.
func (c *ColourPolymorhicPrey) Move() error {
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
	c.pos = newPos
	return nil
}

// Breeder interface:

// MateSearch implements Breeder interface method for ColourPolymorhicPrey:
func (c *ColourPolymorhicPrey) MateSearch(pop []ColourPolymorhicPrey) (bool, *ColourPolymorhicPrey) {

}
