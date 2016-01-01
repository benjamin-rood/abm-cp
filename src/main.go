package main

import (
	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/model"
)

const (
	d              = 1.0
	dimensionality = 2
	cppPopSize     = 25
	vpPopSize      = 0
	vsr            = d
	γ              = 1.0
	cpplife        = -1
	vplife         = -1
	vpS            = 0.0
	vpA            = 1.0
	vκ             = 0.0
	v𝛔             = 0.0
	v𝛂             = 0.0
	cppS           = 0.1
	cppA           = 1.0
	sr             = 1.0
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
	black = colour.RGB{Red: 0.0, Green: 0.0, Blue: 0.0}

	e = model.Environment{
		Bounds: []float64{d, d},
		BG:     black,
	}

	time    = model.Timeframe{Turn: 0, Phase: 0, Action: 0}
	context = model.Context{
		e,
		time,
		dimensionality,
		cppPopSize,
		vpPopSize,
		vpAge,
		vplife,
		vpS,
		vpA,
		vsr,
		γ,
		vκ,
		v𝛔,
		v𝛂,
		cppAge,
		cpplife,
		cppS,
		cppA,
		sr,
		randomAges,
		mf,
		cφ,
		cȣ,
		cκ,
		cβ,
	}
)

func main() {
}
