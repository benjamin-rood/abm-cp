package abm

import "testing"

func TestDefaultModelCreation(t *testing.T) {
	dm := NewModel()
	go dm.Controller()
	go dm.ErrPrinter()
	dm.Start()
	select {
	case <-dm.rc:
		close(dm.Quit)
		return
	default:
	}
}
