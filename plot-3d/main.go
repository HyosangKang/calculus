package main

import (
	"math"
	"multivariable/mvf"
)

func main() {
	f := mvf.NewMvf(func(x []float64) []float64 {
		return []float64{(x[0] + x[1]) / math.Sqrt(x[0]*x[0]+2*x[1]*x[1])}
	}, [][2]float64{{-1, 1}, {-1, 1}}, 2, 1)
	f.GraphSpherical("test.png", 20, math.Pi/4, math.Pi/6)

	// c := mvf.NewMvf(func(t []float64) []float64 {
	// 	return []float64{math.Cos(t[0]), math.Sin(t[0]), t[0]}
	// })
	// c.Graph1to3([]float64{0, 2 * math.Pi})
}
