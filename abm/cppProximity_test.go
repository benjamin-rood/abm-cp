package abm

import (
	"errors"
	"sort"
	"testing"

	"github.com/benjamin-rood/abm-cp/geometry"
)

func TestProximitySort(t *testing.T) {
	prey := []ColourPolymorphicPrey{
		cppTesterAgent(0.0, 0.40),
		cppTesterAgent(0.35, 0.0),
		cppTesterAgent(0.0, -0.3),
		cppTesterAgent(-0.25, 0.0),
		cppTesterAgent(0.1, 0.1),
	}

	predator := vpTesterAgent(0.0, 0.0)

	want := []ColourPolymorphicPrey{
		prey[4],
		prey[3],
		prey[2],
		prey[1],
		prey[0],
	}

	for i := range want {
		want[i].δ, _ = geometry.VectorDistance(want[i].pos, predator.pos)
	}

	got := []ColourPolymorphicPrey{}
	for _, p := range prey {
		p.δ, _ = geometry.VectorDistance(p.pos, predator.pos)
		got = append(got, p)
	}

	sort.Sort(Proximity(got))

	ok, err := proxEquivalence(want, got)
	if err != nil {
		return
	}
	if !ok {
		t.Errorf("ProximitySort(got): %v, %v, %v, %v, %v \n\t\t\twant: %v, %v, %v, %v, %v\n", got[0].δ, got[1].δ, got[2].δ, got[3].δ, got[4].δ, want[0].δ, want[1].δ, want[2].δ, want[3].δ, want[4].δ)
	}
}

func proxEquivalence(p []ColourPolymorphicPrey, q []ColourPolymorphicPrey) (bool, error) {
	if len(p) != len(q) {
		return false, errors.New("input slices not of the same length")
	}
	for i := range p {
		if p[i].δ != q[i].δ {
			return false, nil
		}
	}
	return true, nil
}
