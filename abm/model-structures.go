package abm

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/benjamin-rood/abm-cp/colour"
	"github.com/benjamin-rood/goio"
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
	CppSpawnSize          int       `json:"abm-cpp-spawn-size"`          // possible number of progeny = [1, max]
	CppMutationFactor     float64   `json:"abm-cpp-mf"`                  //	mutation factor
	VpPopulationStart     int       `json:"abm-vp-pop-start"`            //	starting VP agent population size
	VpPopulationCap       int       `json:"abm-vp-pop-cap"`              //
	VpAgeing              bool      `json:"abm-vp-ageing"`               //
	VpLifespan            int       `json:"abm-vp-lifespan"`             //	Visual Predator lifespan
	VpStarvationPoint     int       `json:"abm-vp-starvation-point"`     //
	VpGestation           int       `json:"abm-vp-gestation"`            //	Visual Predator gestation period
	VpSexualRequirement   int       `json:"abm-vp-sex-req"`              //
	VpMovS                float64   `json:"abm-vp-speed"`                // Visual Predator speed
	VpMovA                float64   `json:"abm-vp-acceleration"`         // Visual Predator acceleration
	VpTurn                float64   `json:"abm-vp-turn"`                 //	Visual Predator turn rate / range (in radians)
	Vsr                   float64   `json:"abm-vp-vsr"`                  //	VP agent visual search range
	Vγ                    float64   `json:"abm-vp-visual-acuity"`
	VγBump                float64   `json:"abm-vp-visual-acuity-bump"`
	VpReproductionChance  float64   `json:"abm-vp-reproduction-chance"` //	chance of VP copulation success.
	VpSpawnSize           int       `json:"abm-vp-spawn-size"`
	VpSearchChance        float64   `json:"abm-vp-vsr-chance"`
	VpAttackChance        float64   `json:"abm-vp-attack-chance"`
	VpColImprintFactor    float64   `json:"abm-vp-col-imprinting"`
	Starvation            bool      `json:"abm-starvation"`
	RandomAges            bool      `json:"abm-random-ages"`
	RNGRandomSeed         bool      `json:"abm-rng-random-seed"` //	flag for using server-set random seed val.
	RNGSeedVal            int64     `json:"abm-rng-seedval"`     //	RNG seed value
	Fuzzy                 float64   `json:"abm-fuzziness"`
	Logging               bool      `json:"abm-logging-flag"`  //	log abm on/off
	LogFreq               int       `json:"abm-log-frequency"` // how many turns between writing log files.
	CustomLogPath         bool      `json:"abm-custom-log-filepath"`
	LogPath               string    `json:"abm-log-filepath"`
	Visualise             bool      `json:"abm-visualise-flag"`     //	Visualise on/off
	Duration              int       `json:"abm-duration"`           // fixed abm running length.
	SessionIdentifier     string    `json:"abm-session-identifier"` //	user-friendly string (from client) to identify session
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

// Reset 's the timeframe to 00:00:00
func (t *Timeframe) Reset() {
	t.Turn, t.Phase, t.Action = 0, 0, 0
}

// Model acts as the working instance of the 'game'
type Model struct {
	sessionName string
	timestamp   string
	running     bool
	Dead        bool
	Timeframe
	Environment
	Context
	PopulationCPP
	PopulationVP
	numCppCreated int
	numVpCreated  int
	recordCPP     map[string]ColourPolymorphicPrey
	rcm           sync.RWMutex
	recordVP      map[string]VisualPredator
	rvm           sync.RWMutex
	Om            chan goio.OutMsg
	Im            chan goio.InMsg
	e             chan error
	Quit          chan struct{} //	instance signaling
	r             chan struct{} //	run signalling
}

// AgentDescription used to aid for logging / debugging
type AgentDescription struct {
	AgentType  string `json:"agent-type"`
	AgentNum   int    `json:"agent-num"`
	ParentUUID string `json:"parent"`
	CreatedMT  int    `json:"creation-turn"`
	CreatedAT  string `json:"creation-date"`
}

// NewModel is a constructor for initialising a Model instance
func NewModel() (m *Model) {
	m.sessionName = m.SessionIdentifier
	m.timestamp = fmt.Sprintf("%s", time.Now())
	m.running = false
	m.Timeframe = Timeframe{}
	m.Environment = DefaultEnvironment
	m.Context = DefaultContext
	m.recordCPP = make(map[string]ColourPolymorphicPrey)
	m.recordVP = make(map[string]VisualPredator)
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

func uuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (m *Model) copyCppRecord() map[string]ColourPolymorphicPrey {
	var record = make(map[string]ColourPolymorphicPrey)
	m.rcm.RLock()
	for k, v := range m.recordCPP {
		record[k] = v
	}
	m.rcm.RUnlock()
	return record
}

func (m *Model) copyVpRecord() map[string]VisualPredator {
	var record = make(map[string]VisualPredator)
	m.rvm.RLock()
	for k, v := range m.recordVP {
		record[k] = v
	}
	m.rvm.RUnlock()
	return record
}
