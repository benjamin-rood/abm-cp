package abm

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/render"
	"github.com/benjamin-rood/goio"
	"github.com/davecgh/go-spew/spew"
)

const (
	pause = 1 * time.Second
)

// Controller processes instructions from web client
func (m *Model) Controller() {
	for {
		select {
		case msg := <-m.Im:
			switch msg.Type {
			case "context": //	if context params msg is recieved, (re)start
				err := json.Unmarshal(msg.Data, &m.Context)
				if err != nil {
					log.Println("model Controller(): error: json.Unmarshal:", err)
					break
				}
				m.Timeframe.Reset()
				spew.Dump(m.Context)
				if m.running {
					m.Stop()
				}
				m.Start()
			}
		case <-m.Quit:
			return
		}
	}
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

// Stop the agent-based model
func (m *Model) Stop() {
	close(m.r)
	m.running = false
	m.PopCPP = nil
	m.PopVP = nil
	time.Sleep(pause)
	m.r = make(chan struct{})
}

// Suspend = pause an agent-based model to be resumed later.
func (m *Model) Suspend() {
	close(m.r)
	time.Sleep(pause)
}

// Resume from a suspended agent-based model
func (m *Model) Resume() {
	m.r = make(chan struct{})
	ar := make(chan render.AgentRender)
	turn := make(chan struct{})
	go m.run(ar, turn)
	go m.vis(ar, turn)
}

func (m *Model) run(ar chan<- render.AgentRender, turn chan<- struct{}) {
	time.Sleep(time.Second)
	for {
		select {
		case <-m.r:
			close(ar)
			close(turn)
			time.Sleep(pause)
			return
		case <-m.Quit:
			// clean up?
			time.Sleep(time.Millisecond * 250)
			return
		default: //	PROCEED WITH TURN
			if (len(m.PopCPP) == 0) || (len(m.PopVP) == 0) {
				m.Stop()
			}
			m.turn(ar, turn)
		}
	}
}

func (m *Model) turn(ar chan<- render.AgentRender, turn chan<- struct{}) {
	var am sync.Mutex
	var cppAgentWg sync.WaitGroup
	var vpAgentWg sync.WaitGroup
	var cppAgents []ColourPolymorphicPrey
	// cInterval := time.Now()
	for i := range m.PopCPP {
		cppAgentWg.Add(1)
		// timeMark := time.Now()
		go func(i int) {
			defer cppAgentWg.Done()
			result := m.PopCPP[i].RBB(m.Context, len(m.PopCPP))
			ar <- m.PopCPP[i].GetDrawInfo()
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
	time.Sleep(time.Millisecond * 20)

	var vpAgents []VisualPredator
	for i := range m.PopVP {
		vpAgentWg.Add(1)
		go func(i int) {
			defer vpAgentWg.Done()
			result := m.PopVP[i].RBB(m.Context, len(m.PopVP), m.PopCPP, m.PopVP, i)
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
	time.Sleep(time.Millisecond * 20)
	turn <- struct{}{}

	m.Phase = 0 //	reset at Turn end
	m.Turn++
}

func (m *Model) vis(ar <-chan render.AgentRender, turn <-chan struct{}) {
	msg := goio.OutMsg{Type: "render", Data: nil}
	bg := m.BG.To256()
	dl := render.DrawList{
		CPP:       nil,
		VP:        nil,
		BG:        bg,
		CppPop:    "0",
		VpPop:     "0",
		TurnCount: "0",
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
			dl.CppPop = fmt.Sprintf("cpp %d", len(m.PopCPP))
			dl.VpPop = fmt.Sprintf("vp  %d", len(m.PopVP))
			dl.TurnCount = fmt.Sprintf("%08d", m.Turn)
			msg.Data = dl
			m.Om <- msg
			// reset msg contents
			msg = goio.OutMsg{Type: "render", Data: nil}
			//	reset draw instructions
			dl = render.DrawList{
				CPP:       nil,
				VP:        nil,
				BG:        bg,
				CppPop:    "0",
				VpPop:     "0",
				TurnCount: "0",
			}
		case <-m.r:
			return
		}
	}
}

// Kill off the model and any client bound to it.
func (m *Model) Kill() {
	m.Stop()
	close(m.Quit)
	m.Dead = true
}
