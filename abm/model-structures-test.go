package abm

import "testing"

func TestDefaultModelCreation(t *testing.T) {
	_ = "breakpoint" // godebug
	dm := NewModel()
	go dm.Controller()
	go dm.ErrPrinter()
	dm.Start()
	select {
	case <-dm.rq:
		close(dm.Quit)
		return
	default:
	}
}
