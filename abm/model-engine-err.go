package abm

import (
	"log"
	"time"
)

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
