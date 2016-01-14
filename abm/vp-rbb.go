package abm

import (
	"log"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
)

/* Attack(Hunted) bool
Eat(Hunted) bool */

// RBB : Rule-Based-Behaviour for Visual Predator Agent
func (vp *VisualPredator) RBB(ctxt Context, cppPop []ColourPolymorphicPrey, eaten *ColourPolymorphicPrey) *VisualPredator {
	jump := ""
	jump = vp.Age(ctxt)

	switch jump {
	case "PREY SEARCH":
		target, err := vp.PreySearch(cppPop, ctxt.VsrSearchChance)
		if err != nil {
			log.Println("vp.RBB:", err)
		}
		success := vp.Attack(target, ctxt.VpAttackChance)
		if success {
			eaten = target
			vp.hunger -= 5
			return vp
		}
		fallthrough
	case "PATROL":
		𝚯 := calc.RandFloatIn(-ctxt.VpTurn, ctxt.VpTurn)
		vp.Turn(𝚯)
		vp.Move()
		return vp
	case "DEATH":
		return nil
	default:
		log.Println("vp.RBB Switch: FAIL: jump =", jump)
		return vp
	}
}
