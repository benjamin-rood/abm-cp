package abm

import (
	"sync"
	"time"
)

func (m *Model) run(ec chan<- error) {
	_ = "breakpoint" // godebug
	var turn chan struct{}
	var clash bool
	var signature string

	for {
		signature = uuid()
		turn, clash = m.turnSignal.Register(signature)
		if !clash {
			break
		}
	}

	defer func() {
		m.turnSignal.Deregister(signature)
		// Need to wipe the agent records too? -yes, probably.
	}()

	for {
		select {
		case <-m.rc:
			<-turn // block while waiting for turn to end.
			time.Sleep(pause)
			return
		case <-m.Quit:
			<-turn // block while waiting for turn to end.
			ec <- m.Stop()
			time.Sleep(pause)
			return
		default:
			if (len(m.PopCPP) == 0) || (len(m.PopVP) == 0) {
				ec <- m.Stop()
				return
			}
			m.turn(ec) //	PROCEED WITH TURN
		}
		if m.LimitDuration && m.Turn >= m.FixedDuration {
			ec <- m.Stop()
			return
		}
	}
}

func (m *Model) turn(errCh chan<- error) {
	var am sync.Mutex
	var cppAgentWg sync.WaitGroup
	var vpAgentWg sync.WaitGroup
	var cppAgents []ColourPolymorphicPrey

	for i := range m.PopCPP {
		cppAgentWg.Add(1)
		go func(agent ColourPolymorphicPrey) {
			defer func() {
				cppAgentWg.Done()
				if m.Logging {
					// do this copying to the record in a goroutine once proven stable and safe!
					m.recordCPP[agent.UUID()] = agent
				}
			}()
			result := agent.RBB(m.Context, len(m.PopCPP))
			if m.Visualise {
				m.render <- agent.GetDrawInfo()
			}
			am.Lock()
			cppAgents = append(cppAgents, result...)
			am.Unlock()
			m.Action++
		}(m.PopCPP[i])
	}

	cppAgentWg.Wait()
	m.PopCPP = cppAgents //	update the population based on the results from each agent's rule-based behaviour of the turn.
	m.Phase++
	m.Action = 0 // reset at phase end
	time.Sleep(time.Millisecond * 20)

	var vpAgents []VisualPredator
	for i := range m.PopVP {
		vpAgentWg.Add(1)
		go func(agent VisualPredator, selfIndex int) {
			defer func() {
				vpAgentWg.Done()
				if m.Logging {
					// do this copying to the record in a goroutine once proven stable and safe!
					m.recordVP[agent.UUID()] = agent
				}
			}()
			result := agent.RBB(m.Context, m.numVpCreated, m.Turn, m.PopCPP, m.PopVP, selfIndex)
			if m.Visualise {
				m.render <- agent.GetDrawInfo()
			}
			am.Lock()
			m.numVpCreated += len(result) - 1
			vpAgents = append(vpAgents, result...)
			am.Unlock()
			m.Action++
		}(m.PopVP[i], i)
	}

	vpAgentWg.Wait()
	m.PopVP = vpAgents //	update the population based on the results from each agent's rule-based behaviour of the turn.
	m.Phase++
	m.Action = 0             // reset at phase end
	m.Phase = 0              // reset at Turn end
	_ = "breakpoint"         // godebug
	m.turnSignal.Broadcast() // use blocking version turn ensure synchronisation?
	m.Turn++
	time.Sleep(time.Millisecond * 50) // sync wiggle room
}
