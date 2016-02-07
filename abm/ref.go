package abm

import "github.com/benjamin-rood/abm-cp/colour"

// Default holds baseline parameters for a running model.

const (
	x = iota
	y
	z

	τ         = 6.2831853071795864769252867665590057683943387987502116 //	tau
	twoPi     = τ                                                      //	tau
	quarterpi = 0.7853981633974483096156608458198757210492923498437764
	eigthpi   = 0.3926990816987241548078304229099378605246461749218882

	d              = 1.0
	dimensionality = 2

	dCppPopStart           = 3
	dCppPopCap             = 500
	dCppAgeing             = true
	dCppLifespan           = 50
	dCppMovS               = 0.006
	dCppMovA               = 1.0
	dCppTurn               = quarterpi
	dCppSr                 = dCppMovS
	dCppGestation          = 1   // φ
	dCppSexualCost         = 1   // ȣ
	dCppReproductionChance = 0.1 // cκ
	dCppSpawnSize          = 5   // β
	dVpPopStart            = 3
	dVpPopCap              = 10
	dVpAgeing              = true
	dVpLifespan            = 250
	dVpMovS                = 0.01
	dVpMovA                = 1.0
	dVpTurn                = eigthpi
	dVsr                   = dVpMovS
	dVγ                    = 0.3
	dVγBump                = 1.2
	dVpReproductiveChance  = 1.0
	dVpSexualRequirement   = 50
	dVpGestation           = 5
	dVpSearchChance        = 1.0
	dVpAttackChance        = 1.0
	dVpColImpFactor        = 0.2
	dVpStarvationPoint     = 250
	dVpStarvation          = false
	dCppMf                 = 0.05
	dRandomAges            = true
	dRNGRandomSeed         = true
	dRNGSeedVal            = 0
	dFuzzy                 = 0.1

	tCppPopStart           = 5
	tCppPopCap             = 5
	tCppAgeing             = false
	tCppLifespan           = 1
	tCppMovS               = 0.005
	tCppMovA               = 1.0
	tCppTurn               = quarterpi
	tCppSr                 = tCppMovS
	tCppGestation          = 1
	tCppSexualCost         = 1
	tCppReproductionChance = 1.0
	tCppSpawnSize          = 1
	tVpPopStart            = 5
	tVpPopCap              = 5
	tVpAgeing              = false
	tVpLifespan            = 9999
	tVpStarvationPoint     = 9999
	tStarvation            = false
	tVpMovS                = 0.05
	tVpMovA                = 1.0
	tVpTurn                = eigthpi / 2
	tVpVsr                 = 0.2
	tVγ                    = 1.0
	tVγBump                = 1.2
	tVpReproductiveChance  = 1.0
	tVpSearchChance        = 1.0
	tVpAttackChance        = 1.0
	tVpColImpFactor        = 0.2
	tCppMf                 = 0.1
	tRandomAges            = false
	tRNGSeedVal            = 0
	tFuzzy                 = 0.1
)

var (
	// DefaultBG background for visualisation
	DefaultBG = colour.RGB{Red: 0.423529411765, Green: 0.376470588235, Blue: 0.376470588235}

	// DefaultEnvironment to be used as a baseline example
	DefaultEnvironment = Environment{
		Bounds:         []float64{d, d},
		Dimensionality: dimensionality,
		BG:             DefaultBG,
	}

	// DefaultContext to be used as a baseline example
	DefaultContext = Context{
		Bounds:                DefaultEnvironment.Bounds,
		CppPopulationStart:    dCppPopStart,
		CppPopulationCap:      dCppPopCap,
		CppAgeing:             dCppAgeing,
		CppLifespan:           dCppLifespan,
		CppS:                  dCppMovS,
		CppA:                  1.0,
		CppTurn:               dCppTurn,
		CppSr:                 dCppTurn,
		CppGestation:          dCppGestation,
		CppSexualCost:         dCppSexualCost,
		CppReproductionChance: dCppReproductionChance,
		CppSpawnSize:          dCppSpawnSize,
		VpPopulationStart:     dVpPopStart,
		VpPopulationCap:       dVpPopCap,
		VpAgeing:              dVpAgeing,
		VpLifespan:            dVpLifespan,
		VpStarvationPoint:     dVpStarvationPoint,
		VpMovS:                dVpMovS,
		VpMovA:                1.0,
		VpTurn:                dVpTurn,
		Vsr:                   dVsr,
		Vγ:                    dVγ,
		VγBump:                dVγBump,
		VpReproductionChance:  dVpReproductiveChance,
		VpSexualRequirement:   dVpSexualRequirement,
		VpGestation:           dVpGestation,
		VpSearchChance:        dVpSearchChance,
		VpAttackChance:        dVpAttackChance,
		VpColImprintFactor:    dVpColImpFactor,
		CppMutationFactor:     dCppMf,
		Starvation:            dVpStarvation,
		RandomAges:            dRandomAges,
		RNGRandomSeed:         dRNGRandomSeed,
		RNGSeedVal:            dRNGSeedVal,
		Fuzzy:                 dFuzzy,
	}

	// TestContext to be used for unit testing.
	TestContext = Context{
		Bounds:                DefaultEnvironment.Bounds,
		CppPopulationStart:    tCppPopStart,
		CppPopulationCap:      tCppPopCap,
		CppAgeing:             tCppAgeing,
		CppLifespan:           tCppLifespan,
		CppS:                  tCppMovS,
		CppA:                  1.0,
		CppTurn:               tCppTurn,
		CppSr:                 tCppTurn,
		CppGestation:          tCppGestation,
		CppSexualCost:         tCppSexualCost,
		CppReproductionChance: tCppReproductionChance,
		CppSpawnSize:          tCppSpawnSize,
		VpPopulationStart:     tVpPopStart,
		VpPopulationCap:       tVpPopCap,
		VpAgeing:              tVpAgeing,
		VpLifespan:            tVpLifespan,
		VpStarvationPoint:     tVpStarvationPoint,
		VpMovS:                tVpMovS,
		VpMovA:                1.0,
		VpTurn:                tVpTurn,
		Vsr:                   tVpVsr,
		Vγ:                    tVγ,
		VγBump:                tVγBump,
		VpReproductionChance:  tVpReproductiveChance,
		VpSearchChance:        tVpSearchChance,
		VpAttackChance:        tVpAttackChance,
		VpColImprintFactor:    tVpColImpFactor,
		CppMutationFactor:     tCppMf,
		Starvation:            tStarvation,
		RandomAges:            dRandomAges,
		RNGRandomSeed:         dRNGRandomSeed,
		RNGSeedVal:            dRNGSeedVal,
		Fuzzy:                 dFuzzy,
	}
)
