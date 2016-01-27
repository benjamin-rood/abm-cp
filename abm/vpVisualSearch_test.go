package abm

import (
	"testing"

	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
)

func TestVisualSearch(t *testing.T) {
	predator := vpTesterAgent(0.0, 0.0)
	predator.vsr = 1.0
	predator.colImprint = colour.RGB{Red: 0.31, Green: 0.79, Blue: 0.01}
	prey := []ColourPolymorphicPrey{}
	for i := 1; i < 10; i++ {
		agentA := cppTesterAgent(float64(i)*(0.2), float64(i)*(0.2))
		agentB := cppTesterAgent(float64(-i)*(0.2), float64(-i)*(0.2))
		agentA.colouration = colour.RGB{Red: float64(i) * 0.1, Green: float64(i) * 0.1, Blue: float64(i) * 0.1}
		agentB.colouration = colour.RGB{Red: float64(i*2) * 0.1, Green: float64(i*2) * 0.1, Blue: float64(i*2) * 0.1}
		prey = append(prey, agentA)
		prey = append(prey, agentB)
	}

	// for i := range prey {
	// 	fmt.Printf("%d %v %v %v %p\n", i, prey[i].pos, prey[i].δ, prey[i].𝛘, &prey[i])
	// }

	want := &prey[1] // <- the best match with the least visual difference (distance) from the predator's expectation * the TestContext.VpSearchChance odds of 0.5 (50%).
	got, _ := predator.PreySearch(prey, TestContext.VpSearchChance)
	if got != want {
		t.Errorf("TestVisualSearch:\ngot: %v\nwant: %v\n", got.String(), want.String())
	}
}
