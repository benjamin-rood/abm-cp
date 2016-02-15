package abm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/benjamin-rood/abm-cp/calc"
	"github.com/benjamin-rood/abm-cp/colour"
	"github.com/benjamin-rood/abm-cp/geometry"
	"github.com/benjamin-rood/abm-cp/render"
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
	γ             float64    //	visual seach (colour) bias
	τ             colour.RGB //	imprinted target / colour specialisation value
	ετ            float64    //	imprinting / colour specialisation strength
}

// UUID is just a getter method for the unexported uuid field, which absolutely must not change after agent creation.
func (vp *VisualPredator) UUID() string {
	return vp.uuid
}

// MarshalJSON implements json.Marshaler interface for VisualPredator object
func (vp VisualPredator) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"description":             vp.description,
		"pos":                     vp.pos,
		"speed":                   vp.movS,
		"heading":                 vp.𝚯,
		"turn-rate":               vp.tr,
		"search-range":            vp.vsr,
		"lifespan":                vp.lifespan,
		"hunger":                  vp.hunger,
		"attack-success":          vp.attackSuccess,
		"fertility":               vp.fertility,
		"γ":                       vp.γ,
		"gravid":                  vp.gravid,
		"colour-target-value":     vp.τ,
		"colour-imprint-strength": vp.ετ,
	})
}

// GetDrawInfo exports the data set needed for agent visualisation.
func (vp *VisualPredator) GetDrawInfo() (ar render.AgentRender) {
	ar.Type = "vp"
	ar.X = vp.pos[x]
	ar.Y = vp.pos[y]
	ar.Heading = vp.𝚯
	if vp.attackSuccess {
		// inv := vp.τ.Invert()
		// ar.Colour = inv.To256()
		ar.Colour = colour.RGB256{Red: 0, Green: 0, Blue: 0} // blink black on successful attack!
	} else {
		ar.Colour = vp.τ.To256()
	}
	return
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
		agent.γ = context.Vbγ //	baseline acuity level
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
		agent.𝚯 = parent.𝚯
		agent.dir = parent.dir
		agent.tr = parent.tr
		agent.vsr = parent.vsr
		agent.hunger = context.VpSexualRequirement + 1
		agent.fertility = 1
		agent.gravid = false
		agent.τ = colour.RandRGB()
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
func (vp *VisualPredator) PreySearch(prey []ColourPolymorphicPrey, searchChance float64) (*ColourPolymorphicPrey, float64, error) {
	_ = "breakpoint" // godebug
	c := vp.ετ
	var f = visualSignalStrength(c)
	var 𝛘 float64 // colour sorting value - colour distance/difference between vp.imprimt and cpp.colouration
	var δ float64 // position sorting value - vector distance between vp.pos and cpp.pos
	var err error
	var searchSet []visualRecognition
	for i := range prey { //	exhaustive search 😱
		δ, err = geometry.VectorDistance(vp.pos, prey[i].pos)
		if δ <= vp.vsr { // ∴ only include the prey agent for considertion if within visual range
			𝛘 = colour.RGBDistance(vp.τ, prey[i].colouration)
			fmt.Printf("%v\t%v\t%v\t%v\t%p\n", i, δ, 𝛘, c, &prey[i])
			if 𝛘 < vp.γ { // i.e. if and only if colour distance < predator's visual search bias
				a := visualRecognition{δ, 𝛘, f, c, &prey[i]}
				searchSet = append(searchSet, a)
			}
		}
	}

	for i := range searchSet {
		fmt.Printf("%v\t%v\t%v\t%v\t%p\t%v\t%v\n", i, searchSet[i].δ, searchSet[i].𝛘, c, searchSet[i].ColourPolymorphicPrey, f(searchSet[i].𝛘), f(searchSet[i].𝛘)-searchSet[i].δ)
	}

	sort.Sort(byOptimalAttackVector(searchSet)) //	sort by f(x) - distance

	for i := range searchSet {
		fmt.Printf("%v\t%v\t%v\t%v\t%p\t%v\t%v\n", i, searchSet[i].δ, searchSet[i].𝛘, c, searchSet[i].ColourPolymorphicPrey, f(searchSet[i].𝛘), f(searchSet[i].𝛘)-searchSet[i].δ)
	}

	// search within biased and reduced set
	for i, p := range searchSet {
		return &(*searchSet[i].ColourPolymorphicPrey), p.δ, err
		// if ((1.0 - p.𝛘) * (1.0 - p.δ)) > (1.0 - searchChance) {
		// 	return &(*searchSet[i].ColourPolymorphicPrey), p.δ, err
		// }
	}
	return nil, 0, err
}

// Intercept attempts to turn and move towards target position (as much as vp is able)
// note: generalised to a position vector and distance measurement so that Intercept can be used for any type of targeting.
func (vp *VisualPredator) Intercept(vx geometry.Vector, dist float64) (bool, error) {
	var inRange bool
	Ψ, err := geometry.AngleToIntercept(vp.pos, vp.𝚯, vx)
	if dist < vp.movS {
		inRange = true
		vp.pos = vx
		vp.Turn(calc.ClampFloatIn(Ψ, -vp.tr, vp.tr))
		return inRange, err
	}
	vp.Turn(calc.ClampFloatIn(Ψ, -vp.tr, vp.tr))
	// vp.Turn(Ψ)
	vp.Move()
	return inRange, err
}

// MateSearch searches species population for sexual coupling
func (vp *VisualPredator) MateSearch(predators []VisualPredator, me int) (*VisualPredator, error) {
	min := math.MaxFloat64
	var closest *VisualPredator
	var err error
	var dist float64
	for i := range predators {
		if i == me {
			continue
		}
		dist, err = geometry.VectorDistance(vp.pos, predators[i].pos)
		if dist < min {
			min = dist
			closest = &predators[i]
		}
	}
	return closest, err
}

// Attack VP agent attempts to attack CP prey agent
func (vp *VisualPredator) Attack(prey *ColourPolymorphicPrey, vpAttackChance float64, caf float64, bg float64, bγ float64, bε float64) {
	if prey == nil {
		return
	}
	_ = "breakpoint" // godebug
	vpAttackChance = 1 - vpAttackChance
	α := rand.Float64()
	if α > vpAttackChance {
		vp.colourImprinting(prey.colouration, caf)
		f := visualSignalStrength(math.Pow(2, vp.ετ))
		Vg := f(colour.RGBDistance(vp.τ, prey.colouration)) * bg
		vp.hunger -= int(Vg)
		prey.lifespan = 0 //	i.e. prey agent is flagged for removal at the beginning of next turn and will not be drawn again.
		vp.attackSuccess = true
		vp.ετ++
		vp.γ = bγ //	resetting to context-defined value
		fmt.Println(vp.String())
		return
	}
	if vp.ετ > bε {
		vp.ετ-- //	decrease target colour signal strength factor
	}
}

// colourImprinting updates VP colour / visual recognition bias
// Uses a bias / weighting value, 𝜎 (sigma) to control the degree of
// adaptation VP will make to differences in 'eaten' CPP colours.
func (vp *VisualPredator) colourImprinting(target colour.RGB, 𝜎 float64) {
	𝚫red := (vp.τ.Red - target.Red) * 𝜎
	𝚫green := (vp.τ.Green - target.Green) * 𝜎
	𝚫blue := (vp.τ.Blue - target.Blue) * 𝜎
	vp.τ.Red = vp.τ.Red - 𝚫red
	vp.τ.Green = vp.τ.Green - 𝚫green
	vp.τ.Blue = vp.τ.Blue - 𝚫blue
}

// animal-agent Mortal interface methods:

// Age the vp agent
func (vp *VisualPredator) Age(ctxt Context, popSize int) string {
	_ = "breakpoint" // godebug
	vp.attackSuccess = false
	vp.fertility++
	vp.hunger++

	if ctxt.Starvation {
		if vp.hunger > ctxt.VpPanicPoint { //	if the agent is getting desperate, it lowers its focus and has to start looking harder.
			vp.γ *= ctxt.VγBump // (default is 1.1 == a 10% bump)
			if vp.ετ > ctxt.Vbε {
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
	_ = "breakpoint" // godebug
	switch {
	case vp.lifespan <= 0:
		jump = "DEATH"
	case vp.fertility == 0:
		vp.gravid = false
		jump = "SPAWN"
	case ctxt.Starvation && (vp.hunger > ctxt.VpStarvationPoint):
		jump = "DEATH"
	case (popSize < ctxt.VpPopulationCap) && (vp.fertility > 0) && (vp.hunger < ctxt.VpSexualRequirement):
		jump = "FERTILE"
	default:
		jump = "PREY SEARCH"
	}
	return
}

// Copulation for sexual reproduction between Visual Predator agents
func (vp *VisualPredator) Copulation(mate *VisualPredator, chance float64, gestation int, sexualCost int) bool {
	if mate == nil {
		return false
	}
	if mate.fertility < sexualCost {
		return false
	}
	ω := rand.Float64()
	mate.fertility = -sexualCost // it takes two to tango, buddy!
	if ω <= chance {
		vp.gravid = true
		vp.fertility = -gestation
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

// String returns a clear textual presentation the internal values of the VP agent
func (vp *VisualPredator) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("pos=(%v,%v)\n", vp.pos[x], vp.pos[y]))
	buffer.WriteString(fmt.Sprintf("movS=%v\n", vp.movS))
	buffer.WriteString(fmt.Sprintf("movA=%v\n", vp.movA))
	buffer.WriteString(fmt.Sprintf("𝚯=%v\n", vp.𝚯))
	buffer.WriteString(fmt.Sprintf("dir=(%v,%v)\n", vp.dir[x], vp.dir[y]))
	buffer.WriteString(fmt.Sprintf("tr=%v\n", vp.tr))
	buffer.WriteString(fmt.Sprintf("Vsr=%v\n", vp.vsr))
	buffer.WriteString(fmt.Sprintf("lifespan=%v\n", vp.lifespan))
	buffer.WriteString(fmt.Sprintf("hunger=%v\n", vp.hunger))
	buffer.WriteString(fmt.Sprintf("fertility=%v\n", vp.fertility))
	buffer.WriteString(fmt.Sprintf("gravid=%v\n", vp.gravid))
	buffer.WriteString(fmt.Sprintf("τ=%v\n", vp.τ))
	return buffer.String()
}
