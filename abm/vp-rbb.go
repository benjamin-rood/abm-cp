package abm

import (
	"log"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
)

// RBB : Rule-Based-Behaviour for Visual Predator Agent
func (vp *VisualPredator) RBB(ctxt Context, cppPop []ColourPolymorphicPrey) (returning []VisualPredator) {
	var Φ float64
	var jump string
	jump = vp.Age(ctxt)
	_ = "breakpoint" // godebug
	switch jump {
	case "DEATH":
		goto End
	case "PREY SEARCH":
		target, err := vp.PreySearch(cppPop, ctxt.VpSearchChance) //	will move towards any viable prey it can see.
		attack, err := vp.Intercept(target)
		if err != nil {
			log.Println("vp.RBB:", err)
		}
		if attack {
			vp.Attack(target, ctxt.VpAttackChance, ctxt.VpColImprintFactor)
			goto Add
		}
		goto Patrol
	case "MATE SEARCH":
		goto Patrol //	i.e. not implemented yet.
	default:
		log.Println("vp.RBB Switch: FAIL: jump =", jump)
	}
Patrol:
	Φ = calc.RandFloatIn(-vp.tr, vp.tr)
	vp.Turn(Φ)
	vp.Move()
Add:
	returning = append(returning, *vp)
End:
	return
}
