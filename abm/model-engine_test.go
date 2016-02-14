package abm

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/benjamin-rood/abm-cp/render"
	"github.com/benjamin-rood/gobr"
)

func TestLogMarshalling(t *testing.T) {
	// _ = "breakpoint" // godebug
	// create a new model instance, run for a limited set of turns
	// write marshalled JSON records to log files
	// read from log files, compare with expected values
	tm := newTestModel()
	go tm.ErrPrinter()
	tm.Start()
	select {
	case <-tm.rc:
		close(tm.Quit)
		return
	default:
	}
}

func newTestModel() *Model {
	tm := Model{}
	tm.running = false
	tm.Timeframe = Timeframe{}
	tm.Environment = DefaultEnvironment
	tm.Context = TestContext
	tm.timestamp = testStamp
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

func (tm *Model) testModelStart() error {
	// _ = "breakpoint" // godebug
	if tm.running {
		return errors.New("testModel: Start() failed: model already running")
	}
	tm.running = true
	if tm.RNGRandomSeed {
		rand.Seed(time.Now().UnixNano())
	} else {
		rand.Seed(tm.RNGSeedVal)
	}
	tm.PopCPP = GeneratePopulationCPP(tm.CppPopulationStart, tm.numCppCreated, tm.Turn, tm.Context, tm.timestamp)
	tm.numCppCreated += tm.CppPopulationStart
	tm.PopVP = GeneratePopulationVP(tm.VpPopulationStart, tm.numVpCreated, tm.Turn, tm.Context, tm.timestamp)
	tm.numVpCreated += tm.VpPopulationStart
	go tm.run(tm.e)
	if tm.Visualise {
		go tm.vis(tm.e)
	}
	if tm.Logging {
		go tm.log(tm.e)
	}
	return nil
}
