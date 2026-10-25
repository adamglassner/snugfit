package main

import "testing"

func TestGrid(t *testing.T) {
	for _, a := range []float64{2752.0 / 1536, 16.0 / 9, 1.5, 4.0 / 3} {
		for i := range modes {
			c, r := gridFor(modes[i].Target, a)
			t.Logf("aspect %.3f %-9s %dx%d = %d (piece aspect %.2f)", a, modes[i].Name, c, r, c*r, a*float64(r)/float64(c))
		}
	}
}
