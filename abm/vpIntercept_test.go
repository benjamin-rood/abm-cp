package abm

import (
	"testing"

	"github.com/benjamin-rood/abm-cp/colour"
)

func TestVPIntercept(t *testing.T) {
	predator := vpTesterAgent(0, 0)
	predator.𝚯 = eigthpi
	predator.tr = eigthpi / 2
	predator.vsr = 0.5
	predator.movS = 0.05
	predator.colImprint = colour.RGB{Red: 0.71, Green: 0.1, Blue: 0.39}
	prey := []ColourPolymorphicPrey{cppTesterAgent(0.2, 0.2)}
	prey[0].colouration = colour.RGB{Red: 0.5, Green: 0.5, Blue: 0.5} //  close enough to be recognised.

	want := 0.83841
	got := predator.𝚯
	for {
		target, _ := predator.PreySearch(prey, 1.0)
		got = predator.𝚯
		// fmt.Println(predator.pos, got)
		if target != nil {
			break
		}
	}

	if got != want {
		t.Errorf("TestVPIntercept: got = %v, want = %v\n", got, want)
	}

}
