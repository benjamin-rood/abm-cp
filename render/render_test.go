package render

import "testing"

func Test(t *testing.T) {
	var absToViewTests = []struct {
		p, d, n, want float64
	}{
		{-0.5, 1.0, 640.0, 160.0},
		{-0.713, 0.8, 480.0, 26.1},
		{0.31888, 2.0, 1600, 927.552},
		{0.99913, 1.0, 1024, 1023.555},
	}

	for _, atv := range absToViewTests {
		got := absToView(atv.p, atv.d, atv.n)
		if got != atv.want {
			t.Errorf("absToView(%v, %v, %v) == %v, want %v\n", atv.p, atv.d, atv.n, got, atv.want)
		}
	}
}
