package abm

import (
	"log"

	"github.com/benjamin-rood/abm-cp/calc"
)

// RBB : Rule-Based-Behaviour for Visual Predator Agent
func (vp *VisualPredator) RBB(errCh chan<- error, conditions ConditionParams, start int, turn int, cpPreyPop []ColourPolymorphicPrey, neighbours []VisualPredator, me int) []VisualPredator {
	var returning []VisualPredator
	var Φ float64
	popSize := len(neighbours)
	jump := vp.Age(conditions, popSize)
	switch jump {
	case "DEATH":
		goto End
	case "PREY SEARCH":
		success := vp.SearchAndAttack(cpPreyPop, conditions, errCh)
		if success {
			goto Add
		}
		goto Patrol
	case "FERTILE":
		if len(neighbours) <= 0 {
			goto Patrol
		}
		mate := vp.MateSearch(neighbours, me, errCh)
		success := vp.Copulation(mate, conditions)
		if !success {
			goto Patrol
		}
	case "SPAWN":
		children := vp.Birth(conditions, start, turn)
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
func (vp *VisualPredator) SearchAndAttack(prey []ColourPolymorphicPrey, conditions ConditionParams, errCh chan<- error) bool {
	var attacking bool
	var err error
	target, err := vp.PreySearch(prey) //	will move towards any viable prey it can see.
	errCh <- err
	if target != nil {
		attacking, err = vp.Intercept(target.pos)
		errCh <- err
	}
	if attacking {
		vp.Attack(target, conditions)
		return true
	}
	return false
}
