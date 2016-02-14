package abm

//
// import (
// 	"errors"
// 	"sort"
// 	"testing"
//
// 	"github.com/benjamin-rood/abm-cp/colour"
// 	"github.com/benjamin-rood/abm-cp/geometry"
// )
//
// func TestColourDistanceSort(t *testing.T) {
// 	predator := VisualPredator{}
// 	predator.τ = colour.RGB{Red: 0.5, Green: 0.5, Blue: 0.5}
//
// 	prey := []*ColourPolymorphicPrey{}
// 	for i := 5; i > 0; i-- {
// 		agent := &ColourPolymorphicPrey{}
// 		agent.colouration = colour.RGB{Red: float64(i) * 0.2, Green: float64(i) * 0.2, Blue: float64(i) * 0.2}
// 		agent.𝛘 = colour.RGBDistance(agent.colouration, predator.τ)
// 		prey = append(prey, agent)
// 	}
//
// 	// fmt.Println("Before sorting:")
// 	// for i, p := range prey {
// 	// 	fmt.Println(i, p.colouration, p.𝛘)
// 	// }
//
// 	want := prey
// 	want[0].colouration = colour.RGB{Red: 0.4, Green: 0.4, Blue: 0.4}
// 	want[0].𝛘 = colour.RGBDistance(want[0].colouration, predator.τ)
// 	want[1].colouration = colour.RGB{Red: 0.6, Green: 0.6, Blue: 0.6}
// 	want[1].𝛘 = colour.RGBDistance(want[1].colouration, predator.τ)
// 	want[2].colouration = colour.RGB{Red: 0.2, Green: 0.2, Blue: 0.2}
// 	want[2].𝛘 = colour.RGBDistance(want[2].colouration, predator.τ)
// 	want[3].colouration = colour.RGB{Red: 0.8, Green: 0.8, Blue: 0.8}
// 	want[3].𝛘 = colour.RGBDistance(want[3].colouration, predator.τ)
// 	want[4].colouration = colour.RGB{Red: 1.0, Green: 1.0, Blue: 1.0}
// 	want[4].𝛘 = colour.RGBDistance(want[4].colouration, predator.τ)
//
// 	sort.Sort(byColourDifferentiation(prey))
//
// 	// fmt.Println("After sorting:")
// 	// for i, p := range prey {
// 	// 	fmt.Println(i, p.colouration, p.𝛘)
// 	// }
//
// 	ok, err := colourDiffEquivalence(want, prey)
// 	if err != nil {
// 		return
// 	}
// 	if !ok {
// 		t.Errorf("VisDistSort(got): %v, %v, %v, %v, %v \n\t\t\twant: %v, %v, %v, %v, %v\n", prey[0].𝛘, prey[1].𝛘, prey[2].𝛘, prey[3].𝛘, prey[4].𝛘, want[0].𝛘, want[1].𝛘, want[2].𝛘, want[3].𝛘, want[4].𝛘)
// 	}
//
// 	predator.τ = colour.RGB{Red: 0.31, Green: 0.79, Blue: 0.01}
// 	prey = []*ColourPolymorphicPrey{}
// 	for i := 0; i < 10; i++ {
// 		agent := &ColourPolymorphicPrey{}
// 		agent.colouration = colour.RGB{Red: float64(i) * 0.1, Green: float64(i) * 0.1, Blue: float64(i) * 0.1}
// 		agent.𝛘 = colour.RGBDistance(agent.colouration, predator.τ)
// 		prey = append(prey, agent)
// 	}
//
// 	// fmt.Println("Before sorting:")
// 	// for i, p := range prey {
// 	// 	fmt.Println(i, p.colouration, p.𝛘)
// 	// }
//
// 	copy := []*ColourPolymorphicPrey{}
// 	for _, p := range prey {
// 		copy = append(copy, p)
// 	}
//
// 	want = []*ColourPolymorphicPrey{}
// 	want = append(want, copy[3])
// 	want = append(want, copy[4])
// 	want = append(want, copy[2])
// 	want = append(want, copy[5])
// 	want = append(want, copy[1])
// 	want = append(want, copy[6])
// 	want = append(want, copy[0])
// 	want = append(want, copy[7])
// 	want = append(want, copy[8])
// 	want = append(want, copy[9])
//
// 	sort.Sort(byColourDifferentiation(prey))
//
// 	// fmt.Println("After sorting:")
// 	// for i, p := range prey {
// 	// 	fmt.Println(i, p.colouration, p.𝛘)
// 	// }
//
// 	ok, err = colourDiffEquivalence(want, prey)
// 	if err != nil {
// 		return
// 	}
// 	if !ok {
// 		t.Errorf("VisDistSort(got):\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\nwant:\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n", prey[0].𝛘, prey[1].𝛘, prey[2].𝛘, prey[3].𝛘, prey[4].𝛘, prey[5].𝛘, prey[6].𝛘, prey[7].𝛘, prey[8].𝛘, prey[9].𝛘, want[0].𝛘, want[1].𝛘, want[2].𝛘, want[3].𝛘, want[4].𝛘, want[5].𝛘, want[6].𝛘, want[7].𝛘, want[8].𝛘, want[9].𝛘)
// 	}
//
// }
//
// func colourDiffEquivalence(p []*ColourPolymorphicPrey, q []*ColourPolymorphicPrey) (bool, error) {
// 	if len(p) != len(q) {
// 		return false, errors.New("input slices not of the same length")
// 	}
// 	for i := range p {
// 		if p[i].𝛘 != q[i].𝛘 {
// 			return false, nil
// 		}
// 	}
// 	return true, nil
// }
//
// func TestProximitySort(t *testing.T) {
// 	prey := []*ColourPolymorphicPrey{
// 		newCppTesterAgent(0.0, 0.40),
// 		newCppTesterAgent(0.35, 0.0),
// 		newCppTesterAgent(0.0, -0.3),
// 		newCppTesterAgent(-0.25, 0.0),
// 		newCppTesterAgent(0.1, 0.1),
// 	}
//
// 	predator := vpTesterAgent(0.0, 0.0)
//
// 	want := []*ColourPolymorphicPrey{
// 		prey[4],
// 		prey[3],
// 		prey[2],
// 		prey[1],
// 		prey[0],
// 	}
//
// 	for i := range want {
// 		want[i].δ, _ = geometry.VectorDistance(want[i].pos, predator.pos)
// 	}
//
// 	got := []*ColourPolymorphicPrey{}
// 	for _, p := range prey {
// 		p.δ, _ = geometry.VectorDistance(p.pos, predator.pos)
// 		got = append(got, p)
// 	}
//
// 	sort.Sort(byProximity(got))
//
// 	ok, err := proxEquivalence(want, got)
// 	if err != nil {
// 		return
// 	}
// 	if !ok {
// 		t.Errorf("ProximitySort(got): %v, %v, %v, %v, %v \n\t\t\twant: %v, %v, %v, %v, %v\n", got[0].δ, got[1].δ, got[2].δ, got[3].δ, got[4].δ, want[0].δ, want[1].δ, want[2].δ, want[3].δ, want[4].δ)
// 	}
// }
//
// func proxEquivalence(p []*ColourPolymorphicPrey, q []*ColourPolymorphicPrey) (bool, error) {
// 	if len(p) != len(q) {
// 		return false, errors.New("input slices not of the same length")
// 	}
// 	for i := range p {
// 		if p[i].δ != q[i].δ {
// 			return false, nil
// 		}
// 	}
// 	return true, nil
// }
