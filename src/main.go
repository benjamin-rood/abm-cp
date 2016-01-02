package main

import (
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
	cppPopSize     = 25
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
	cppS           = 0.1
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
		Bounds: []float64{d, d},
		BG:     black,
	}

	time    = abm.Timeframe{Turn: 0, Phase: 0, Action: 0}
	context = abm.Context{
		e,
		time,
		dimensionality,
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

func cppRBB(pop []abm.ColourPolymorphicPrey, queue <-chan render.AgentRender) {
	for {
		for i := 0; i < len(pop); i++ {
			c := &pop[i]
			𝚯 := calc.RandFloatIn(-c.Rτ, c.Rτ)
			c.Turn(𝚯)
			c.Move()
			c.Log()
		}
	}
}

func runningModel(m abm.Model, context chan<- abm.Context) {

}

func initModel(m abm.Model, context abm.Context) {
	m.PopCPP = abm.GeneratePopulation(25, context)
	m.DefinitionCPP = []string{"mover"}
}

func visualiseModel(view chan<- render.Viewport, queue chan<- render.AgentRender, render <-chan render.Msg) {

}

func main() {
}
