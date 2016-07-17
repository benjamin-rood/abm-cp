package abm

import "testing"

func TestDefaultModelCreation(t *testing.T) {
	_ = "breakpoint" // godebug
	dm := NewModel()
	go dm.Controller()
	go dm.ErrPrinter()
	dm.Start()
	select {
	case <-dm.halt:
		close(dm.Quit)
		return
	default:
	}
}
