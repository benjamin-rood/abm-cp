package abm

import (
  "encoding/json"
  "errors"
  "fmt"
  "math/rand"
  "time"

  "github.com/benjamin-rood/gobr"
  "github.com/davecgh/go-spew/spew"
)

const (
  pause       = 1 * time.Second
  blocking    = false
  nonblocking = true
)

// Controller processes instructions from client (web, command-line)
// for now, we just send errors to the general error channel for the model instance (m.e)
func (m *Model) Controller() {
  signature := "CONTROLLER_" + m.SessionIdentifier
  for {
    select {
    case msg := <-m.Im:
      switch msg.Type {
      case "conditions": //	if conditions params msg is recieved, (re)start
        if m.running {
        register:
          clash := gobr.WaitForSignalOnce(signature, m.turnSync)
          if clash {
            time.Sleep(pause)
            goto register
          } //	will block until receiving turn broadcast once.
          m.e <- m.Stop()
        }
        err := json.Unmarshal(msg.Data, &m.ConditionParams)
        if err != nil {
          errString := fmt.Sprintf("model Controller(): error: json.Unmarshal: %s", err)
          m.e <- errors.New(errString)
          break
        }
        m.Timeframe.Reset()
        spew.Dump(m.ConditionParams)
        m.e <- m.Start()
      case "pause":
        m.e <- m.Suspend()
      }
    case <-m.Quit:
      gobr.WaitForSignalOnce(signature, m.turnSync) //	will block until receiving turn broadcast once.
      return
    }
  }
}

// Start the agent-based model
func (m *Model) Start() error {
  if m.running {
    return errors.New("Model: Start() failed: model already running")
  }
  m.running = true
  if m.RNGRandomSeed {
    rand.Seed(time.Now().UnixNano())
  } else {
    rand.Seed(m.RNGSeedVal)
  }
  timestamp := fmt.Sprintf("%s", time.Now())
  m.popCpPrey = GenerateCpPreyPopulation(m.CpPreyPopulationStart, m.numCpPreyCreated, m.Turn, m.ConditionParams, timestamp)
  m.numCpPreyCreated += m.CpPreyPopulationStart
  m.popVisualPredator = GenerateVPredatorPopulation(m.VpPopulationStart, m.numVpCreated, m.Turn, m.ConditionParams, timestamp)
  m.numVpCreated += m.VpPopulationStart
  if m.Logging {
    go m.log(m.e)
  }
  if m.Visualise {
    go m.vis(m.e)
  }
  time.Sleep(pause)
  go m.run(m.e)
  return nil
}

// Stop the agent-based model
func (m *Model) Stop() error {
  if !m.running {
    return errors.New("Model: Stop() failed: model not currently running")
  }
  close(m.halt)
  m.running = false
  m.popCpPrey = nil
  m.popVisualPredator = nil
  time.Sleep(pause)
  m.halt = make(chan struct{})
  return nil
}

// Suspend = pause a running agent-based model to be resumed later.
func (m *Model) Suspend() error {
  if !m.running {
    return errors.New("Model: Suspend() failed: model not currently running")
  }
  close(m.halt)
  m.running = false
  time.Sleep(pause)
  return nil
}

// Resume from a suspended agent-based model
func (m *Model) Resume() error {
  if m.running {
    return errors.New("Model: Resume() failed: model already running")
  }
  m.halt = make(chan struct{})
  go m.run(m.e)
  m.running = true
  if m.Visualise {
    go m.vis(m.e)
  }
  if m.Logging {
    go m.log(m.e)
  }
  return nil
}
