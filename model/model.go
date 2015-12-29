package model

import "github.com/benjamin-rood/abm-colour-polymorphism/colour"

// Model acts as the working instance of the 'game'
type Model struct {
	Timeframe
	Environment
	CPP cppPopulation
	VP  vpPopulation
}

type cppPopulation struct {
	Pop        []ColourPolymorhicPrey
	Definition []string //	lists agent interfaces which define the behaviour of this type
}

type vpPopulation struct {
	Pop        []VisualPredator
	Definition []string //	lists agent interfaces which define the behaviour of this type
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
