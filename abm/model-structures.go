package abm

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/benjamin-rood/abm-cp/colour"
	"github.com/benjamin-rood/abm-cp/render"
	"github.com/benjamin-rood/gobr"
)

/*
Environment specifies the boundary / dimensions of the working model. They
extend in both positive and negative directions, oriented at the center. Setting
any field (eg. zBounds) to zero will reduce the dimensionality of the model. For
most cases, a 2D environment will be sufficient.
In the future it may include some environmental factors etc.
*/
type Environment struct {
	Bounds         []float64  `json:"abm-environment-bounds"` // d value for each axis
	Dimensionality int        `json:"abm-environment-dimensionality"`
	BG             colour.RGB `json:"abm-environment-background"`
}

// Context contains the local model context;
type Context struct {
	Environment           `json:"abm-environment-bounds"` // d value for each axis
	CppPopulationStart    int                             `json:"abm-cpp-pop-start"`                // starting CPP agent population size
	CppPopulationCap      int                             `json:"abm-cpp-pop-cap"`                  //
	CppAgeing             bool                            `json:"abm-cpp-ageing"`                   //
	CppLifespan           int                             `json:"abm-cpp-lifespan"`                 //	CPP agent lifespan
	CppS                  float64                         `json:"abm-cpp-speed"`                    // CPP agent speed
	CppA                  float64                         `json:"abm-cpp-acceleration"`             // CPP agent acceleration
	CppTurn               float64                         `json:"abm-cpp-turn"`                     //	CPP agent turn rate / range (in radians)
	CppSr                 float64                         `json:"abm-cpp-sr"`                       // CPP agent search range for mating
	CppGestation          int                             `json:"abm-cpp-gestation"`                //	CPP gestation period
	CppSexualCost         int                             `json:"abm-cpp-sexual-cost"`              //	CPP sexual rest cost
	CppReproductionChance float64                         `json:"abm-cpp-reproduction-chance"`      //	chance of CPP copulation success.
	CppSpawnSize          int                             `json:"abm-cpp-spawn-size"`               // possible number of progeny = [1, max]
	CppMutationFactor     float64                         `json:"abm-cpp-mf"`                       //	mutation factor
	VpPopulationStart     int                             `json:"abm-vp-pop-start"`                 //	starting VP agent population size
	VpPopulationCap       int                             `json:"abm-vp-pop-cap"`                   //
	VpAgeing              bool                            `json:"abm-vp-ageing"`                    //
	VpLifespan            int                             `json:"abm-vp-lifespan"`                  //	Visual Predator lifespan
	VpStarvationPoint     int                             `json:"abm-vp-starvation-point"`          //
	VpPanicPoint          int                             `json:"abm-vp-panic-point"`               //
	VpGestation           int                             `json:"abm-vp-gestation"`                 //	Visual Predator gestation period
	VpSexualRequirement   int                             `json:"abm-vp-sex-req"`                   //
	VpMovS                float64                         `json:"abm-vp-speed"`                     // Visual Predator speed
	VpMovA                float64                         `json:"abm-vp-acceleration"`              // Visual Predator acceleration
	VpTurn                float64                         `json:"abm-vp-turn"`                      //	Visual Predator turn rate / range (in radians)
	Vsr                   float64                         `json:"abm-vp-vsr"`                       //	VP agent visual search range
	Vb𝛄                   float64                         `json:"abm-vp-visual-acuity"`             //
	V𝛄Bump                float64                         `json:"abm-vp-visual-acuity-bump"`        //
	Vbε                   float64                         `json:"abm-vp-baseline-col-sig-strength"` // 	baseline colour signal strength factor
	Vmε                   float64                         `json:"abm-vp-max-col-sig-strength"`      // 	max limit colour signal strength factor
	VpReproductionChance  float64                         `json:"abm-vp-reproduction-chance"`       //	chance of VP copulation success.
	VpSpawnSize           int                             `json:"abm-vp-spawn-size"`                //
	VpSearchChance        float64                         `json:"abm-vp-vsr-chance"`                //
	VpAttackChance        float64                         `json:"abm-vp-attack-chance"`             //
	Vbg                   float64                         `json:"abm-vp-baseline-attack-gain"`      //
	VpCaf                 float64                         `json:"abm-vp-col-adaptation-factor"`     //
	VpStarvation          bool                            `json:"abm-vp-starvation"`                //
	RandomAges            bool                            `json:"abm-random-ages"`                  //
	RNGRandomSeed         bool                            `json:"abm-rng-random-seed"`              //	flag for using server-set random seed val.
	RNGSeedVal            int64                           `json:"abm-rng-seedval"`                  //	RNG seed value
	Fuzzy                 float64                         `json:"abm-rng-fuzziness"`                //
	Logging               bool                            `json:"abm-logging-flag"`                 //	log abm on/off
	LogFreq               int                             `json:"abm-log-frequency"`                // how many turns between writing log files.
	UseCustomLogPath      bool                            `json:"abm-use-custom-log-filepath"`      //
	CustomLogPath         string                          `json:"abm-custom-log-filepath"`          //
	LogPath               string                          `json:"abm-log-filepath"`                 //
	Visualise             bool                            `json:"abm-visualise-flag"`               //	Visualise on/off
	LimitDuration         bool                            `json:"abm-limit-duration"`               //
	FixedDuration         int                             `json:"abm-fixed-duration"`               // fixed abm running length.
	SessionIdentifier     string                          `json:"abm-session-identifier"`           //	user-friendly string (from client) to identify session
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
	timestamp string
	running   bool
	Dead      bool
	Timeframe
	Environment
	Context
	PopulationCPP
	PopulationVP
	numCppCreated int
	numVpCreated  int
	recordCPP     map[string]ColourPolymorphicPrey
	rcppRW        sync.RWMutex
	recordVP      map[string]VisualPredator
	rvpRW         sync.RWMutex
	Om            chan gobr.OutMsg
	Im            chan gobr.InMsg
	e             chan error              //	error message channel (general)
	Quit          chan struct{}           //	instance signaling
	rc            chan struct{}           //	run signalling
	render        chan render.AgentRender //	visualisation message channel
	turnSignal    *gobr.SignalHub         //	turn signalling and broadcasting
}

// AgentDescription used to aid for logging / debugging - used at time of agent creation
type AgentDescription struct {
	AgentType  string `json:"agent-type"`
	AgentNum   int    `json:"agent-num"`
	ParentUUID string `json:"parent"`
	CreatedMT  int    `json:"creation-turn"`
	CreatedAT  string `json:"creation-date"`
}

// NewModel is a constructor for initialising a Model instance
func NewModel() *Model {
	_ = "breakpoint" // godebug
	m := Model{}
	m.timestamp = fmt.Sprintf("%s", time.Now())
	m.running = false
	m.Timeframe = Timeframe{}
	m.Environment = DefaultEnvironment
	m.Context = DefaultContext
	m.LogPath = path.Join(os.Getenv("HOME")+os.Getenv("HOMEPATH"), abmlogPath, m.SessionIdentifier, m.timestamp)
	m.recordCPP = make(map[string]ColourPolymorphicPrey)
	m.recordVP = make(map[string]VisualPredator)
	m.Om = make(chan gobr.OutMsg)
	m.Im = make(chan gobr.InMsg)
	m.e = make(chan error)
	m.Quit = make(chan struct{})
	m.rc = make(chan struct{})
	m.render = make(chan render.AgentRender)
	_ = "breakpoint" // godebug
	m.turnSignal = gobr.NewSignalHub()
	return &m
}

// PopLog prints the current time and populations
// shit version
func (m *Model) PopLog() {
	log.Printf("%04dT : %04dP : %04dA\n", m.Turn, m.Phase, m.Action)
	log.Printf("cpp population size = %v\n", len(m.PopCPP))
	log.Printf("vp population size = %v\n", len(m.PopVP))
}

func uuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (m *Model) cppRecordCopy() map[string]ColourPolymorphicPrey {
	defer m.rcppRW.RUnlock()
	m.rcppRW.RLock()
	var record = make(map[string]ColourPolymorphicPrey)
	for k, v := range m.recordCPP {
		record[k] = v
	}
	return record
}

func (m *Model) cppRecordAssignValue(key string, value ColourPolymorphicPrey) error {
	defer m.rcppRW.Unlock()
	m.rcppRW.Lock()
	m.recordCPP[key] = value
	return nil
}

func (m *Model) vpRecordCopy() map[string]VisualPredator {
	defer m.rvpRW.RUnlock()
	m.rvpRW.RLock()
	var record = make(map[string]VisualPredator)
	for k, v := range m.recordVP {
		record[k] = v
	}
	return record
}

func (m *Model) vpRecordAssignValue(key string, value VisualPredator) error {
	defer m.rvpRW.Unlock()
	m.rvpRW.Lock()
	m.recordVP[key] = value
	return nil
}
