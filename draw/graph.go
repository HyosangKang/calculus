package draw

import (
	"math"
)

var MAX float64 = 1e6

type graph struct {
	xr, yr []float64
}

func NewGraphArr(xr, yr []float64) graph {
	if len(xr) != len(yr) {
		panic("Invalid graph array.")
	}
	return graph{
		xr: xr,
		yr: yr,
	}
}

func NewGraph(fx, fy func(float64) float64, t0, t1 float64, n int) graph {
	f := []func(float64) float64{fx, fy}
	r := make([][]float64, 2)
	for i := 0; i < n+1; i++ {
		t := t0 + (t1-t0)*float64(i)/float64(n)
		for j := 0; j < 2; j++ {
			s := f[j](t)
			if math.Abs(s) > MAX {
				s = math.NaN()
			}
			r[j] = append(r[j], s)
		}
	}
	return graph{
		xr: r[0],
		yr: r[1],
	}
}
