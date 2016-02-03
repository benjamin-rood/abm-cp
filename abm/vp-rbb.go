package abm

import (
	"log"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
)

// RBB : Rule-Based-Behaviour for Visual Predator Agent
func (vp *VisualPredator) RBB(ctxt Context, popSize int, cppPop []ColourPolymorphicPrey, vpPop []VisualPredator, me int) (returning []VisualPredator) {
	var Φ float64
	jump := vp.Age(ctxt)
	_ = "breakpoint" // godebug
	switch jump {
	case "DEATH":
		goto End
	case "PREY SEARCH":
		var attack bool
		var err error
		target, err := vp.PreySearch(cppPop, ctxt.VpSearchChance) //	will move towards any viable prey it can see.
		if target != nil {
			attack, err = vp.Intercept(target.pos, target.δ)
		}
		if err != nil {
			log.Println("vp.RBB:", err) // ARGH
			goto Patrol
		}
		if attack {
			vp.Attack(target, ctxt.VpAttackChance, ctxt.VpColImprintFactor, ctxt.Vγ)
			goto Add
		}
	case "FERTILE":
		if popSize >= ctxt.VpPopulationCap {
			goto Patrol
		}
		mate, err := vp.MateSearch(vpPop, me)
		if err != nil {
			log.Println("vp.RBB:", err) // ARGH
		}
		vp.Copulation(mate, ctxt.VpReproductionChance, ctxt.VpGestation, ctxt.VpSexualRequirement)
	case "SPAWN":
		children := vp.Birth(ctxt)
		returning = append(returning, children...)
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
