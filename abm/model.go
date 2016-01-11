package abm

import (
	"log"
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
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
	Cφ            int     //	CPP gestation period
	Cȣ            int     //	CPP sexual rest cost
	Cκ            float64 //	chance of CPP copulation success.
	Cβ            int     // 	CPP max spawn size (birth range)
}

func cppRBB(ctxt Context, time Timeframe, pop []ColourPolymorphicPrey, queue chan<- render.AgentRender) (newpop []ColourPolymorphicPrey, newtime Timeframe) {
	newkids := []ColourPolymorphicPrey{}
	newtime = time
	for i := range pop {
		jump := ""
		// BEGIN
		if ctxt.CppAgeing {
			jump = pop[i].Age()
			switch jump {
			case "DEATH":
				goto End
			}
		}

		jump = pop[i].Fertility(ctxt.Cȣ)
		// fmt.Println("JUMP to", jump)
		switch jump {
		case "SPAWN":
			progeny := pop[i].Birth(ctxt) //	max spawn size, mutation factor
			newkids = append(newkids, progeny...)
		case "MATE SEARCH":
			if len(pop) < maxPopSize {
				mate, err := pop[i].MateSearch(pop, i) // need to exclude self from search :-)
				if err != nil {
					log.Fatalln("cppRBB: MateSearch: Error:", err)
				}
				// fmt.Println("mate = ", mate)
				// ATTEMPT REPRODUCE
				success := pop[i].Copulation(mate, ctxt.Cκ, ctxt.Cφ, ctxt.Cȣ)
				if success {
					// fmt.Println("successful sex act!")
					goto Add
				}
			}
			fallthrough // no nearby sexual partners, EXPLORE
		case "EXPLORE":
			𝚯 := calc.RandFloatIn(-ctxt.Cτ, ctxt.Cτ)
			pop[i].Turn(𝚯)
			pop[i].Move()
		}

	Add:
		newpop = append(newpop, pop[i])
		queue <- pop[i].GetDrawInfo()

	End:
		newtime.Action++
	}
	newpop = append(newpop, newkids...) // add the newly created children to the returning population
	return
}

func runningModel(m Model, rc chan<- render.AgentRender, quit <-chan struct{}, phase chan<- struct{}) {
	for {
		m.PopCPP, m.Timeframe = cppRBB(m.Context, m.Timeframe, m.PopCPP, rc) //	returns a replacement
		m.Action = 0                                                         // reset at phase end.
		m.Phase++
		m.Log()
		phase <- struct{}{}
		time.Sleep(20)
	}
}

// insufficient hack
func InitModel(ctxt Context, e Environment, om chan goio.OutMsg, view chan render.Viewport, phase chan struct{}) {
	simple := setModel(ctxt, e)
	quit := make(chan struct{})
	rc := make(chan render.AgentRender)
	go runningModel(simple, rc, quit, phase)
	go visualiseModel(ctxt, view, rc, om, phase)
}

func setModel(ctxt Context, e Environment) (m Model) {
	m.PopCPP = GeneratePopulation(cppPopSize, ctxt)
	m.DefinitionCPP = []string{"mover", "breeder", "mortal"}
	m.Environment = e
	m.Context = ctxt
	return
}

func visualiseModel(ctxt Context, view <-chan render.Viewport, queue <-chan render.AgentRender, out chan<- goio.OutMsg, phase <-chan struct{}) {
	v := DemoViewport
	msg := goio.OutMsg{Type: "render", Data: nil}
	dl := render.DrawList{
		CPP: nil,
		VP:  nil,
		BG:  colour.RGB256{Red: 0, Green: 0, Blue: 0},
	}
	for {
		select {
		case job := <-queue:
			job.TranslateToViewport(v, ctxt.Bounds[0], ctxt.Bounds[1])
			switch job.Type {
			case "cpp":
				dl.CPP = append(dl.CPP, job)
			case "vp":
				dl.VP = append(dl.VP, job)
			default:
				log.Fatalf("viz: failed to determine agent-render job type!")
			}
		case <-phase:
			msg.Data = dl
			out <- msg
			// reset msg contents
			msg = goio.OutMsg{Type: "render", Data: nil}
			//	reset draw instructions
			dl = render.DrawList{
				CPP: nil,
				VP:  nil,
				BG:  colour.RGB256{Red: 0, Green: 0, Blue: 0},
			}
		case v = <-view:
		}
	}
}
