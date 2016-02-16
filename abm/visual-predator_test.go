package abm

import (
	"math/rand"
	"testing"

	"github.com/benjamin-rood/abm-cp/calc"
	"github.com/benjamin-rood/abm-cp/colour"
)

func TestVPIntercept(t *testing.T) {
	predator := vpTesterAgent(0, 0)
	predator.𝚯 = eigthpi
	predator.tr = eigthpi / 2
	predator.vsr = 0.5
	predator.movS = 0.1
	predator.τ = colour.RGB{Red: 0.71, Green: 0.1, Blue: 0.39}
	prey := []ColourPolymorphicPrey{cppTesterAgent(0.05, 0.05)}
	prey[0].colouration = colour.RGB{Red: 0.6, Green: 0.2, Blue: 0.4} //  close enough to be recognised.

	want := 0.3927
	target, _, _ := predator.PreySearch(prey, 1.0)
	if target == nil {
		t.Errorf("No target found.")
	}
	got := predator.𝚯
	// fmt.Println(predator.pos, got)

	got = calc.ToFixed(predator.𝚯, 5)

	if got != want {
		t.Errorf("TestVPIntercept: got = %v, want = %v\n", got, want)
	}
}

func shuffle(arr []ColourPolymorphicPrey) {
	rand.Seed(12345) // no shuffling without this line

	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func TestSearchAndAttack(t *testing.T) {
	rand.Seed(0)
	predator := vpTesterAgent(0, 0)
	// fmt.Println(predator.String())
	prey := cppTestPop(1000)
	want := &prey[539]

	// for i := range prey {
	// 	fmt.Printf("%v\t%p\n", i, &prey[i])
	// }

	var err error
	target, δ, err := predator.PreySearch(prey, TestContext.VpSearchChance) //	previous test made sure this was true.
	if target == nil {
		t.Errorf("No target found in PreySearch")
		return
	}
	if target != want {
		t.Errorf("Targeting failed, want: %p\tgot: %p\n", want, target)
	}
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	intercepted, err := predator.Intercept(target.pos, δ)
	if !intercepted {
		t.Errorf("Failed to intercept target.")
	}
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	success := predator.Attack(target, TestContext.VpAttackChance, TestContext.VpCaf, TestContext.Vbg, TestContext.Vb𝛄, TestContext.Vbε)
	if !success {
		t.Errorf("Attack unsuccessful.")
	}
}
