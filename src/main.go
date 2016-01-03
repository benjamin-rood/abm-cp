package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/abm"
	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/render"
)

const (
	quarterpi      = 0.7853981633974483096156608458198757210492923498437764
	eigthpi        = 0.3926990816987241548078304229099378605246461749218882
	d              = 1.0
	dimensionality = 2
	cppPopSize     = 1
	vpPopSize      = 0
	vsr            = d / 4
	γ              = 1.0
	cpplife        = -1
	vplife         = -1
	vpS            = 0.0
	vpA            = 1.0
	vτ             = quarterpi
	vκ             = 0.0
	v𝛔             = 0.0
	v𝛂             = 0.0
	cppS           = 0.01
	cppA           = 1.0
	cτ             = quarterpi
	sr             = d / 8
	randomAges     = false
	mf             = 0.5
	cφ             = 5
	cȣ             = 2
	cκ             = 1.0
	cβ             = 5
	vpAge          = false
	cppAge         = false
)

var (
	black = colour.Black
	white = colour.White

	e = abm.Environment{
		Bounds:         []float64{d, d},
		Dimensionality: dimensionality,
		BG:             black,
	}

	timeframe = abm.Timeframe{Turn: 0, Phase: 0, Action: 0}
	context   = abm.Context{
		e.Bounds,
		cppPopSize,
		vpPopSize,
		vpAge,
		vplife,
		vpS,
		vpA,
		vτ,
		vsr,
		γ,
		vκ,
		v𝛔,
		v𝛂,
		cppAge,
		cpplife,
		cppS,
		cppA,
		cτ,
		sr,
		randomAges,
		mf,
		cφ,
		cȣ,
		cκ,
		cβ,
	}
)

func cppRBB(pop []abm.ColourPolymorphicPrey, queue chan<- render.AgentRender) {
	for {
		for i := 0; i < len(pop); i++ {
			c := &pop[i]
			𝚯 := calc.RandFloatIn(-c.Rτ, c.Rτ)
			c.Turn(𝚯)
			c.Move()
			c.Log()
			queue <- c.GetDrawInfo()
			time.Sleep(time.Second * 10)
		}
	}
}

func runningModel(m abm.Model, queue chan<- render.AgentRender) {
	for {
		cppRBB(m.PopCPP, queue)
	}
}

func initModel(context abm.Context) (m abm.Model) {
	m.PopCPP = abm.GeneratePopulation(cppPopSize, context)
	m.DefinitionCPP = []string{"mover"}
	m.Timeframe = timeframe
	m.Environment = e
	m.Context = context
	return
}

func visualiseModel(view <-chan render.Viewport, queue <-chan render.AgentRender, out chan<- render.Msg) {
	v := render.Viewport{"Viewport", 640, 480}
	tick := time.Tick(time.Second)
	var msg render.Msg
	msg.Type = "Render"
	for {
		select {
		case job := <-queue:
			job = render.TranslateToViewport(job, v)
			msg.CPP = append(msg.CPP, job)
			out <- msg
			msg = render.Msg{"Render", nil, nil} // reset msg contents
		case <-tick:
			// msg = render.Msg{} // reset msg contents
		case v = <-view:
		}
	}
}

func writer(out <-chan render.Msg) {
	for {
		select {
		case msg := <-out:
			m, err := json.MarshalIndent(&msg, "", "  ")
			if err != nil {
				log.Fatalf("writer: failed when trying to marshal: %q", err)
			}
			err = ioutil.WriteFile("/tmp/dat1", m, 0644)
		}
	}
}

func main() {
	simple := initModel(context)
	q := make(chan render.AgentRender)
	viz := make(chan render.Viewport)
	renderOut := make(chan render.Msg)
	go runningModel(simple, q)
	go visualiseModel(viz, q, renderOut)
	writer(renderOut)
}
