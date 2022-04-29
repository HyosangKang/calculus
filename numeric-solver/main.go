package main

import (
	"main/plot"
	"main/solver"
)

func main() {
	// fs := []func(float64, []float64) float64{
	// 	func(t float64, ys []float64) float64 {
	// 		return -ys[1]
	// 	},
	// 	func(t float64, ys []float64) float64 {
	// 		return -ys[0]
	// 	},
	// }

	// xb := [2]float64{0, 10}
	// init := []float64{10, 5}
	// h := 0.1
	fs := []func(float64, []float64) float64{
		func(t float64, ys []float64) float64 {
			return t
		},
	}
	tb := [2]float64{0, .3}
	init := []float64{0}
	h := 0.1
	tl, ylr := solver.RungeKutta(fs, tb, init, h)
	_, yle := solver.Euler(fs, tb, init, h)
	n := len(yle[0])
	yt := make([]float64, n)
	for i := 0; i < n; i++ {
		t := (tb[0] + (tb[1]-tb[0])/float64(n-1)*float64(i))
		yt[i] = t * t / 2
	}
	label := []string{
		"Runge-Kutta",
		"Euler",
		"Exact",
	}
	plot.Plot(tl, label, ylr, yle, [][]float64{yt})
}
