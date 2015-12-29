package model

import (
	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/geometry"
)

// Model acts as the working instance of the 'game'
type Model struct {
	Timeframe
	Environment
	CppPopulation []ColourPolymorhicPrey
	VpPopulation  []VisualPredator
}

/*
Timeframe holds the model's representation of the time metrics.
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
type Timeframe struct {
	Turn   int
	Phase  int
	Action int
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
In the future it may include some environmental factors etc.
*/
type Environment struct {
	bounds [][2]int //	from -d to d for each axis.
	bg     colour.RGB
}

// Context contains the local model context;
type Context struct {
	𝐄                  Environment
	time               Timeframe
	cppPopulation      uint    // CPP agent population size
	vpPopulation       uint    //	VP agent population size
	vsr                float64 //	visual search range
	γ                  float64 //	visual acuity in environments
	vpLifespan         uint    //	Visual Predator lifespan
	cppLifespan        uint    //	Colour Polymorphic Prey lifespan
	φ                  uint    //	CPP incubation cost
	ȣ                  uint    //	CPP sexual rest cost
	vpAgeing           bool
	cppAgeing          bool
	cppReproduceChance float64
	vpReproduceChance  float64
	vsrSearchChance    float64
	vpAttackChance     float64
}

// ColourPolymorhicPrey – Prey agent type for Predator-Prey ABM
type ColourPolymorhicPrey struct {
	populationIndex uint            //	index to the master population array.
	pos             geometry.Vector //	position in the environment
	movS            float64         //	speed
	movA            float64         //	acceleration
	dir             geometry.Vector //	must be implemented as a unit vector
	dir𝚯            float64         //	 heading angle
	hunger          uint            //	counter for interval between needing food
	fertility       uint            //	counter for interval between birth and sex
	gravid          bool            //	i.e. pregnant
	colouration     colour.RGB      //	colour
	𝛘               float64         //	 colour sorting value - colour distance/difference between vp.imprimt and cpp.colouration
	ϸ               float64         //  position sorting value - vector distance between vp.pos and cpp.pos
}

// ProximitySort implements sort.Interface for []ColourPolymorhicPrey
// based on δ field.
type ProximitySort []ColourPolymorhicPrey

func (ps ProximitySort) Len() int           { return len(ps) }
func (ps ProximitySort) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
func (ps ProximitySort) Less(i, j int) bool { return ps[i].ϸ < ps[j].ϸ }

// VisualSort implements sort.Interface for []ColourPolymorhicPrey
// based on 𝛘 field – to assert visual bias of a VisualPredator based on it's colour imprinting value.
type VisualSort []ColourPolymorhicPrey

func (vs VisualSort) Len() int           { return len(vs) }
func (vs VisualSort) Swap(i, j int)      { vs[i], vs[j] = vs[j], vs[i] }
func (vs VisualSort) Less(i, j int) bool { return vs[i].𝛘 < vs[j].𝛘 }

// VisualPredator - Predator agent type for Predator-Prey ABM
type VisualPredator struct {
	populationIndex uint            //	index to the master population array.
	pos             geometry.Vector //	position in the environment
	movS            float64         //	speed
	movA            float64         //	acceleration
	dir             geometry.Vector //	must be implemented as a unit vector
	dir𝚯            float64         //	 heading angle
	hunger          uint            //	counter for interval between needing food
	fertility       uint            //	counter for interval between birth and sex
	gravid          bool            //	i.e. pregnant
	vsr             float64         //	visual search range
	γ               float64         //	visual acuity (initially, use 1.0)
	colImprint      colour.RGB
}
