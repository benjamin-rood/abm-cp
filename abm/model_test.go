package abm

import (
	"fmt"
	"testing"
	"time"

	"github.com/benjamin-rood/abm-cp/render"
	"github.com/benjamin-rood/gobr"
)

func newTestModel() (m *Model) {
	m.running = false
	m.Timeframe = Timeframe{}
	m.Environment = DefaultEnvironment
	m.Context = TestContext
	m.timestamp = fmt.Sprintf("%s", time.Now())
	m.recordCPP = make(map[string]ColourPolymorphicPrey)
	m.recordVP = make(map[string]VisualPredator)
	m.Om = make(chan gobr.OutMsg)
	m.Im = make(chan gobr.InMsg)
	m.e = make(chan error)
	m.Quit = make(chan struct{})
	m.rc = make(chan struct{})
	m.render = make(chan render.AgentRender)
	return
}

func newTestAgentPopulations(m *Model, timestamp string) {

}

func TestLogMarshalling(t *testing.T) {
	// create a new model instance, run for a limited set of turns
	// write marshalled JSON records to log files
	// read from log files, compare with expected values
	//model := abm.NewModel()

}
