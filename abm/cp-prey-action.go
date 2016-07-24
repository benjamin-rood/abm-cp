package abm

import "github.com/benjamin-rood/abm-cp/calc"

// Action = Rule Based Behaviour that each cpPrey agent engages in once per turn, counts as the agent's action for that turn/phase.
func (c *ColourPolymorphicPrey) Action(conditions ConditionParams, popSize int) (newpop []ColourPolymorphicPrey) {
  newkids := []ColourPolymorphicPrey{}
  jump := ""
  // BEGIN
  jump = c.Age(conditions)
  switch jump {
  case "DEATH":
    goto End
  case "SPAWN":
    progeny := c.Birth(conditions) //	max spawn size, mutation factor
    newkids = append(newkids, progeny...)
  case "FERTILE":
    if popSize <= conditions.CpPreyPopulationCap {
      c.Reproduction(conditions.CpPreyReproductionChance, conditions.CpPreyGestation)
    }
    fallthrough
  case "EXPLORE":
    𝚯 := calc.RandFloatIn(-conditions.CpPreyTurn, conditions.CpPreyTurn)
    c.Turn(𝚯)
    c.Move()
  }

  newpop = append(newpop, *c)

End:
  newpop = append(newpop, newkids...) // add the newly created children to the returning population
  return
}
