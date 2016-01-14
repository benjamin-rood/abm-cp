package abm

import (
	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
	"github.com/benjamin-rood/abm-colour-polymorphism/render"
)

/* Attack(Hunted) bool
Eat(Hunted) bool */

// RBB : Rule-Based-Behaviour for Visual Predator Agent
func (vp *VisualPredator) RBB(ctxt Context, render <-chan render.AgentRender, cppPop []ColourPolymorphicPrey, eaten *ColourPolymorphicPrey) {
	jump := ""
	jump = vp.Age(ctxt)

	switch jump {
	case "PREY SEARCH":
		target := vp.PreySearch(cppPop, ctxt.VsrSearchChance)
		success := vp.Attack(target, ctxt.VpAttackChance)
		if success {
			eaten = target
			vp.hunger -= 5
		}
		fallthrough
	case "PATROL":
		𝚯 := calc.RandFloatIn(-ctxt.VpTurn, ctxt.VpTurn)
		vp.Turn(𝚯)
		vp.Move()
		return
	case "DEATH":
		return
	}
}
