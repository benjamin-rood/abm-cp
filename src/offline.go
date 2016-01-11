package main

import (
	"math/rand"
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/abm"
	"github.com/benjamin-rood/abm-colour-polymorphism/render"
	"github.com/benjamin-rood/goio"
)

var (
	om    = make(chan goio.OutMsg)
	phase = make(chan struct{})
	view  = make(chan render.Viewport)
	ctxt  = make(chan abm.Context)
)

func writer(om <-chan goio.OutMsg) {
	for {
		select {
		case <-om:
			// jsonMsg, _ := json.MarshalIndent(msg, "", " ")
			// ioutil.WriteFile("/tmp/outMsg", jsonMsg, 0644)
			// if err := ws.WriteJSON(msg); err != nil {
			// 	log.Println("writer: failed to WriteJSON:", err)
			// }
			_ = "breakpoint" // godebug
			time.Sleep(time.Millisecond * 5)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	abm.InitModel(abm.DemoContext, abm.DemoEnvironment, om, view, phase)
	writer(om)
}
