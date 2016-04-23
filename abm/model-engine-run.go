package abm

import (
	"sync"
	"time"

	"github.com/benjamin-rood/gobr"
)

func (m *Model) run(ec chan<- error) {
	signature := "RUN_" + m.SessionIdentifier
	for {
		select {
		case <-m.rc:
			_ = "breakpoint"                                // godebug
			gobr.WaitForSignalOnce(signature, m.turnSignal) // block while waiting for turn to end.
			time.Sleep(pause)
			return
		case <-m.Quit:
			_ = "breakpoint"                                // godebug
			gobr.WaitForSignalOnce(signature, m.turnSignal) // block while waiting for turn to end.
			ec <- m.Stop()
			time.Sleep(pause)
			return
		default:
			_ = "breakpoint" // godebug
			if m.LimitDuration && m.Turn >= m.FixedDuration {
				ec <- m.Stop()
				return
			}
			if (len(m.PopCPP) == 0) || (len(m.PopVP) == 0) {
				ec <- m.Stop()
				return
			}
			m.turn(ec) //	PROCEED WITH TURN
			// m.PopLog()
			// time.Sleep(100 * time.Millisecond)
		}
	}
}

func (m *Model) turn(errCh chan<- error) {
	var am sync.Mutex
	var cpPreyAgentWg sync.WaitGroup
	var cpPreyAgents []ColourPolymorphicPrey

	for i := range m.PopCPP {
		cpPreyAgentWg.Add(1)
		go func(agent ColourPolymorphicPrey) {
			defer func() {
				cpPreyAgentWg.Done()
				if m.Logging {
					// do this copying to the record in a goroutine once proven stable and safe!
					errCh <- m.cpPreyRecordAssignValue(agent.UUID(), agent)
				}
			}()
			result := agent.RBB(m.Condition, len(m.PopCPP))
			if m.Visualise {
				m.render <- agent.GetDrawInfo()
			}
			am.Lock()
			cpPreyAgents = append(cpPreyAgents, result...)
			am.Unlock()
			m.Action++
		}(m.PopCPP[i])
	}

	cpPreyAgentWg.Wait()
	m.PopCPP = cpPreyAgents //	update the population based on the results from each agent's rule-based behaviour of the turn.
	m.Phase++
	m.Action = 0 // reset at phase end
	time.Sleep(time.Millisecond * 20)
	_ = "breakpoint" // godebug
	// var vpAgentWg sync.WaitGroup
	var vpAgents []VisualPredator

	for i := range m.PopVP {
		// vpAgentWg.Add(1)
		func(agent VisualPredator) {
			_ = "breakpoint" // godebug
			defer func() {
				// vpAgentWg.Done()
				if m.Logging {
					// do this copying to the record in a goroutine once proven stable and safe!
					errCh <- m.vpRecordAssignValue(agent.UUID(), agent)
				}
			}()
			result := agent.RBB(errCh, m.Condition, m.numVpCreated, m.Turn, m.PopCPP, m.PopVP, i)
			if m.Visualise {
				m.render <- agent.GetDrawInfo()
			}
			am.Lock()
			m.numVpCreated += len(result) - 1
			vpAgents = append(vpAgents, result...)
			am.Unlock()
			m.Action++
		}(m.PopVP[i])
	}

	// vpAgentWg.Wait()
	m.PopVP = vpAgents //	update the population based on the results from each agent's rule-based behaviour of the turn.
	m.Phase++
	m.Action = 0                  // reset at phase end
	m.Phase = 0                   // reset at Turn end
	_ = "breakpoint"              // godebug
	m.turnSignal.Broadcast(false) // using blocking version to ensure turn synchronisation
	m.Turn++
	time.Sleep(time.Millisecond * 50) // sync wiggle room
}
