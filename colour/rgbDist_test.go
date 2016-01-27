package colour

import "testing"

const twoThirds = 2 / 3

func TestRGBDist(t *testing.T) {
	var rgbDistTests = []struct {
		ca   RGB
		cb   RGB
		want float64
	}{
		{RGB{1.0, 1.0, 1.0}, RGB{1.0, 1.0, 1.0}, 0.0},
		{RGB{1.0, 0.0, 0.0}, RGB{0.0, 0.0, 1.0}, 0.667},
		{RGB{0.3, 0.6, 0.9}, RGB{1.0, 0.7, 0.4}, 0.433},
	}

	for _, rdt := range rgbDistTests {
		got := RGBDistance(rdt.ca, rdt.cb)
		if rdt.want != got {
			t.Errorf("RGBDistance(%v, %v): got = %v, want = %v\n", rdt.ca, rdt.cb, got, rdt.want)
		}
	}
}
