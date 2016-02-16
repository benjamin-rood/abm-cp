package abm

import (
	"log"

	"github.com/benjamin-rood/abm-cp/calc"
)

// RBB : Rule-Based-Behaviour for Visual Predator Agent
func (vp *VisualPredator) RBB(errCh chan<- error, ctxt Context, start int, turn int, cppPop []ColourPolymorphicPrey, vpPop []VisualPredator, me int) []VisualPredator {
	var returning []VisualPredator
	var Φ float64
	popSize := len(vpPop)
	jump := vp.Age(ctxt, popSize)
	_ = "breakpoint" // godebug
	switch jump {
	case "DEATH":
		goto End
	case "PREY SEARCH":
		success, err := vp.SearchAndAttack(cppPop, ctxt)
		errCh <- err
		if success {
			goto Add
		}
		goto Patrol
	case "FERTILE":
		mate, err := vp.MateSearch(vpPop, me)
		errCh <- err
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
	return returning
}

// SearchAndAttack gathers the logic for these steps of the VP RBB
func (vp *VisualPredator) SearchAndAttack(prey []ColourPolymorphicPrey, ctxt Context) (bool, error) {
	var attacking bool
	var err error
	target, δ, err := vp.PreySearch(prey) //	will move towards any viable prey it can see.
	if target != nil {
		attacking, err = vp.Intercept(target.pos, δ)
	}
	if err != nil {
		return false, err
	}
	if attacking {
		vp.Attack(target, ctxt)
		return true, nil
	}
	return false, nil
}
