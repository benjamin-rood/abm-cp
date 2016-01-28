package abm

import (
	"log"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
)

/* Attack(Hunted) bool
Eat(Hunted) bool */

// RBB : Rule-Based-Behaviour for Visual Predator Agent
func (vp *VisualPredator) RBB(ctxt Context, cppPop []ColourPolymorphicPrey) (returning []VisualPredator) {
	jump := ""
	jump = vp.Age(ctxt)
	_ = "breakpoint" // godebug
	switch jump {
	case "PREY SEARCH":
		target, err := vp.PreySearch(cppPop, ctxt.VpSearchChance) //	will move towards any viable prey it can see.
		if err != nil {
			log.Println("vp.RBB:", err)
		}
		if target != nil {
			vp.Attack(target, ctxt.VpAttackChance, ctxt.VpColImprintFactor)
			goto Add
		}
		fallthrough
	case "PATROL":
		𝚯 := calc.RandFloatIn(-vp.tr, vp.tr)
		vp.Turn(𝚯)
		vp.Move()
	case "DEATH":
		goto End
	default:
		log.Println("vp.RBB Switch: FAIL: jump =", jump)
	}
Add:
	returning = append(returning, *vp)
End:
	return
}
