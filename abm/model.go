package abm

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/render"
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
	CppSpawnSize          int       `json:"abm-cpp-spawn-size"`          // 	CPP max spawn size s.t. possible number of progeny = [1, max]
	VpPopulationStart     int       `json:"abm-vp-pop-start"`            //	starting VP agent population size
	VpPopulationCap       int       `json:"abm-vp-pop-cap"`
	VpAgeing              bool      `json:"abm-vp-ageing"`
	VpLifespan            int       `json:"abm-vp-lifespan"` //	Visual Predator lifespan
	VS                    float64   `json:"abm-vp-speed"`    // Visual Predator speed
	VA                    float64   // Visual Predator acceleration
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

// Model acts as the working instance of the 'game'
type Model struct {
	running bool
	Dead    bool
	Timeframe
	Environment
	Context
	PopulationCPP
	PopulationVP
	Om chan goio.OutMsg
	Im chan goio.InMsg
	e  chan error
	Q  chan struct{}
	r  chan struct{}
}

// NewModel is a constructor for initialising a Model instance
func NewModel() (m Model) {
	m.running = false
	m.Timeframe = Timeframe{}
	m.Environment = Environment{
		Bounds:         []float64{1.0, 1.0},
		Dimensionality: 2,
		BG:             colour.RandRGB(),
	}
	m.Context = DemoContext
	m.PopulationCPP = PopulationCPP{}
	m.PopulationVP = PopulationVP{}
	m.Om = make(chan goio.OutMsg)
	m.Im = make(chan goio.InMsg)
	m.e = make(chan error)
	m.Q = make(chan struct{})
	m.r = make(chan struct{})
	return
}

// Controller processes instructions from web client
func (m *Model) Controller() {
	for {
		select {
		case msg := <-m.Im:
			switch msg.Type {
			case "start":
				json.Unmarshal(msg.Data, &m.Context)
				if m.running {
					m.Stop()
					m.Reset()
				}
				m.Start()
			}
		}
	}
}

// Reset resets all non-channel fields
func (m *Model) Reset() {
	m.running = false
	m.Timeframe = Timeframe{}
	m.Environment = Environment{
		Bounds:         []float64{1.0, 1.0},
		Dimensionality: 2,
		BG:             colour.RandRGB(),
	}
	m.Context = DemoContext
	m.PopulationCPP = PopulationCPP{}
	m.PopulationVP = PopulationVP{}
}

// Log prints the current state of time
func (m *Model) Log() {
	log.Printf("%04dT : %04dP : %04dA\n", m.Turn, m.Phase, m.Action)
	log.Printf("cpp population size = %d\n", len(m.PopCPP))
	log.Printf("vp population size = %d\n", len(m.PopVP))
}

func (m *Model) Stop() {
	close(m.r)
	m.running = false
	time.Sleep(1 * time.Second)
	m.r = make(chan struct{})
}

// Start the agent-based model
func (m *Model) Start() {
	m.running = true
	ar := make(chan render.AgentRender)
	turn := make(chan struct{})
	if m.RNGRandomSeed {
		rand.Seed(time.Now().UnixNano())
	} else {
		rand.Seed(m.RNGSeedVal)
	}
	m.PopCPP = GeneratePopulationCPP(m.CppPopulationStart, m.Context)
	m.PopVP = GeneratePopulationVP(m.VpPopulationStart, m.Context)
	go m.run(ar, turn)
	go m.vis(ar, turn)
}

func (m *Model) run(ar chan<- render.AgentRender, turn chan<- struct{}) {

	var am sync.Mutex
	var cppAgentWg sync.WaitGroup
	var vpAgentWg sync.WaitGroup

	for {
		select {
		case <-m.r:
			time.Sleep(time.Second)
			close(ar)
			close(turn)
			return
		default: //	PROCEED WITH TURN
			var cppAgents []ColourPolymorphicPrey
			cInterval := time.Now()
			for i := range m.PopCPP {
				cppAgentWg.Add(1)
				timeMark := time.Now()
				go func(i int) {
					defer cppAgentWg.Done()
					result := m.PopCPP[i].RBB(m.Context, len(m.PopCPP))
					ar <- m.PopCPP[i].GetDrawInfo()
					am.Lock()
					cppAgents = append(cppAgents, result...)
					am.Unlock()
					time.Sleep(time.Millisecond * 10)
					m.Action++
				}(i)
				fmt.Printf("cp-rbb: %04d elapsed: %v\n", i, time.Since(timeMark))
			}
			fmt.Printf("m.PopCPP: %04d total cp-rbb elapsed: %v\n", len(m.PopCPP), time.Since(cInterval))
			cppAgentWg.Wait()
			m.PopCPP = cppAgents //	update the population based on the results from each agent's rule-based behaviour of the turn.
			m.Phase++
			m.Action = 0 // reset at phase end
			// m.Log()
			time.Sleep(time.Millisecond * 50)

			var vpAgents []VisualPredator
			for i := range m.PopVP {
				vpAgentWg.Add(1)
				go func(i int) {
					defer vpAgentWg.Done()
					result := m.PopVP[i].RBB(m.Context, m.PopCPP)
					ar <- m.PopVP[i].GetDrawInfo()
					am.Lock()
					vpAgents = append(vpAgents, result...)
					am.Unlock()
					m.Action++
				}(i)
			}
			vpAgentWg.Wait()
			m.PopVP = vpAgents //	update the population based on the results from each agent's rule-based behaviour of the turn.

			m.Phase++
			m.Action = 0 // reset at phase end
			// m.Log()
			time.Sleep(time.Millisecond * 50)
			turn <- struct{}{}

			m.Phase = 0 //	reset at Turn end
			m.Turn++
			// m.Log()
		}
	}
}

func (m *Model) vis(ar <-chan render.AgentRender, turn <-chan struct{}) {
	msg := goio.OutMsg{Type: "render", Data: nil}
	bg := m.BG.To256()
	dl := render.DrawList{
		CPP: nil,
		VP:  nil,
		BG:  bg,
	}
	for {
		select {
		case job := <-ar:
			switch job.Type {
			case "cpp":
				dl.CPP = append(dl.CPP, job)
			case "vp":
				dl.VP = append(dl.VP, job)
			}
		case <-turn:
			msg.Data = dl
			m.Om <- msg
			// reset msg contents
			msg = goio.OutMsg{Type: "render", Data: nil}
			//	reset draw instructions
			dl = render.DrawList{
				CPP: nil,
				VP:  nil,
				BG:  bg,
			}
		case <-m.r:
			return
		}
	}
}
