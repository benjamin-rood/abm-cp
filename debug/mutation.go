package main

import (
	"math/rand"
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/abm"
	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
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

	context = abm.Context{
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

const (
	quarterpi      = 0.7853981633974483096156608458198757210492923498437764
	eigthpi        = 0.3926990816987241548078304229099378605246461749218882
	d              = 1.0
	dimensionality = 2
	cppPopSize     = 30
	vpPopSize      = 0
	vsr            = d / 4
	γ              = 1.0
	cpplife        = 100
	vplife         = -1
	vpS            = 0.0
	vpA            = 1.0
	vτ             = quarterpi
	vκ             = 0.0
	v𝛔             = 0.0
	v𝛂             = 0.0
	cppS           = 0.02
	cppA           = 1.0
	cτ             = quarterpi
	sr             = 0.02
	randomAges     = true
	mf             = 0.1
	cφ             = 5
	cȣ             = 5
	cκ             = 0.2
	cβ             = 10
	vpAge          = false
	cppAge         = true
)

func main() {
	rand.Seed(time.Now().UnixNano())
	agents := abm.GeneratePopulation(1, context)
	for {
		for i := range agents {
			abm.Mutation(&agents[i], context.Mf)
		}
	}
}
