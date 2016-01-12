package abm

import (
	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/render"
)

// demo holds baseline parameters for a running model.

const (
	x = iota
	y
	z

	maxPopSize     = 1000
	quarterpi      = 0.7853981633974483096156608458198757210492923498437764
	eigthpi        = 0.3926990816987241548078304229099378605246461749218882
	d              = 1.0
	dimensionality = 2
	cppPopSize     = 30
	vpPopSize      = 0
	vsr            = d / 4
	γ              = 1.0
	cpplife        = 25
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
	sr             = 0.02
	randomAges     = true
	mf             = 0.1
	cφ             = 3
	cȣ             = 5
	cκ             = 0.1
	cβ             = 5
	vpAgeing       = false
	cppAgeing      = true
)

var (

	// DemoEnvironment to be used as a baseline example
	DemoEnvironment = Environment{
		Bounds:         []float64{d, d},
		Dimensionality: dimensionality,
		BG:             colour.RandRGB(),
	}

	DemoViewport = render.Viewport{Width: 1600, Height: 1300}

	// DemoContext to be used as a baseline example
	DemoContext = Context{
		Bounds:                DemoEnvironment.Bounds,
		CppPopulation:         cppPopSize,
		VpPopulation:          vpPopSize,
		VpAgeing:              vpAgeing,
		VpLifespan:            vplife,
		VS:                    vpS,
		VA:                    vpA,
		Vτ:                    vτ,
		Vsr:                   vsr,
		Vγ:                    γ,
		Vκ:                    vκ,
		V𝛔:                    v𝛔,
		V𝛂:                    v𝛂,
		CppAgeing:             cppAgeing,
		CppLifespan:           cpplife,
		CppS:                  cppS,
		CppA:                  cppA,
		CppTurn:               cτ,
		CppSr:                 sr,
		RandomAges:            randomAges,
		MutationFactor:        mf,
		CppGestation:          cφ,
		CppSexualCost:         cȣ,
		CppReproductiveChance: cκ,
		CppSpawnSize:          cβ,
	}
)
