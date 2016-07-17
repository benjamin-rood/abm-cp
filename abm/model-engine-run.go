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
		case <-m.rq:
			gobr.WaitForSignalOnce(signature, m.turnSync) // block while waiting for turn to end.
			time.Sleep(pause)
			return
		case <-m.Quit:
			gobr.WaitForSignalOnce(signature, m.turnSync) // block while waiting for turn to end.
			ec <- m.Stop()
			time.Sleep(pause)
			return
		default:
			if m.LimitDuration && m.Turn >= m.FixedDuration {
				ec <- m.Stop()
				return
			}
			if (len(m.popCpPrey) == 0) || (len(m.popVisualPredator) == 0) {
				ec <- m.Stop()
				return
			}
			//	PROCEED WITH TURN
			m.turn(ec)
		}
	}
}

func (m *Model) turn(errCh chan<- error) {
	var am sync.Mutex
	var cpPreyAgentWg sync.WaitGroup
	var cpPreyAgents []ColourPolymorphicPrey

	for i := range m.popCpPrey {
		cpPreyAgentWg.Add(1)
		go func(agent ColourPolymorphicPrey) {
			defer func() {
				cpPreyAgentWg.Done()
				if m.Logging {
					// do this copying to the record in a goroutine once proven stable and safe!
					errCh <- m.cpPreyRecordAssignValue(agent.UUID(), agent)
				}
			}()
			result := agent.RBB(m.ConditionParams, len(m.popCpPrey))
			if m.Visualise {
				m.render <- agent.GetDrawInfo()
			}
			am.Lock()
			cpPreyAgents = append(cpPreyAgents, result...)
			am.Unlock()
			m.Action++
		}(m.popCpPrey[i])
	}

	cpPreyAgentWg.Wait()
	m.popCpPrey = cpPreyAgents //	update the population based on the results from each agent's rule-based behaviour of the turn.
	m.Phase++
	m.Action = 0 // reset at phase end
	time.Sleep(time.Millisecond * 20)
	_ = "breakpoint" // godebug

	var vpAgents []VisualPredator

	for i := range m.popVisualPredator {
		func(agent VisualPredator) {
			defer func() {
				if m.Logging {
					// do this copying to the record in a goroutine once proven stable and safe!
					errCh <- m.vpRecordAssignValue(agent.UUID(), agent)
				}
			}()
			result := agent.RBB(errCh, m.ConditionParams, m.numVpCreated, m.Turn, m.popCpPrey, m.popVisualPredator, i)
			if m.Visualise {
				m.render <- agent.GetDrawInfo()
			}
			am.Lock()
			m.numVpCreated += len(result) - 1
			vpAgents = append(vpAgents, result...)
			am.Unlock()
			m.Action++
		}(m.popVisualPredator[i])
	}

	m.popVisualPredator = vpAgents //	update the population based on the results from each agent's rule-based behaviour of the turn.
	m.Phase++
	m.Action = 0                // reset at phase end
	m.Phase = 0                 // reset at Turn end
	_ = "breakpoint"            // godebug
	m.turnSync.Broadcast(false) // using blocking version to ensure turn synchronisation
	m.Turn++
	time.Sleep(time.Millisecond * 50) // sync wiggle room
}
