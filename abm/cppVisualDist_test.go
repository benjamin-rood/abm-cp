package abm

import (
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
)

func TestVisualDistSort(t *testing.T) {
	predator := VisualPredator{}
	predator.colImprint = colour.RGB{Red: 0.5, Green: 0.5, Blue: 0.5}

	prey := []ColourPolymorphicPrey{}
	for i := 5; i > 0; i-- {
		agent := ColourPolymorphicPrey{}
		agent.colouration = colour.RGB{Red: float64(i) * 0.2, Green: float64(i) * 0.2, Blue: float64(i) * 0.2}
		agent.𝛘 = colour.RGBDistance(agent.colouration, predator.colImprint)
		prey = append(prey, agent)
	}

	want := prey
	want[0].colouration = colour.RGB{Red: 0.4, Green: 0.4, Blue: 0.4}
	want[0].𝛘 = colour.RGBDistance(want[0].colouration, predator.colImprint)
	want[1].colouration = colour.RGB{Red: 0.6, Green: 0.6, Blue: 0.6}
	want[1].𝛘 = colour.RGBDistance(want[1].colouration, predator.colImprint)
	want[2].colouration = colour.RGB{Red: 0.2, Green: 0.2, Blue: 0.2}
	want[2].𝛘 = colour.RGBDistance(want[2].colouration, predator.colImprint)
	want[3].colouration = colour.RGB{Red: 0.8, Green: 0.8, Blue: 0.8}
	want[3].𝛘 = colour.RGBDistance(want[3].colouration, predator.colImprint)
	want[4].colouration = colour.RGB{Red: 1.0, Green: 1.0, Blue: 1.0}
	want[4].𝛘 = colour.RGBDistance(want[4].colouration, predator.colImprint)

	sort.Sort(VisualDifference(prey))

	ok, err := visualDiffEquivalence(want, prey)
	if err != nil {
		return
	}
	if !ok {
		t.Errorf("VisDistSort(got): %v, %v, %v, %v, %v \n\t\t\twant: %v, %v, %v, %v, %v\n", prey[0].𝛘, prey[1].𝛘, prey[2].𝛘, prey[3].𝛘, prey[4].𝛘, want[0].𝛘, want[1].𝛘, want[2].𝛘, want[3].𝛘, want[4].𝛘)
	}

	predator.colImprint = colour.RGB{Red: 0.31, Green: 0.79, Blue: 0.01}
	prey = []ColourPolymorphicPrey{}
	for i := 0; i < 10; i++ {
		agent := ColourPolymorphicPrey{}
		agent.colouration = colour.RGB{Red: float64(i) * 0.1, Green: float64(i) * 0.1, Blue: float64(i) * 0.1}
		agent.𝛘 = colour.RGBDistance(agent.colouration, predator.colImprint)
		fmt.Println(i, agent.𝛘)
		prey = append(prey, agent)
	}

	copy := []ColourPolymorphicPrey{}
	for _, p := range prey {
		copy = append(copy, p)
	}

	want = []ColourPolymorphicPrey{}
	want = append(want, copy[4])
	want = append(want, copy[3])
	want = append(want, copy[5])
	want = append(want, copy[2])
	want = append(want, copy[6])
	want = append(want, copy[1])
	want = append(want, copy[7])
	want = append(want, copy[0])
	want = append(want, copy[8])
	want = append(want, copy[9])

	sort.Sort(VisualDifference(prey))

	ok, err = visualDiffEquivalence(want, prey)
	if err != nil {
		return
	}
	if !ok {
		t.Errorf("VisDistSort(got):\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\nwant:\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n", prey[0].𝛘, prey[1].𝛘, prey[2].𝛘, prey[3].𝛘, prey[4].𝛘, prey[5].𝛘, prey[6].𝛘, prey[7].𝛘, prey[8].𝛘, prey[9].𝛘, want[0].𝛘, want[1].𝛘, want[2].𝛘, want[3].𝛘, want[4].𝛘, want[5].𝛘, want[6].𝛘, want[7].𝛘, want[8].𝛘, want[9].𝛘)
	}

}

func visualDiffEquivalence(p []ColourPolymorphicPrey, q []ColourPolymorphicPrey) (bool, error) {
	if len(p) != len(q) {
		return false, errors.New("input slices not of the same length")
	}
	for i := range p {
		if p[i].𝛘 != q[i].𝛘 {
			return false, nil
		}
	}
	return true, nil
}
