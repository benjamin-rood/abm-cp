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

	abmlogPath = "abmlog"

	dCpPreyPopStart           = 3
	dCpPreyPopCap             = 500
	dCpPreyAgeing             = true
	dCpPreyLifespan           = 50
	dCpPreyMovS               = 0.006
	dCpPreyMovA               = 1.0
	dCpPreyTurn               = quarterpi
	dCpPreySr                 = dCpPreyMovS
	dCpPreyGestation          = 1   // φ
	dCpPreySexualCost         = 1   // ȣ
	dCpPreyReproductionChance = 0.1 // cκ
	dCpPreySpawnSize          = 5   // β
	dVpPopStart               = 3
	dVpPopCap                 = 10
	dVpAgeing                 = true
	dVpLifespan               = 250
	dVpMovS                   = 0.01
	dVpMovA                   = 1.0
	dVpTurn                   = eigthpi
	dVpVsr                      = dVpMovS
	dV𝛄                       = 0.3
	dVpV𝛄Bump                   = 1.2
	dVpReproductiveChance     = 1.0
	dVpSexualRequirement      = 50
	dVpGestation              = 5
	dVpSearchChance           = 1.0
	dVpAttackChance           = 1.0
	dVpColAdaptationFactor    = 0.2
	dVpStarvationPoint        = 250
	dVpStarvation             = false
	dCpPreyMf                 = 0.05
	dRandomAges               = true
	dRNGRandomSeed            = true
	dRNGSeedVal               = 0
	dFuzzy                    = 0.1
	dLogging                  = false
	dVisualise                = true
	dLimitDuration            = false
	dSessionIdentifier        = "DefaultConditionParamsSession"

	testStamp = "TESTING ONLY"

	tCpPreyPopStart           = 25
	tCpPreyPopCap             = 100
	tCpPreyAgeing             = false
	tCpPreyLifespan           = 1
	tCpPreyMovS               = 0.005
	tCpPreyMovA               = 1.0
	tCpPreyTurn               = quarterpi
	tCpPreySr                 = tCpPreyMovS
	tCpPreyGestation          = 1
	tCpPreySexualCost         = 1
	tCpPreyReproductionChance = 1.0
	tCpPreySpawnSize          = 1
	tVpPopStart               = 5
	tVpPopCap                 = 5
	tVpAgeing                 = false
	tVpLifespan               = 9999
	tVpStarvationPoint        = 9999
	tVpStarvation             = false
	tVpMovS                   = 0.2
	tVpMovA                   = 1.0
	tVpTurn                   = eigthpi / 2
	tVpVpVsr                    = 0.2
	tV𝛄                       = 0.2
	tVpV𝛄Bump                   = 1.2
	tVpReproductiveChance     = 1.0
	tVpSearchChance           = 1.0
	tVpAttackChance           = 1.0
	tVpColAdaptationFactor    = 0.2
	tCpPreyMf                 = 0.1
	tRNGRandomSeed            = false
	tRandomAges               = false
	tRNGSeedVal               = 0
	tFuzzy                    = 0.1
	tLogging                  = true
	tLogFreq                  = 0 // write every turn
	tUseCustomLogPath         = false
	tCustomLogPath            = ""
	tLimitDuration            = true
	tFixedDuration            = 10 // two turns only
	tSessionIdentifier        = "TestConditionParamsSession"
	tVpBaseAttackGain                      = 20
	tVpVbε                      = 3
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

	// DefaultConditionParams to be used as a baseline example
	DefaultConditionParams = ConditionParams{
		Environment:              DefaultEnvironment,
		CpPreyPopulationStart:    dCpPreyPopStart,
		CpPreyPopulationCap:      dCpPreyPopCap,
		CpPreyAgeing:             dCpPreyAgeing,
		CpPreyLifespan:           dCpPreyLifespan,
		CpPreyS:                  dCpPreyMovS,
		CpPreyA:                  1.0,
		CpPreyTurn:               dCpPreyTurn,
		CpPreySr:                 dCpPreyTurn,
		CpPreyGestation:          dCpPreyGestation,
		CpPreySexualCost:         dCpPreySexualCost,
		CpPreyReproductionChance: dCpPreyReproductionChance,
		CpPreySpawnSize:          dCpPreySpawnSize,
		VpPopulationStart:        dVpPopStart,
		VpPopulationCap:          dVpPopCap,
		VpAgeing:                 dVpAgeing,
		VpLifespan:               dVpLifespan,
		VpStarvationPoint:        dVpStarvationPoint,
		VpMovS:                   dVpMovS,
		VpMovA:                   1.0,
		VpTurn:                   dVpTurn,
		VpVsr:                      dVpVsr,
		VpVb𝛄:                      dV𝛄,
		VpV𝛄Bump:                   dVpV𝛄Bump,
		VpReproductionChance:     dVpReproductiveChance,
		VpSexualRequirement:      dVpSexualRequirement,
		VpGestation:              dVpGestation,
		VpSearchChance:           dVpSearchChance,
		VpAttackChance:           dVpAttackChance,
		VpCaf:                    dVpColAdaptationFactor,
		CpPreyMutationFactor:     dCpPreyMf,
		VpStarvation:             dVpStarvation,
		RandomAges:               dRandomAges,
		RNGRandomSeed:            dRNGRandomSeed,
		RNGSeedVal:               dRNGSeedVal,
		Fuzzy:                    dFuzzy,
	}

	// TestConditionParams to be used for unit testing.
	TestConditionParams = ConditionParams{
		Environment:              DefaultEnvironment,
		CpPreyPopulationStart:    tCpPreyPopStart,
		CpPreyPopulationCap:      tCpPreyPopCap,
		CpPreyAgeing:             tCpPreyAgeing,
		CpPreyLifespan:           tCpPreyLifespan,
		CpPreyS:                  tCpPreyMovS,
		CpPreyA:                  1.0,
		CpPreyTurn:               tCpPreyTurn,
		CpPreySr:                 tCpPreyTurn,
		CpPreyGestation:          tCpPreyGestation,
		CpPreySexualCost:         tCpPreySexualCost,
		CpPreyReproductionChance: tCpPreyReproductionChance,
		CpPreySpawnSize:          tCpPreySpawnSize,
		VpPopulationStart:        tVpPopStart,
		VpPopulationCap:          tVpPopCap,
		VpAgeing:                 tVpAgeing,
		VpLifespan:               tVpLifespan,
		VpStarvationPoint:        tVpStarvationPoint,
		VpMovS:                   tVpMovS,
		VpMovA:                   1.0,
		VpTurn:                   tVpTurn,
		VpVsr:                      tVpVpVsr,
		VpVb𝛄:                      tV𝛄,
		VpV𝛄Bump:                   tVpV𝛄Bump,
		VpVbε:                      tVpVbε,
		VpBaseAttackGain:                      tVpBaseAttackGain,
		VpReproductionChance:     tVpReproductiveChance,
		VpSearchChance:           tVpSearchChance,
		VpAttackChance:           tVpAttackChance,
		VpCaf:                    tVpColAdaptationFactor,
		CpPreyMutationFactor:     tCpPreyMf,
		VpStarvation:             tVpStarvation,
		RandomAges:               tRandomAges,
		RNGRandomSeed:            tRNGRandomSeed,
		RNGSeedVal:               tRNGSeedVal,
		Fuzzy:                    tFuzzy,
		Logging:                  tLogging,
		LogFreq:                  tLogFreq,
		UseCustomLogPath:         tUseCustomLogPath,
		CustomLogPath:            tCustomLogPath,
		LimitDuration:            tLimitDuration,
		FixedDuration:            tFixedDuration,
		SessionIdentifier:        tSessionIdentifier,
	}
)
