package abm

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/benjamin-rood/abm-cp/calc"
	"github.com/benjamin-rood/abm-cp/colour"
	"github.com/benjamin-rood/abm-cp/geometry"
)

func TestCPPComparitorSort(t *testing.T) {
	rand.Seed(0)
	predator := vpTesterAgent(0.0, 0.0)
	predator.vsr = 1.0
	predator.τ = colour.RGB{Red: 0.71, Green: 0.1, Blue: 0.39}

	prey := []ColourPolymorphicPrey{}

	for i := 1; i <= 10; i++ {
		agentA := cpPreyTesterAgent(calc.RandFloatIn(-1, 1), calc.RandFloatIn(-1, 1))
		agentB := cpPreyTesterAgent(calc.RandFloatIn(-1, 1), calc.RandFloatIn(-1, 1))
		agentA.colouration = colour.RGB{Red: rand.Float64(), Green: rand.Float64(), Blue: rand.Float64()}
		agentB.colouration = colour.RGB{Red: rand.Float64(), Green: rand.Float64(), Blue: rand.Float64()}
		prey = append(prey, agentA)
		prey = append(prey, agentB)
	}

	f := func(v geometry.Vector) func(geometry.Vector) float64 {
		return func(u geometry.Vector) float64 {
			dist, _ := geometry.VectorDistance(v, u)
			return dist
		}
	}(predator.pos)

	compers := []compCPP{}

	for i := range prey {
		// fmt.Printf("%v\t%v\t%p\n", i, prey[i].pos, &prey[i])
		a := compCPP{f, &prey[i]}
		compers = append(compers, a)
	}

	sort.Sort(byComparitor(compers))

	// for i := range compers {
	// 	fmt.Printf("%v\t%v\t%p\t%v\n", i, compers[i].pos, compers[i].ColourPolymorphicPrey, f(compers[i].pos))
	// }

	want := fmt.Sprintf("%p", &prey[19])
	got := fmt.Sprintf("%p", &(*compers[0].ColourPolymorphicPrey))

	if want != got {
		t.Errorf("want:\n%q\ngot:\n%q\n", want, got)
	}
}

func TestOptimalAttackVectorSort(t *testing.T) {
	rand.Seed(0)
	predator := vpTesterAgent(0.0, 0.0)
	predator.vsr = 1.0
	predator.τ = colour.RGB{Red: 0.71, Green: 0.1, Blue: 0.39}

	prey := []ColourPolymorphicPrey{}

	for i := 1; i <= 10; i++ {
		agentA := cpPreyTesterAgent(calc.RandFloatIn(-1, 1), calc.RandFloatIn(-1, 1))
		agentB := cpPreyTesterAgent(calc.RandFloatIn(-1, 1), calc.RandFloatIn(-1, 1))
		agentA.colouration = colour.RGB{Red: rand.Float64(), Green: rand.Float64(), Blue: rand.Float64()}
		agentB.colouration = colour.RGB{Red: rand.Float64(), Green: rand.Float64(), Blue: rand.Float64()}
		prey = append(prey, agentA)
		prey = append(prey, agentB)
	}

	c := math.Pow(2, 3)
	f := visualSignalStrength(c)
	optimals := []visualRecognition{}
	for i := range prey {
		𝛘 := colour.RGBDistance(predator.τ, prey[i].colouration)
		δ, _ := geometry.VectorDistance(predator.pos, prey[i].pos)
		// fmt.Printf("%v\t%v\t%v\t%p\n", i, 𝛘, δ, &prey[i])
		a := visualRecognition{δ, 𝛘, f, c, &prey[i]}
		optimals = append(optimals, a)
	}

	sort.Sort(byOptimalAttackVector(optimals))

	// for i := range optimals {
	// 	fmt.Printf("%v\t%v\t%v\t%p\t%v\n", i, optimals[i].𝛘, optimals[i].δ, optimals[i].ColourPolymorphicPrey, f(optimals[i].𝛘)-optimals[i].δ)
	// }

	want := fmt.Sprintf("%p", &prey[19])
	got := fmt.Sprintf("%p", &(*optimals[0].ColourPolymorphicPrey))

	if want != got {
		t.Errorf("want:\n%q\ngot:\n%q\n", want, got)
	}
}
