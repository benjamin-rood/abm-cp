package colour

import (
	"testing"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
)

func ColourFunctionTesting(t *testing.T) {
	var rgbDistTests = []struct {
		ca   RGB
		cb   RGB
		want float64
	}{
		{RGB{1.0, 1.0, 1.0}, RGB{1.0, 1.0, 1.0}, 0.0},
		{RGB{1.0, 0.0, 0.0}, RGB{0.0, 0.0, 1.0}, calc.ToFixed((2 / 3), 3)},
		{RGB{0.3, 0.6, 0.9}, RGB{1.0, 0.7, 0.4}, calc.ToFixed((1.9 / 3), 3)},
	}

	for _, rdt := range rgbDistTests {
		got := RGBDistance(rdt.ca, rdt.cb)
		if rdt.want != got {
			t.Errorf("RGBDistance(%v, %v): got = %v, want = %v\n", rdt.ca, rdt.cb, got, rdt.want)
		}
	}
}

/*
||1.0 - 0.0|| = 1.0
||0.0 - 0.0|| = 0.0
||0.0 - 1.0|| = 1.0



||0.3 - 1.0|| = 1.3
||0.6 - 0.7|| = 0.1
||0.9 - 0.4|| = 0.5

(1.3 + 0.1 + 0.5) = (1.9 / 3)
*/
