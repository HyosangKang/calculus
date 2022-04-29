package main

import (
	"main/plot"
	"main/rk"
)

func main() {
	f1 := func(t float64, ys []float64) float64 {
		return -ys[1]
	}
	f2 := func(t float64, ys []float64) float64 {
		return -ys[0]
	}
	fs := []func(float64, []float64) float64{f1, f2}
	tl, yl := rk.RungeKutta(fs, [2]float64{0, 1}, []float64{10, 5}, 0.01)
	plot.Plot(tl, yl)
}
