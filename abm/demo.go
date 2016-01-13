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

	maxPopSize     = 200
	quarterpi      = 0.7853981633974483096156608458198757210492923498437764
	eigthpi        = 0.3926990816987241548078304229099378605246461749218882
	d              = 1.0
	dimensionality = 2
	cppPopSize     = 100
	vpPopSize      = 0
	vsr            = d / 4
	γ              = 1.0
	cpplife        = 25
	vplife         = 250
	vpS            = 0.1
	vpA            = 1.0
	vτ             = quarterpi
	vκ             = 0.0
	v𝛔             = 0.0
	v𝛂             = 0.0
	cppS           = 0.001
	cppA           = 1.0
	cτ             = quarterpi
	sr             = 0.001
	randomAges     = true
	fuzzy          = 0.3
	mf             = 0.1
	cφ             = 3
	cȣ             = 3
	cκ             = 0.1
	cβ             = 5
	vpAgeing       = true
	cppAgeing      = true
)

var (
	// DemoEnvironment to be used as a baseline example
	DemoEnvironment = Environment{
		Bounds:         []float64{d, d},
		Dimensionality: dimensionality,
		BG:             colour.RandRGB(),
	}

	DemoViewport = render.Viewport{Width: 1440, Height: 900}

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
		MutationFactor:        mf,
		CppGestation:          cφ,
		CppSexualCost:         cȣ,
		CppReproductiveChance: cκ,
		CppSpawnSize:          cβ,
		RandomAges:            randomAges,
		Fuzzy:                 fuzzy,
	}
)
