package abm

import (
	"fmt"
	"testing"
	"time"

	"github.com/benjamin-rood/abm-cp/render"
	"github.com/benjamin-rood/gobr"
)

func newTestModel() *Model {
	tm := Model{}
	tm.running = false
	tm.Timeframe = Timeframe{}
	tm.Environment = DefaultEnvironment
	tm.Context = TestContext
	tm.timestamp = fmt.Sprintf("%s", time.Now())
	tm.recordCPP = make(map[string]ColourPolymorphicPrey)
	tm.recordVP = make(map[string]VisualPredator)
	tm.Om = make(chan gobr.OutMsg)
	tm.Im = make(chan gobr.InMsg)
	tm.e = make(chan error)
	tm.Quit = make(chan struct{})
	tm.rc = make(chan struct{})
	tm.render = make(chan render.AgentRender)
	tm.turnSignal = gobr.NewSignalHub()
	return &tm
}

func newTestAgentPopulations() ([]ColourPolymorphicPrey, []VisualPredator) {
	popCpp := GeneratePopulationCPP(10, 0, 0, TestContext, testStamp)
	popVp := GeneratePopulationVP(4, 0, 0, TestContext, testStamp)

	return popCpp, popVp
}

func (tm *Model) testModelStart() {
	tm.PopCPP, tm.PopVP = newTestAgentPopulations()
	tm.run(tm.e)
	go tm.logging(tm.e)
}

func TestLogMarshalling(t *testing.T) {
	_ = "breakpoint" // godebug
	// create a new model instance, run for a limited set of turns
	// write marshalled JSON records to log files
	// read from log files, compare with expected values
	tm := newTestModel()
	go tm.ErrPrinter()
	tm.Start()
	select {
	case <-tm.rc:
		return
	default:
	}
}
