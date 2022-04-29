package main

import (
	"fmt"
	"main/rk"
	"math"
)

func main() {
	fs := []func(float64, []float64) float64{
		func(t float64, ys []float64) float64 {
			return math.Sin(t)
		},
	}
	tl, yl := rk.RungeKutta(fs, [2]float64{0, 1}, []float64{0}, 0.01)
	fmt.Println(tl)
}
