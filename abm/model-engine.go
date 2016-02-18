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
	pause = 1 * time.Second
)

// Controller processes instructions from client (web, command-line)
// for now, we just send errors to the general error channel for the model instance (m.e)
func (m *Model) Controller() {
	signature := "CONTROLLER_" + m.SessionIdentifier
	for {
		select {
		case msg := <-m.Im:
			switch msg.Type {
			case "context": //	if context params msg is recieved, (re)start
				if m.running {
				register:
					clash := gobr.WaitForSignalOnce(signature, m.turnSignal)
					if clash {
						time.Sleep(pause)
						goto register
					} //	will block until receiving turn broadcast once.
					m.e <- m.Stop()
				}
				err := json.Unmarshal(msg.Data, &m.Context)
				if err != nil {
					errString := fmt.Sprintf("model Controller(): error: json.Unmarshal: %s", err)
					m.e <- errors.New(errString)
					break
				}
				m.Timeframe.Reset()
				spew.Dump(m.Context)
				m.e <- m.Start()
			case "pause":
				m.e <- m.Suspend()
			}
		case <-m.Quit:
			gobr.WaitForSignalOnce(signature, m.turnSignal) //	will block until receiving turn broadcast once.
			return
		}
	}
}

// Start the agent-based model
func (m *Model) Start() error {
	// _ = "breakpoint" // godebug
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
	m.PopCPP = GeneratePopulationCPP(m.CppPopulationStart, m.numCppCreated, m.Turn, m.Context, timestamp)
	m.numCppCreated += m.CppPopulationStart
	m.PopVP = GeneratePopulationVP(m.VpPopulationStart, m.numVpCreated, m.Turn, m.Context, timestamp)
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
		return errors.New("Model: Stop() failed: model not currently running!")
	}
	close(m.rc)
	m.running = false
	m.PopCPP = nil
	m.PopVP = nil
	time.Sleep(pause)
	m.rc = make(chan struct{})
	return nil
}

// Suspend = pause an agent-based model to be resumed later.
func (m *Model) Suspend() error {
	if !m.running {
		return errors.New("Model: Suspend() failed: model not currently running!")
	}
	close(m.rc)
	m.running = false
	time.Sleep(pause)
	return nil
}

// Resume from a suspended agent-based model
func (m *Model) Resume() error {
	if m.running {
		return errors.New("Model: Resume() failed: model already running")
	}
	m.rc = make(chan struct{})
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

// kill off the model and any client bound to it - internal killoff switch
func (m *Model) kill() {
	m.Stop()
	close(m.Quit)
	m.Dead = true
}
