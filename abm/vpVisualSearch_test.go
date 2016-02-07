package abm

import (
	"log"
	"testing"

	"github.com/benjamin-rood/abm-cp/colour"
)

func TestVisualSearch(t *testing.T) {
	predator := vpTesterAgent(0.0, 0.0)
	predator.vsr = 1.0
	predator.colImprint = colour.RGB{Red: 0.71, Green: 0.1, Blue: 0.39}
	prey := []ColourPolymorphicPrey{}
	for i := 1; i <= 10; i++ {
		agentA := cppTesterAgent(float64(i)*(0.01), float64(i)*(0.01))
		agentB := cppTesterAgent(float64(-i)*(0.01), float64(-i)*(0.01))
		agentA.colouration = colour.RGB{Red: float64(i) * 0.1, Green: float64(i) * 0.1, Blue: float64(i) * 0.1}
		agentB.colouration = colour.RGB{Red: 1 - float64(i)*0.1, Green: 1 - float64(i)*0.1, Blue: 1 - float64(i)*0.1}
		prey = append(prey, agentA)
		prey = append(prey, agentB)
	}

	// for i := range prey {
	// 	fmt.Printf("%v %v %v %v %v %p\n", i, prey[i].pos, prey[i].δ, prey[i].colouration, prey[i].𝛘, &prey[i])
	// }

	_ = "breakpoint" //	godebug
	want := &prey[4] // <- the best match with the least visual difference (distance) from the predator's expectation * the TestContext.VpSearchChance odds of 0.5 (50%).

	got, err := predator.PreySearch(prey, TestContext.VpSearchChance)

	switch {
	case err != nil:
		log.Fatalln("TestVisualSearch:", err)
	case got == nil:
		t.Errorf("TestVisualSearch: got = %p\n", got)
		return
	case got != want:
		t.Errorf("TestVisualSearch:\ngot: %p\n%v\nwant: %p\n%v\n", got, got.String(), want, want.String())
	}
}
