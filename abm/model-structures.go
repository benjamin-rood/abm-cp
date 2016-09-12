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

// Model acts as the working instance of the modelling session / 'game'
type Model struct {
	timestamp        string // instance inception time
	running          bool
	Timeframe        // embedded Model clock
	Environment      // embedded environment attributes
	ConditionParams  //	embedded local model conditions and constraints
	AgentPopulations //	embedded slices of each agent type

	Om     chan gobr.OutMsg        // Outgoing comm channel – dispatches batch render instructions
	Im     chan gobr.InMsg         // Incoming comm channel – receives user control messages
	e      chan error              // error message channel - general
	Quit   chan struct{}           // WebSckt monitor signal - external stop signal on ch close
	halt   chan struct{}           // exec engine halt signal on ch close
	render chan render.AgentRender // VIS message channel

	turnSync *gobr.SignalHub // synchronisation

	Stats  //	embedded global agent population statistics
	DatBuf //	embedded buffer of last turn agent pop record for LOG
}

// AgentPopulations collects slices of agent types of the `abm` package active in a model instance.
type AgentPopulations struct {
	popCpPrey         []ColourPolymorphicPrey // current prey agent population
	popVisualPredator []VisualPredator        // current predator agent population
}

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

// ConditionParams groups the CONSTANT LOCAL model conditions and constraints into a single set
type ConditionParams struct {
	Environment              `json:"abm-environment"` //	embedded model environment
	CpPreyPopulationStart    int                      `json:"abm-cp-prey-pop-start"`               // starting Prey agent population size
	CpPreyPopulationCap      int                      `json:"abm-cp-prey-pop-cap"`                 //
	CpPreyAgeing             bool                     `json:"abm-cp-prey-ageing"`                  //
	CpPreyLifespan           int                      `json:"abm-cp-prey-lifespan"`                // Prey agent lifespan
	CpPreyS                  float64                  `json:"abm-cp-prey-speed"`                   // Prey agent speed
	CpPreyA                  float64                  `json:"abm-cp-prey-acceleration"`            // Prey agent acceleration
	CpPreyTurn               float64                  `json:"abm-cp-prey-turn"`                    // Prey agent turn rate / range (in radians)
	CpPreySr                 float64                  `json:"abm-cp-prey-sr"`                      // Prey agent search range for mating
	CpPreyGestation          int                      `json:"abm-cp-prey-gestation"`               // Prey gestation period
	CpPreySexualCost         int                      `json:"abm-cp-prey-sexual-cost"`             // Prey sexual rest cost
	CpPreyReproductionChance float64                  `json:"abm-cp-prey-reproduction-chance"`     // chance of CP Prey  copulation success.
	CpPreySpawnSize          int                      `json:"abm-cp-prey-spawn-size"`              // possible number of progeny = [1, max]
	CpPreyMutationFactor     float64                  `json:"abm-cp-prey-mf"`                      // mutation factor
	VpPopulationStart        int                      `json:"abm-vp-pop-start"`                    // starting Predator agent population size
	VpPopulationCap          int                      `json:"abm-vp-pop-cap"`                      //
	VpAgeing                 bool                     `json:"abm-vp-ageing"`                       //
	VpLifespan               int                      `json:"abm-vp-lifespan"`                     // Visual Predator lifespan
	VpStarvationPoint        int                      `json:"abm-vp-starvation-point"`             //
	VpPanicPoint             int                      `json:"abm-vp-panic-point"`                  //
	VpGestation              int                      `json:"abm-vp-gestation"`                    // Visual Predator gestation period
	VpSexualRequirement      int                      `json:"abm-vp-sex-req"`                      //
	VpMovS                   float64                  `json:"abm-vp-speed"`                        // Visual Predator speed
	VpMovA                   float64                  `json:"abm-vp-acceleration"`                 // Visual Predator acceleration
	VpTurn                   float64                  `json:"abm-vp-turn"`                         // Visual Predator turn rate / range (in radians)
	VpVsr                    float64                  `json:"abm-vp-vsr"`                          // Visual Predator visual search range
	VpVb𝛄                    float64                  `json:"abm-vp-visual-search-tolerance"`      //
	VpV𝛄Bump                 float64                  `json:"abm-vp-visual-search-tolerance-bump"` //
	VpVbε                    float64                  `json:"abm-vp-baseline-col-sig-strength"`    // baseline colour signal strength factor
	VpVmε                    float64                  `json:"abm-vp-max-col-sig-strength"`         // max limit colour signal strength factor
	VpReproductionChance     float64                  `json:"abm-vp-reproduction-chance"`          // chance of VP copulation success.
	VpSpawnSize              int                      `json:"abm-vp-spawn-size"`                   //
	VpSearchChance           float64                  `json:"abm-vp-vsr-chance"`                   //
	VpAttackChance           float64                  `json:"abm-vp-attack-chance"`                //
	VpBaseAttackGain         float64                  `json:"abm-vp-baseline-attack-gain"`         //
	VpCaf                    float64                  `json:"abm-vp-col-adaptation-factor"`        //
	VpStarvation             bool                     `json:"abm-vp-starvation"`                   //
	RandomAges               bool                     `json:"abm-random-ages"`                     //	flag determining if agent ages are randomised
	RNGRandomSeed            bool                     `json:"abm-rng-random-seed"`                 // flag for using server-set random seed val.
	RNGSeedVal               int64                    `json:"abm-rng-seedval"`                     // RNG seed value
	Fuzzy                    float64                  `json:"abm-rng-fuzziness"`                   //	random 'fuzziness' offset
	Logging                  bool                     `json:"abm-logging-flag"`                    // log abm on/off
	LogFreq                  int                      `json:"abm-log-frequency"`                   // # of turns between writing log files. Default = 0
	UseCustomLogPath         bool                     `json:"abm-use-custom-log-filepath"`         //
	CustomLogPath            string                   `json:"abm-custom-log-filepath"`             //
	LogPath                  string                   `json:"abm-log-filepath"`                    //	Default logging filepath unless UseCustomLogPath is ON
	Visualise                bool                     `json:"abm-visualise-flag"`                  // Visualisation on/off
	VisFreq                  int                      `json:"abm-visualise-freq"`                  //	# of turns between sending draw instructions to web client. Default = 0
	LimitDuration            bool                     `json:"abm-limit-duration"`                  //
	FixedDuration            int                      `json:"abm-fixed-duration"`                  // fixed abm running length.
	SessionIdentifier        string                   `json:"abm-session-identifier"`              // user-friendly string (from client) to identify session
}

/*
Timeframe holds the model's representation of the time metrics.
Turn – The cycle length for all agents ∈ 𝐄 to perform 1 (and only 1) Action.
Phase – Division of a Turn, between agent sets, environmental effects/factors,
        and updates to populations and model conditionss (via external).
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

// Stats holds global statistics of the model instance.
type Stats struct {
	numCpPreyCreated int
	numCpPreyEaten   int
	numCpPreyDeath   int
	numVpCreated     int
	numVpDeath       int
}

// DatBuf is a wrapper for the buffered agent data saved for logging.
type DatBuf struct {
	recordCPP map[string]ColourPolymorphicPrey
	rcpPreyRW sync.RWMutex
	recordVP  map[string]VisualPredator
	rvpRW     sync.RWMutex
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
	m := Model{}
	m.timestamp = fmt.Sprintf("%s", time.Now())
	m.running = false
	m.Timeframe = Timeframe{}
	m.Environment = DefaultEnvironment
	m.ConditionParams = PresetParams
	m.LogPath = path.Join(os.Getenv("HOME")+os.Getenv("HOMEPATH"), abmlogPath, m.SessionIdentifier, m.timestamp)
	m.recordCPP = make(map[string]ColourPolymorphicPrey)
	m.recordVP = make(map[string]VisualPredator)
	m.Om = make(chan gobr.OutMsg)
	m.Im = make(chan gobr.InMsg)
	m.e = make(chan error)
	m.Quit = make(chan struct{})
	m.halt = make(chan struct{})
	m.render = make(chan render.AgentRender)
	m.turnSync = gobr.NewSignalHub()
	return &m
}

// PopLog prints the current time and populations
func (m *Model) PopLog() {
	log.Printf("%04dT : %04dP : %04dA\n", m.Turn, m.Phase, m.Action)
	log.Printf("cpPrey population size = %v\n", len(m.popCpPrey))
	log.Printf("vp population size = %v\n", len(m.popVisualPredator))
}

func uuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
