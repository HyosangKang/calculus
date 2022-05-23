package main

import (
	"calculus/draw"
	"calculus/solver"
	"image/color"
	"math"
)

func main() {
	// System of ODE
	// fs := []func(float64, []float64) float64{
	// 	func(t float64, ys []float64) float64 {
	// 		return -ys[1]
	// 	},
	// 	func(t float64, ys []float64) float64 {
	// 		return -ys[0]
	// 	},
	// }

	// tb := [2]float64{0, .6}
	// init := []float64{10, 5}
	// nn := 100
	// tl, yl := solver.RungeKutta(fs, tb, init, nn)
	// plot.Plot(tl, [][]string{{"A", "B"}}, yl)

	// Comparision between RK, Euler, and exact
	fs := []func(float64, []float64) float64{
		func(t float64, ys []float64) float64 {
			return ys[0] - t*t + 1
		},
	}
	tb := [2]float64{0, 2}
	init := []float64{0}
	nn := 2
	tl, ylr := solver.RungeKutta(fs, tb, init, nn)
	_, yle := solver.Euler(fs, tb, init, nn)
	n := len(yle[0])
	yt := make([]float64, n)
	for i := 0; i < n; i++ {
		t := tb[0] + (tb[1]-tb[0])/float64(nn)*float64(i)
		yt[i] = t*t + 2*t + 1 - math.Exp(t)
	}
	c := draw.NewCanvas(600, 600)
	m := draw.NewImg(0, 0, 600, 600)
	g1 := draw.NewGraphArr(tl, yt)
	g2 := draw.NewGraphArr(tl, yle[0])
	g3 := draw.NewGraphArr(tl, ylr[0])
	bk := color.Black
	rd := color.RGBA{255, 0, 0, 255}
	bl := color.RGBA{0, 0, 255, 255}
	m.Add(g1, bk, "Exact")
	m.Add(g2, rd, "Euler")
	m.Add(g3, bl, "Runge-Kutta")
	c.Add(m)
	c.Draw()
	c.Save("graph-numeric.png")
}
