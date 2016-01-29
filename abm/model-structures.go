package abm

import (
	"log"
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/goio"
)

const (
	pause = 1 * time.Second
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
	Bounds                []float64 // d value for each axis
	CppPopulationStart    int       `json:"abm-cpp-pop-start"` // starting CPP agent population size
	CppPopulationCap      int       `json:"abm-cpp-pop-cap"`
	CppAgeing             bool      `json:"abm-cpp-ageing"`
	CppLifespan           int       `json:"abm-cpp-lifespan"` //	CPP agent lifespan
	CppS                  float64   `json:"abm-cpp-speed"`    // CPP agent speed
	CppA                  float64   // CPP agent acceleration
	CppTurn               float64   `json:"abm-cpp-turn"` //	CPP agent turn rate / range (in radians)
	CppSr                 float64   // CPP agent search range for mating
	CppGestation          int       `json:"abm-cpp-gestation"`           //	CPP gestation period
	CppSexualCost         int       `json:"abm-cpp-sexual-cost"`         //	CPP sexual rest cost
	CppReproductionChance float64   `json:"abm-cpp-reproduction-chance"` //	chance of CPP copulation success.
	CppSpawnSize          int       `json:"abm-cpp-spawn-size"`          // 	CPP max spawn size s.t. possible number of progeny = [1, max]
	VpPopulationStart     int       `json:"abm-vp-pop-start"`            //	starting VP agent population size
	VpPopulationCap       int       `json:"abm-vp-pop-cap"`
	VpAgeing              bool      `json:"abm-vp-ageing"`
	VpLifespan            int       `json:"abm-vp-lifespan"` //	Visual Predator lifespan
	VpStarvation          int       `json:"abm-vp-starvation"`
	VpMovS                float64   `json:"abm-vp-speed"` // Visual Predator speed
	VpMovA                float64   // Visual Predator acceleration
	VpTurn                float64   `json:"abm-vp-turn"` //	Visual Predator turn rate / range (in radians)
	Vsr                   float64   `json:"abm-vp-vsr"`  //	VP agent visual search range
	Vγ                    float64   //	visual acuity in environments
	VpReproductiveChance  float64   //	chance of VP copulation success.
	VpSearchChance        float64   `json:"abm-vp-vsr-chance"`
	VpAttackChance        float64   `json:"abm-vp-attack-chance"`
	VpColImprintFactor    float64   `json:"abm-vp-col-imprinting"`
	CppMutationFactor     float64   `json:"abm-cpp-mf"` //	mutation factor
	RandomAges            bool      `json:"abm-random-ages"`
	RNGRandomSeed         bool      `json:"abm-rng-random-seed"` //	flag for using server-set random seed val.
	RNGSeedVal            int64     `json:"abm-rng-seedval"`     //	RNG seed value
	Fuzzy                 float64
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
	Turn   uint64
	Phase  uint64
	Action uint64
}

// Reset 's the timeframe to 00:00:00
func (t *Timeframe) Reset() {
	t.Turn, t.Phase, t.Action = 0, 0, 0
}

// Model acts as the working instance of the 'game'
type Model struct {
	running bool
	Dead    bool
	Timeframe
	Environment
	Context
	PopulationCPP
	PopulationVP
	Om   chan goio.OutMsg
	Im   chan goio.InMsg
	e    chan error
	Quit chan struct{}
	r    chan struct{}
}

// NewModel is a constructor for initialising a Model instance
func NewModel() (m Model) {
	m.running = false
	m.Timeframe = Timeframe{}
	m.Environment = Environment{
		Bounds:         []float64{1.0, 1.0},
		Dimensionality: 2,
		BG:             colour.RGB{Red: 0.1, Green: 0.1, Blue: 0.1},
	}
	m.Context = DefaultContext
	m.Om = make(chan goio.OutMsg)
	m.Im = make(chan goio.InMsg)
	m.e = make(chan error)
	m.Quit = make(chan struct{})
	m.r = make(chan struct{})
	return
}

// Log prints the current state of time
// shit version
func (m *Model) Log() {
	log.Printf("%04dT : %04dP : %04dA\n", m.Turn, m.Phase, m.Action)
	log.Printf("cpp population size = %d\n", len(m.PopCPP))
	log.Printf("vp population size = %d\n", len(m.PopVP))
}
