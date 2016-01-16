package abm

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/render"
	"github.com/benjamin-rood/goio"
)

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

// Log prints the current state of time
func (m *Model) Log() {
	log.Printf("%04dT : %04dP : %04dA\n", m.Turn, m.Phase, m.Action)
	log.Printf("cpp population size = %d\n", len(m.PopCPP))
	log.Printf("vp population size = %d\n", len(m.PopVP))
}

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
	RandomSeed            bool      //	flag for using server-set random seed val.
	Seed                  int       //	RNG seed value
	Bounds                []float64 // d value for each axis
	MaxCppPopSize         int
	MaxVpPopSize          int
	StartCppPopSize       int // starting CPP agent population size
	StartVpPopSize        int //	starting VP agent population size
	VpAgeing              bool
	VpLifespan            int     //	Visual Predator lifespan
	VS                    float64 // Visual Predator speed
	VA                    float64 // Visual Predator acceleration
	VpTurn                float64 //	Visual Predator turn rate / range (in radians)
	Vsr                   float64 //	VP agent visual search range
	Vγ                    float64 //	visual acuity in environments
	VpReproductiveChance  float64 //	chance of VP copulation success.
	VsrSearchChance       float64
	VpAttackChance        float64
	VpColImprintFactor    float64
	CppAgeing             bool
	CppLifespan           int     //	CPP agent lifespan
	CppS                  float64 // CPP agent speed
	CppA                  float64 // CPP agent acceleration
	CppTurn               float64 //	CPP agent turn rate / range (in radians)
	CppSr                 float64 // CPP agent search range for mating
	MutationFactor        float64 //	mutation factor
	CppGestation          int     //	CPP gestation period
	CppSexualCost         int     //	CPP sexual rest cost
	CppReproductiveChance float64 //	chance of CPP copulation success.
	CppSpawnSize          int     // 	CPP max spawn size s.t. possible number of progeny = [1, max]
	RandomAges            bool
	Fuzzy                 float64
}

func runningModel(m Model, viz chan<- render.AgentRender, quit <-chan struct{}, phase chan<- struct{}) {
	var am sync.Mutex
	var cppAgentWg sync.WaitGroup
	// var vpAgentWg sync.WaitGroup
	for {
		select {
		case <-quit:
			// clean up, then...
			return
		default: //	PROCEED WITH TURN
			var cppAgents []ColourPolymorphicPrey
			for i := range m.PopCPP {
				cppAgentWg.Add(1)
				go func(i int) {
					defer cppAgentWg.Done()
					result := m.PopCPP[i].RBB(m.Context, len(m.PopCPP))
					viz <- m.PopCPP[i].GetDrawInfo()
					time.Sleep(time.Millisecond * 25)
					am.Lock()
					cppAgents = append(cppAgents, result...)
					am.Unlock()
					m.Action++
				}(i)
			}
			cppAgentWg.Wait()
			m.PopCPP = cppAgents //	update the population based on the results from each agent's rule-based behaviour of the turn.
			m.Phase++
			m.Action = 0 // reset at phase end
			m.Log()
			time.Sleep(time.Millisecond * 50)

			var vpAgents []VisualPredator
			for i := range m.PopVP {
				// vpAgentWg.Add(1)
				func(i int) {
					// defer vpAgentWg.Done()
					result := m.PopVP[i].RBB(m.Context, m.PopCPP)
					viz <- m.PopVP[i].GetDrawInfo()
					am.Lock()
					vpAgents = append(vpAgents, result...)
					am.Unlock()
					m.Action++
				}(i)
			}
			// vpAgentWg.Wait()
			m.PopVP = vpAgents //	update the population based on the results from each agent's rule-based behaviour of the turn.

			m.Phase++
			m.Action = 0 // reset at phase end
			m.Log()
			time.Sleep(time.Millisecond * 50)
			phase <- struct{}{}

			m.Phase = 0 //	reset at Turn end
			m.Turn++
			m.Log()
		}
	}
}

// insufficient hack
func InitModel(ctxt Context, e Environment, om chan goio.OutMsg, phase chan struct{}) {
	simple := setModel(ctxt, e)
	quit := make(chan struct{})
	ar := make(chan render.AgentRender, 2000)
	go runningModel(simple, ar, quit, phase)
	go visualiseModel(ctxt, ar, om, phase)
}

func setModel(ctxt Context, e Environment) (m Model) {
	m.PopCPP = GeneratePopulationCPP(ctxt.StartCppPopSize, ctxt)
	//m.PopVP = GeneratePopulationVP(ctxt.StartVpPopSize, ctxt)
	m.DefinitionCPP = []string{"mover", "breeder", "mortal"}
	m.DefinitionVP = []string{"mover", "hunter", "mortal"}
	m.Environment = e
	m.Context = ctxt
	return
}

func visualiseModel(ctxt Context, ar <-chan render.AgentRender, out chan<- goio.OutMsg, turn <-chan struct{}) {
	// v := DemoViewport
	rand.Seed(time.Now().UnixNano())
	bg := colour.RGB256{Red: 30, Green: 30, Blue: 30}
	msg := goio.OutMsg{Type: "render", Data: nil}
	dl := render.DrawList{
		CPP: nil,
		VP:  nil,
		BG:  bg,
	}
	for {
		select {
		case job := <-ar:
			// offload the viewport translation to client (browser) side.
			//job.TranslateToViewport(v, ctxt.Bounds[0], ctxt.Bounds[1])
			switch job.Type {
			case "cpp":
				dl.CPP = append(dl.CPP, job)
			case "vp":
				dl.VP = append(dl.VP, job)
			default:
				log.Fatalf("viz: failed to determine agent-render job type!")
			}
		case <-turn:
			msg.Data = dl
			out <- msg
			// reset msg contents
			msg = goio.OutMsg{Type: "render", Data: nil}
			//	reset draw instructions
			dl = render.DrawList{
				CPP: nil,
				VP:  nil,
				BG:  bg,
			}
			// removed viewport update case – switching this responsibility completely to client side.
		}
	}
}
