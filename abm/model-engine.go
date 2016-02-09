package abm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/benjamin-rood/abm-cp/render"
	"github.com/benjamin-rood/gobr"
	"github.com/davecgh/go-spew/spew"
)

const (
	pause = 1 * time.Second
)

// Controller processes instructions from client (web, command-line)
// for now, we just send errors to the general error channel for the model instance (m.e)
func (m *Model) Controller() {
	var turn chan struct{}
	var clash bool
	var signature string
	defer func() {
		m.turnSignal.Deregister(signature)
		// Need to wipe the agent records too? -yes, probably.
	}()

Registration: // register to receive from m.turnSignal - loop until no clash with existing entry in map.
	signature = uuid()
	turn, clash = m.turnSignal.Register(signature)
	if clash {
		goto Registration
	}

	for {
		select {
		case msg := <-m.Im:
			switch msg.Type {
			case "context": //	if context params msg is recieved, (re)start
				time.Sleep(pause)
				err := json.Unmarshal(msg.Data, &m.Context)
				if err != nil {
					errString := fmt.Sprintf("model Controller(): error: json.Unmarshal: %s", err)
					m.e <- errors.New(errString)
					break
				}
				<-turn //	block while waiting for turn to end.
				m.Timeframe.Reset()
				spew.Dump(m.Context)
				m.e <- m.Stop()
				m.e <- m.Start()
			case "pause":
				m.e <- m.Suspend() // if Suspend returns nil error, then the handler will just discard! Simple!
			}
		case <-m.Quit:
			<-turn //	block while waiting for turn to end.
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
	m.PopCPP = GeneratePopulationCPP(m.CppPopulationStart, m.numCppCreated, m.Turn, m.Context)
	m.numCppCreated += m.CppPopulationStart
	m.PopVP = GeneratePopulationVP(m.VpPopulationStart, m.numVpCreated, m.Turn, m.Context)
	m.numVpCreated += m.VpPopulationStart
	go m.run(m.e)
	if m.Visualise {
		go m.vis(m.e)
	}
	if m.Logging {
		go m.logging(m.e)
	}
	return nil
}

// Data Logging process local to the model instance.
func (m *Model) logging(ec chan<- error) {
	var turn chan struct{}
	var clash bool
	var signature string
	defer func() {
		m.turnSignal.Deregister(signature)
		// Need to wipe the agent records too? -yes, probably.
	}()

Registration: // register to receive from m.turnSignal - loop until no clash with existing entry in map.
	signature = uuid()
	turn, clash = m.turnSignal.Register(signature)
	if clash {
		goto Registration
	}

	for {
		select {
		case <-m.rc: // run finished as rc channel closed!
			time.Sleep(time.Second)
			// clean up?
			return
		case <-turn:
			func() {
				reccpp := m.copyCppRecord()
				recvp := m.copyVpRecord()
				go func(rc map[string]ColourPolymorphicPrey, errCh chan<- error) {
					dir := "m.LogPath" + string(filepath.Separator) + "abmlog" + string(filepath.Separator) + m.SessionIdentifier + string(filepath.Separator) + m.timestamp
					path := dir + string(filepath.Separator) + "0" + "_vp_pop_record.dat"

					msg, err := json.MarshalIndent(rc, "", "  ")
					if err != nil {
						log.Printf("model: logging: json.Marshal failed, error: %v\n source: %s : %s : %v\n", err, m.SessionIdentifier, m.timestamp, m.Turn)
						errCh <- err
						return
					}
					var buff []byte
					out := bytes.NewBuffer(buff)
					out.Write(msg)
					output := make([]byte, 1024*10)
					n, rerr := out.Read(output)
					if n == 0 || rerr != nil {
						fmt.Println("n:", n, "rerr:", rerr.Error())
						errCh <- err
						return
					}
					err = os.MkdirAll(dir, 0777)
					if err != nil {
						log.Println(err)
						errCh <- err
						return
					}
					err = ioutil.WriteFile(path, output, 0777)
					if err != nil {
						fmt.Println(err)
						errCh <- err
						return
					}
				}(reccpp, ec)
				go func(rv map[string]VisualPredator) {
					// write map as json to file.
				}(recvp)
			}()
		}
	}
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
		go m.logging(m.e)
	}
	return nil
}

func (m *Model) run(ec chan<- error) {
	var turn chan struct{}
	var clash bool
	var signature string
	defer func() {
		m.turnSignal.Deregister(signature)
		// Need to wipe the agent records too? -yes, probably.
	}()

Registration: // register to receive from m.turnSignal - loop until no clash with existing entry in map.
	signature = uuid()
	turn, clash = m.turnSignal.Register(signature)
	if clash {
		goto Registration
	}

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
			}
			m.turn(ec) //	PROCEED WITH TURN
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
	m.turnSignal.Broadcast() // use blocking version turn ensure synchronisation
	m.Turn++
	time.Sleep(time.Millisecond * 50) // sync wiggle room
}

// Visualisation process local to the model instance.
func (m *Model) vis(ec chan<- error) {
	var turn chan struct{}
	var clash bool
	var signature string
	defer m.turnSignal.Deregister(signature)

Registration: // register to receive from m.turnSignal - loop until no clash with existing entry in map.
	signature = uuid()
	turn, clash = m.turnSignal.Register(signature)
	if clash {
		goto Registration
	}

	msg := gobr.OutMsg{Type: "render", Data: nil}
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
		case job := <-m.render:
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
			msg = gobr.OutMsg{Type: "render", Data: nil}
			//	reset draw instructions
			dl = render.DrawList{ //	reset
				CPP:       nil,
				VP:        nil,
				BG:        bg,
				CppPop:    "0",
				VpPop:     "0",
				TurnCount: "0",
			}
		case <-m.rc: //	run channel closed!
			time.Sleep(time.Second)
			return
		}
	}
}

// kill off the model and any client bound to it - internal killoff switch
func (m *Model) kill() {
	m.Stop()
	close(m.Quit)
	m.Dead = true
}
