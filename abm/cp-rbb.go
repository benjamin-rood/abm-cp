package abm

import "github.com/benjamin-rood/abm-colour-polymorphism/calc"

// RBB = Rule Based Behaviour that each cpp agent engages in once per turn, counts as the agent's action for that turn/phase.
func (c *ColourPolymorphicPrey) RBB(ctxt Context, popSize int) (newpop []ColourPolymorphicPrey) {
	newkids := []ColourPolymorphicPrey{}
	jump := ""
	// BEGIN
	jump = c.Age(ctxt)
	switch jump {
	case "DEATH":
		goto End
	case "SPAWN":
		progeny := c.Birth(ctxt) //	max spawn size, mutation factor
		newkids = append(newkids, progeny...)
	case "FERTILE":
		if popSize <= demoMaxCPP {
			c.Reproduction(ctxt.CppReproductiveChance, ctxt.CppGestation)
		}
		fallthrough
	case "EXPLORE":
		𝚯 := calc.RandFloatIn(-ctxt.CppTurn, ctxt.CppTurn)
		c.Turn(𝚯)
		c.Move()
	}

	newpop = append(newpop, *c)

End:
	newpop = append(newpop, newkids...) // add the newly created children to the returning population
	return
}
