package abm

import "github.com/benjamin-rood/abm-colour-polymorphism/colour"

// Model acts as the working instance of the 'game'
type Model struct {
	Timeframe
	Environment
	Context
	PopulationCPP
	PopulationVP
}

// PopulationCPP holds the agent population
type PopulationCPP struct {
	PopCPP        []ColourPolymorphicPrey
	DefinitionCPP []string //	lists agent interfaces which define the behaviour of this type
}

// PopulationVP holds the agent population
type PopulationVP struct {
	PopVP        []VisualPredator
	DefinitionVP []string //	lists agent interfaces which define the behaviour of this type
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
	Bounds         []float64 // d value for each axis
	Dimensionality int
	BG             colour.RGB
}

// Context contains the local model context;
type Context struct {
	// Type          string    `json:"type"` //	json flag for deserialisation
	Bounds        []float64 // d value for each axis
	CppPopulation int       // starting CPP agent population size
	VpPopulation  uint      //	starting VP agent population size
	VpAgeing      bool
	VpLifespan    int     //	Visual Predator lifespan
	VS            float64 // Visual Predator speed
	VA            float64 // Visual Predator acceleration
	Vτ            float64 //	Visual Predator turn rate / range (in radians)
	Vsr           float64 //	VP agent visual search range
	Vγ            float64 //	visual acuity in environments
	Vκ            float64 //	chance of VP copulation success.
	V𝛔            float64 // VsrSearchChance
	V𝛂            float64 // VpAttackChance
	CppAgeing     bool
	CppLifespan   int     //	CPP agent lifespan
	CppS          float64 // CPP agent speed
	CppA          float64 // CPP agent acceleration
	Cτ            float64 //	CPP agent turn rate / range (in radians)
	CppSr         float64 // CPP agent search range for mating
	RandomAges    bool
	Mf            float64 //	mutation factor
	Cφ            int     //	CPP incubation cost
	Cȣ            int     //	CPP sexual rest cost
	Cκ            float64 //	chance of CPP copulation success.
	Cβ            int     // 	CPP max spawn size (birth range)
}
