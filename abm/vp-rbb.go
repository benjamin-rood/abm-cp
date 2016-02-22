package abm

import (
	"log"

	"github.com/benjamin-rood/abm-cp/calc"
)

// RBB : Rule-Based-Behaviour for Visual Predator Agent
func (vp *VisualPredator) RBB(errCh chan<- error, ctxt Context, start int, turn int, cppPop []ColourPolymorphicPrey, neighbours []VisualPredator, me int) []VisualPredator {
	var returning []VisualPredator
	var Φ float64
	popSize := len(neighbours)
	jump := vp.BeginTurn(ctxt, popSize)
	_ = "breakpoint" // godebug
	switch jump {
	case "DEATH":
		goto End
	case "PREY SEARCH":
		success := vp.SearchAndAttack(cppPop, ctxt, errCh)
		if success {
			goto Add
		}
		goto Patrol
	case "MATE SEARCH":
		if len(neighbours) <= 0 {
			goto Patrol
		}
		mate := vp.MateSearch(neighbours, me, errCh)
		success := vp.Copulation(mate, ctxt)
		if !success {
			goto Patrol
		}
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
func (vp *VisualPredator) SearchAndAttack(prey []ColourPolymorphicPrey, ctxt Context, errCh chan<- error) bool {
	var attacking bool
	var err error
	target, err := vp.PreySearch(prey) //	will move towards any viable prey it can see.
	errCh <- err
	if target != nil {
		attacking, err = vp.Intercept(target.pos)
		errCh <- err
	}
	if attacking {
		vp.Attack(target, ctxt)
		return true
	}
	return false
}
