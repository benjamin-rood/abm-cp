package abm

import (
	"log"
	"time"
)

// ErrPrinter is the most basic error channel handler – it simply prints error messages.
func (m *Model) ErrPrinter() {
	for {
		select {
		case <-m.Quit:
			return
		case err := <-m.e:
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second)
			}
		}
	}
}
