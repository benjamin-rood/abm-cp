package abm

import (
	"log"

	"github.com/benjamin-rood/abm-cp/calc"
)

// RBB : Rule-Based-Behaviour for Visual Predator Agent
func (vp *VisualPredator) RBB(ctxt Context, start int, turn int, cppPop []ColourPolymorphicPrey, vpPop []VisualPredator, me int) (returning []VisualPredator) {
	var Φ float64
	popSize := len(vpPop)
	jump := vp.Age(ctxt, popSize)
	_ = "breakpoint" // godebug
	switch jump {
	case "DEATH":
		goto End
	case "PREY SEARCH":
		var attack bool
		var err error
		target, δ, err := vp.PreySearch(cppPop, ctxt.VpSearchChance) //	will move towards any viable prey it can see.
		if target != nil {
			attack, err = vp.Intercept(target.pos, δ)
		}
		if err != nil {
			log.Println("vp.RBB:", err) // ARGH
		}
		if attack {
			vp.Attack(target, ctxt.VpAttackChance, ctxt.VpCaf, ctxt.Vbg, ctxt.Vbγ, ctxt.Vbε)
			goto Add
		}
		goto Patrol
	case "FERTILE":
		mate, err := vp.MateSearch(vpPop, me)
		if err != nil {
			log.Println("vp.RBB:", err) // ARGH
		}
		vp.Copulation(mate, ctxt.VpReproductionChance, ctxt.VpGestation, ctxt.VpSexualRequirement)
	case "SPAWN":
		// func (vp *VisualPredator) Birth(ctxt Context, start int, mt int)
		children := vp.Birth(ctxt, start, turn)
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
