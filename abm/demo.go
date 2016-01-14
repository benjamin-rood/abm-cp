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

	demoMaxCPP = 500
	demoMaxVP  = 10

	quarterpi      = 0.7853981633974483096156608458198757210492923498437764
	eigthpi        = 0.3926990816987241548078304229099378605246461749218882
	d              = 1.0
	dimensionality = 2
	cppPopSize     = 25
	vpPopSize      = 0
	vsr            = d / 4
	γ              = 1.0
	cpplife        = 30
	vplife         = 250
	vpS            = 0.1
	vpA            = 1.0
	vτ             = quarterpi
	vκ             = 0.0
	v𝛔             = 0.0
	v𝛂             = 0.0
	cppS           = 0.004
	cppA           = 1.0
	cτ             = quarterpi
	sr             = 0.004
	randomAges     = true
	fuzzy          = 0.3
	mf             = 0.05
	cφ             = 3
	cȣ             = 3
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
	// DemoViewport to be used as a baseline reference
	DemoViewport = render.Viewport{Width: 600, Height: 400}

	// DemoContext to be used as a baseline example
	DemoContext = Context{
		Bounds:                DemoEnvironment.Bounds,
		CppPopulation:         cppPopSize,
		VpPopulation:          vpPopSize,
		VpAgeing:              vpAgeing,
		VpLifespan:            vplife,
		VS:                    vpS,
		VA:                    vpA,
		VpTurn:                vτ,
		Vsr:                   vsr,
		Vγ:                    γ,
		VpReproductiveChance:  vκ,
		VsrSearchChance:       v𝛔,
		VpAttackChance:        v𝛂,
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
