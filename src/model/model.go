package main

import "errors"

/*
MTime holds the model's representation of the time metrics.
Turn – The cycle length for all agents ∈ 𝐄 to perform 1 (and only 1) Action.
Phase – Division of a Turn, between agent sets, environmental effects/factors,
				and updates to populations and model conditions (via external).
				One Phase is complete when all members of a set have performed an Action
				or all requirements for the model's continuation have been fulfilled.
Action – An individual 'step' in the model. All Actions have a cost:
				the period (number of turns) before that specific Action can be
				performed again. For most actions this is zero.
				Some Actions could also *stop* any other behaviour by that agent
				for a period.
*/
type MTime struct {
	turn   int
	phase  int
	action int
}

const (
	x = iota
	y
	z
)

/*
Environment specifies the boundary / dimensions of the working model. They
extend in both positive and negative directions, oriented at the center. Setting
any field (eg. zBounds) to zero will reduce the dimensionality of the model. For
most cases, a 2D environment will be sufficient.
*/
type Environment struct {
	xBounds float64
	yBounds float64
	zBounds float64
}

// ColourPolymorhicPrey – Prey agent type for Predator-Prey ABM
type ColourPolymorhicPrey struct {
	populationIndex int
	pos             Vector
	movS            float64
	movA            float64
	heading         float64
	direction       Vector
	lifetime        int32
	hunger          int32
	gravid          bool //	i.e. pregnant
	colour          ColRGB
	𝛘               float64 //	colour sorting value
	δ               float64 // position sorting value
}

// ProximitySort implements sort.Interface for []ColourPolymorhicPrey
// based on δ field.
type ProximitySort []ColourPolymorhicPrey

func (ps ProximitySort) Len() int           { return len(ps) }
func (ps ProximitySort) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
func (ps ProximitySort) Less(i, j int) bool { return ps[i].δ < ps[j].δ }

// VisualSort implements sort.Interface for []ColourPolymorhicPrey
// based on 𝛘 field.
type VisualSort []ColourPolymorhicPrey

func (vs VisualSort) Len() int           { return len(vs) }
func (vs VisualSort) Swap(i, j int)      { vs[i], vs[j] = vs[j], vs[i] }
func (vs VisualSort) Less(i, j int) bool { return vs[i].𝛘 < vs[j].𝛘 }

// VisualPredator - Predator agent type for Predator-Prey ABM
type VisualPredator struct {
	populationIndex int
	pos             Vector
	movS            float64
	movA            float64
	heading         float64 //	angle of direction relative to origin ∈ 𝐄
	direction       Vector  //	unit vector for
	lifetime        int32   //	counter for number of turns agent exists in the model
	hunger          int32   //	measurement reflecting
	fertility       int32   //	interval measurement between birth and sex
	gravid          bool    //	i.e. pregnant
	visRange        float64
	visAcuity       float64
	colImprint      ColRGB
}

// AgentActions interface for general agent behaviours
type AgentActions interface {
	Turn(𝝧 float64)
	Move()
	Death()
}

// VPbehaviour – set of actions only VisualPredator agents will perform.
type VPbehaviour interface {
	VisualSearch([]ColourPolymorhicPrey, float64) (*ColourPolymorhicPrey, error)
	// ColourImprinting updates VP colour / visual recognition bias
	ColourImprinting(ColRGB, float64) error
	Eat(*ColourPolymorhicPrey)
}

// VisualSearch tries to 'recognise' a nearby prey agent to attack.
func (vp *VisualPredator) VisualSearch(population []ColourPolymorhicPrey, vsrSearchChance float64) (*ColourPolymorhicPrey, error) {
	for i := range population {
		population[i].𝛘 = ColourDistance(vp.colImprint, population[i].colour)
	}

	population = VisualSort(population)

	for i := range population {
		distanceToTarget, err := VectorDistance(vp.pos, population[i].pos)
		if err != nil {
			return nil, err
		}
		if distanceToTarget > vp.visRange {
			return nil, errors.New("VisualSearch failed")
		}
		if (distanceToTarget * vp.visAcuity * population[i].𝛘) > vsrSearchChance {
			return &population[i], nil
		}
	}

	return nil, errors.New("VisualSearch failed")
}

// ColourImprinting updates VP colour / visual recognition bias
func (vp *VisualPredator) ColourImprinting(target ColRGB, colourWeighting float64) error {
	𝚫red := byte(float64(vp.colImprint.red-target.red) * colourWeighting)
	𝚫green := byte(float64(vp.colImprint.green-target.green) * colourWeighting)
	𝚫blue := byte(float64(vp.colImprint.blue-target.blue) * colourWeighting)
	vp.colImprint = ColRGB{
		vp.colImprint.red - 𝚫red,
		vp.colImprint.green - 𝚫green,
		vp.colImprint.blue - 𝚫blue}
	return nil
}

// Eat absorbs energy from ColourPolymorhicPrey
func (vp *VisualPredator) Eat(cpp *ColourPolymorhicPrey) {

}
