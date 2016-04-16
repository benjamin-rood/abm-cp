package abm

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/benjamin-rood/abm-cp/calc"
	"github.com/benjamin-rood/abm-cp/colour"
	"github.com/benjamin-rood/abm-cp/geometry"
)

// VisualPredator - Predator agent type for Predator-Prey ABM
type VisualPredator struct {
	uuid          string //	do not export this field!
	description   AgentDescription
	pos           geometry.Vector //	position in the environment
	movS          float64         //	speed	/ movement range per turn
	movA          float64         //	acceleration
	tr            float64         // turn rate / range (in radians)
	dir           geometry.Vector //	must be implemented as a unit vector
	𝚯             float64         //	 heading angle
	lifespan      int
	hunger        int        //	counter for interval between needing food
	attackSuccess bool       //	if during the turn, the VP agent successfully ate a CP prey agent
	fertility     int        //	counter for interval between birth and sex
	gravid        bool       //	i.e. pregnant
	vsr           float64    //	visual search range
	𝛄             float64    //	visual seach (colour) bias
	τ             colour.RGB //	imprinted target / colour specialisation value
	ετ            float64    //	imprinting / colour specialisation strength
}

// GeneratePopulationVP will create `size` number of Visual Predator agents
func GeneratePopulationVP(size int, start int, mt int, context Context, timestamp string) []VisualPredator {
	pop := []VisualPredator{}
	for i := 0; i < size; i++ {
		agent := VisualPredator{}
		agent.uuid = uuid()
		agent.description = AgentDescription{AgentType: "vp", AgentNum: start + i, ParentUUID: "", CreatedMT: mt, CreatedAT: timestamp}
		agent.pos = geometry.RandVector(context.Bounds)
		if context.VpAgeing {
			if context.RandomAges {
				agent.lifespan = calc.RandIntIn(int(float64(context.VpLifespan)*0.7), int(float64(context.VpLifespan)*1.3))
			} else {
				agent.lifespan = context.VpLifespan
			}
		} else {
			agent.lifespan = 99999
		}
		agent.movS = context.VpMovS
		agent.movA = context.VpMovA
		agent.𝚯 = rand.Float64() * (2 * math.Pi)
		agent.dir = geometry.UnitVector(agent.𝚯)
		agent.tr = context.VpTurn
		agent.vsr = context.Vsr
		agent.𝛄 = context.Vb𝛄 //	baseline acuity level
		agent.hunger = context.VpSexualRequirement + 1
		agent.fertility = 1
		agent.gravid = false
		agent.τ = colour.RandRGB()
		agent.ετ = context.Vbε
		pop = append(pop, agent)
	}
	return pop
}

func vpSpawn(size int, start int, mt int, parent VisualPredator, context Context, timestamp string) []VisualPredator {
	pop := []VisualPredator{}
	for i := 0; i < size; i++ {
		agent := parent
		agent.uuid = uuid()
		agent.description = AgentDescription{AgentType: "vp", AgentNum: start + i, ParentUUID: parent.uuid, CreatedMT: mt, CreatedAT: timestamp}
		agent.pos = parent.pos
		if context.VpAgeing {
			if context.RandomAges {
				agent.lifespan = calc.RandIntIn(int(float64(context.VpLifespan)*0.7), int(float64(context.VpLifespan)*1.3))
			} else {
				agent.lifespan = context.VpLifespan
			}
		} else {
			agent.lifespan = 99999
		}
		agent.movS = parent.movS
		agent.movA = parent.movA
		agent.𝚯 = rand.Float64() * (2 * math.Pi)
		agent.dir = parent.dir
		agent.tr = parent.tr
		agent.vsr = parent.vsr
		agent.hunger = context.VpSexualRequirement + 1
		agent.fertility = 1
		agent.gravid = false
		agent.τ = colour.RandRGBClamped(parent.τ, 0.5) //	random offset (up to 50%) deviation from parent's target colour
		agent.ετ = context.Vbε
		pop = append(pop, agent)
	}
	return pop
}

// Turn updates 𝚯 and dir vector to the new heading offset by 𝚯
func (vp *VisualPredator) Turn(𝚯 float64) {
	newHeading := geometry.UnitAngle(vp.𝚯 + 𝚯)
	vp.dir[x] = math.Cos(newHeading)
	vp.dir[y] = math.Sin(newHeading)
	vp.𝚯 = newHeading
}

// Move updates the agent's position if it doesn't encounter any errors.
func (vp *VisualPredator) Move() error {
	var posOffset, newPos geometry.Vector
	var err error
	posOffset, err = geometry.VecScalarMultiply(vp.dir, vp.movS*vp.movA)
	if err != nil {
		return errors.New("agent move failed: " + err.Error())
	}
	newPos, err = geometry.VecAddition(vp.pos, posOffset)
	if err != nil {
		return errors.New("agent move failed: " + err.Error())
	}
	newPos[x] = calc.WrapFloatIn(newPos[x], -1.0, 1.0)
	newPos[y] = calc.WrapFloatIn(newPos[y], -1.0, 1.0)
	vp.pos = newPos
	return nil
}

// PreySearch – uses Visual Search to try to 'recognise' a nearby prey agent within model Environment to target
func (vp *VisualPredator) PreySearch(prey []ColourPolymorphicPrey) (*ColourPolymorphicPrey, error) {
	c := vp.ετ
	// var 𝒇 = visualSignalStrength(c)
	var 𝒇 = visualSignalStrength2(c)
	var 𝛘 float64 // colour sorting value - colour distance/difference between vp.imprimt and cpp.colouration
	var δ float64 // position sorting value - vector distance between vp.pos and cpp.pos
	var err error
	var searchSet []visualRecognition
	for i := range prey { //	exhaustive search 😱
		δ, err = geometry.VectorDistance(vp.pos, prey[i].pos)
		// fmt.Printf("δ=%v\t\tvsr=%v\n", δ, vp.vsr)
		if δ <= vp.vsr { // ∴ only include the prey agent for considertion if within visual range
			𝛘 = colour.RGBDistance(vp.τ, prey[i].colouration)
			// fmt.Printf("𝛘=%v\t\t𝛄=%v\n", 𝛘, vp.𝛄)
			if 𝛘 < vp.𝛄 { // i.e. if and only if colour distance falls within predator's current acuity
				a := visualRecognition{δ, 𝛘, 𝒇, c, &prey[i]}
				searchSet = append(searchSet, a)
			}
		}
	}

	// for i := range searchSet {
	// 	fmt.Printf("%v\tδ=%v\t𝛘=%v\tc=%v\t%p\t%v\t%v\n", i, searchSet[i].δ, searchSet[i].𝛘, c, searchSet[i].ColourPolymorphicPrey, 𝒇(searchSet[i].𝛘), 𝒇(searchSet[i].𝛘)-searchSet[i].δ)
	// }

	sort.Sort(byOptimalAttackVector(searchSet)) //	sort by 𝒇(x) - distance

	// for i := range searchSet {
	// 	fmt.Printf("%v\tδ=%v\t𝛘=%v\tc=%v\t%p\t%v\t%v\n", i, searchSet[i].δ, searchSet[i].𝛘, c, searchSet[i].ColourPolymorphicPrey, 𝒇(searchSet[i].𝛘), 𝒇(searchSet[i].𝛘)-searchSet[i].δ)
	// }

	// search within biased and reduced set
	for i, p := range searchSet {
		if 𝒇(p.𝛘) > (1 - vp.𝛄) { // i.e. is the colour detection strength sufficiently great
			return &(*searchSet[i].ColourPolymorphicPrey), err
		}
	}
	return nil, err
}

// Attack VP agent attempts to attack CP prey agent
func (vp *VisualPredator) Attack(prey *ColourPolymorphicPrey, ctxt Context) bool {
	if prey == nil {
		return false
	}
	α := rand.Float64()
	if α > (1 - ctxt.VpAttackChance) {
		vp.attackSuccess = true
		vp.colourImprinting(prey.colouration, ctxt.VpCaf)
		c := vp.ετ
		𝒇 := visualSignalStrength(c)
		𝛘 := colour.RGBDistance(vp.τ, prey.colouration)
		Vg := 𝒇(𝛘) * ctxt.Vbg
		vp.hunger -= int(Vg)
		if vp.hunger < 0 {
			vp.hunger = 0
		}
		prey.lifespan = 0 //	i.e. prey agent is flagged for removal at the beginning of next turn and will not be drawn again.
		if ctxt.Vmε > vp.ετ {
			vp.ετ++
		}
		if vp.𝛄 > ctxt.Vb𝛄 {
			vp.𝛄 *= (1 - ctxt.V𝛄Bump) //	returning towards context-defined value
		}
		return vp.attackSuccess
	}
	// FAILURE
	vp.attackSuccess = false
	// MAYBE THIS SHOULD BE DETERMINED IF STARVING OR NOT?
	if vp.ετ > ctxt.Vbε {
		vp.ετ-- //	decrease target colour signal strength factor
	}
	return vp.attackSuccess
}

// Intercept attempts to turn and move towards target position (as much as vp is able)
func (vp *VisualPredator) Intercept(target geometry.Vector) (bool, error) {
	dist, _ := geometry.VectorDistance(vp.pos, target)
	Ψ, err := geometry.AngleToIntercept(vp.pos, vp.𝚯, target)
	if dist < vp.movS {
		vp.pos = target
		vp.Turn(calc.ClampFloatIn(Ψ, -vp.tr, vp.tr))
		return true, err
	}
	vp.Turn(calc.ClampFloatIn(Ψ, -vp.tr, vp.tr))
	// vp.Turn(Ψ)
	vp.Move()
	return false, err
}

// MateSearch searches species population for sexual coupling
func (vp *VisualPredator) MateSearch(neighbours []VisualPredator, me int, errCh chan<- error) *VisualPredator {
	if len(neighbours) == 0 {
		return nil
	}

	var searchSet []proxVP
	f := func(u geometry.Vector, errCh chan<- error) func(geometry.Vector) float64 {
		return func(v geometry.Vector) float64 {
			δ, err := geometry.VectorDistance(u, v)
			errCh <- err
			return δ
		}
	}(vp.pos, errCh)

	for i := range neighbours {
		if i == me { //	SEXUAL not asexual reproduction! 😘
			continue
		}
		searchSet = append(searchSet, proxVP{f, &neighbours[i]})
	}

	if len(searchSet) == 0 {
		return nil
	}

	// fmt.Println()
	// for i := range searchSet {
	// 	fmt.Printf("%v\tδ=%v\t%v\t%p\n", i, searchSet[i].comp(searchSet[i].pos), searchSet[i].pos, searchSet[i].VisualPredator)
	// }

	sort.Sort(byProximityVp(searchSet))

	// fmt.Println()
	// for i := range searchSet {
	// 	fmt.Printf("%v\tδ=%v\t%v\t%p\n", i, searchSet[i].comp(searchSet[i].pos), searchSet[i].pos, searchSet[i].VisualPredator)
	// }
	target := searchSet[0].pos // guaranteed to exist by test on test of searchSet length above

	// fmt.Printf("Before Intercept:\n%v\t%p\n", vp.pos, vp)

	inRange, err := vp.Intercept(target)
	errCh <- err
	if inRange {
		// fmt.Printf("After Intercept:\n%v\t%p\n", vp.pos, vp)
		return searchSet[0].VisualPredator
	}

	return nil
}

// animal-agent Mortal interface methods:

// Age the vp agent
func (vp *VisualPredator) Age(ctxt Context, popSize int) string {
	vp.attackSuccess = false
	vp.fertility++
	vp.hunger++

	if ctxt.VpStarvation {
		if vp.hunger > ctxt.VpPanicPoint { //	if the agent is getting desperate, it lowers its focus and has to start looking harder.
			vp.𝛄 *= ctxt.V𝛄Bump // (default is 1.1 == a 10% bump)
			if (vp.hunger%5 == 0) && (vp.ετ > ctxt.Vbε) {
				vp.ετ-- //	the energy gain from attack success reduces because it costs more energy to look harder!
			}
		}
	}

	if ctxt.VpAgeing {
		vp.lifespan--
	}

	return vp.jump(ctxt, popSize)
}

func (vp *VisualPredator) jump(ctxt Context, popSize int) (jump string) {
	switch {
	case vp.lifespan <= 0:
		jump = "DEATH"
	case vp.fertility == 0:
		vp.gravid = false
		jump = "SPAWN"
	case ctxt.VpStarvation && (vp.hunger > ctxt.VpStarvationPoint):
		jump = "DEATH"
	case (popSize < ctxt.VpPopulationCap) && (vp.fertility > ctxt.VpSexualRequirement/2) && (vp.hunger < ctxt.VpSexualRequirement):
		jump = "FERTILE"
	default:
		jump = "PREY SEARCH"
	}
	return
}

// Copulation for sexual reproduction between Visual Predator agents
func (vp *VisualPredator) Copulation(mate *VisualPredator, ctxt Context) bool {
	if mate == nil {
		return false
	}
	if mate.fertility < ctxt.VpSexualRequirement {
		return false
	}
	ω := rand.Float64()
	mate.fertility = 1 // it takes two to tango, buddy!
	if ω <= ctxt.VpReproductionChance {
		vp.gravid = true
		vp.fertility = -ctxt.VpGestation
		return true
	}
	vp.fertility = 1
	return false
}

// Birth spawns Visual Predator children
func (vp *VisualPredator) Birth(ctxt Context, start int, mt int) []VisualPredator {
	n := 1
	if ctxt.VpSpawnSize > 1 {
		n = rand.Intn(ctxt.VpSpawnSize) + 1
	}
	// func vpSpawn(size int, start int, mt int, parent VisualPredator, context Context)
	timestamp := fmt.Sprintf("%s", time.Now())
	progeny := vpSpawn(n, start, mt, *vp, ctxt, timestamp)
	vp.hunger++
	vp.gravid = false
	return progeny
}
