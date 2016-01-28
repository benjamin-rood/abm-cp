package abm

import (
	"fmt"
	"log"
	"testing"

	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
)

func TestVisualSearch(t *testing.T) {
	predator := vpTesterAgent(0.0, 0.0)
	predator.vsr = 1.0
	predator.colImprint = colour.RGB{Red: 0.5, Green: 0.5, Blue: 0.5}
	prey := []ColourPolymorphicPrey{}
	for i := 1; i <= 5; i++ {
		agentA := cppTesterAgent(float64(i)*(0.1), float64(i)*(0.1))
		agentB := cppTesterAgent(float64(-i)*(0.1), float64(-i)*(0.1))
		agentA.colouration = colour.RGB{Red: float64(i) * 0.1, Green: float64(i) * 0.1, Blue: float64(i) * 0.1}
		agentB.colouration = colour.RGB{Red: float64(i*2) * 0.1, Green: float64(i*2) * 0.1, Blue: float64(i*2) * 0.1}
		prey = append(prey, agentA)
		prey = append(prey, agentB)
	}

	for i := range prey {
		fmt.Printf("%v %v %v %v %p\n", i, prey[i].pos, prey[i].δ, prey[i].𝛘, &prey[i])
	}

	_ = "breakpoint" //	godebug
	want := &prey[5] // <- the best match with the least visual difference (distance) from the predator's expectation * the TestContext.VpSearchChance odds of 0.5 (50%).
	fmt.Println(predator.pos, predator.𝚯)
	got, err := predator.PreySearch(prey, TestContext.VpSearchChance)
	if err != nil {
		fmt.Println(err)
	}
	if got == nil {
		t.Errorf("TestVisualSearch: got = %p\n", got)
		return
	}
	fmt.Printf("got: %v %v %v %p\n", got.pos, got.δ, got.𝛘, got)
	if err != nil {
		log.Fatalln(err)
	}
	if got != want {
		t.Errorf("TestVisualSearch:\ngot: %v\nwant: %v\n", got.String(), want.String())
	}
}
