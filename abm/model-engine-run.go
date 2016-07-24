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
    case <-m.halt:
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

func (m *Model) cpPreyPhase(errCh chan<- error) []ColourPolymorphicPrey {
  var mutex sync.Mutex
  var Wg sync.WaitGroup
  var agentsUpdate []ColourPolymorphicPrey

  for i := range m.popCpPrey {
    Wg.Add(1)
    go func(agent ColourPolymorphicPrey) {
      defer func() {
        Wg.Done()
        if m.Logging {
          // do this copying to the record in a goroutine once proven stable and safe!
          errCh <- m.cpPreyRecordAssignValue(agent.UUID(), agent)
        }
      }()
      result := agent.Action(m.ConditionParams, len(m.popCpPrey))
      if m.Visualise {
        m.render <- agent.GetDrawInfo()
      }
      mutex.Lock()
      m.numCpPreyCreated += len(result) - 1
      agentsUpdate = append(agentsUpdate, result...)
      mutex.Unlock()
      m.Action++
    }(m.popCpPrey[i])
  }
  Wg.Wait()
  return agentsUpdate
}

func (m *Model) visualPredatorPhase(errCh chan<- error) []VisualPredator {
  var mutex sync.Mutex
  var agentsUpdate []VisualPredator
  for i := range m.popVisualPredator {
    func(agent VisualPredator) {
      defer func() {
        if m.Logging {
          // do this copying to the record in a seperate goroutine once proven stable and safe!
          errCh <- m.vpRecordAssignValue(agent.UUID(), agent)
        }
      }()
      result := agent.Action(errCh, m.ConditionParams, m.numVpCreated, m.Turn, m.popCpPrey, m.popVisualPredator, i)
      if m.Visualise {
        m.render <- agent.GetDrawInfo()
      }
      mutex.Lock()
      m.numVpCreated += len(result) - 1
      agentsUpdate = append(agentsUpdate, result...)
      mutex.Unlock()
      m.Action++
    }(m.popVisualPredator[i])
  }
  return agentsUpdate
}

func (m *Model) turn(errCh chan<- error) {
  m.popCpPrey = m.cpPreyPhase(errCh) // update the population based on the results from all Prey agents rule-based behaviour in the phase.
  m.Phase++
  m.Action = 0                                       // reset at phase end
  m.popVisualPredator = m.visualPredatorPhase(errCh) // update the population based on the results from all Predators rule-based behaviour in the phase.
  m.Phase++
  m.Action = 0                   // reset at phase end
  m.Phase = 0                    // reset at Turn end
  m.turnSync.Broadcast(blocking) // using blocking version to ensure synchronisation with the other processes in the active Engine Set.
  m.Turn++
}
