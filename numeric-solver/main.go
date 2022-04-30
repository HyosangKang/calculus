package main

import (
	// "main/plot"
	// "main/solver"
	// "math"
	"main/game"
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
	// fs := []func(float64, []float64) float64{
	// 	func(t float64, ys []float64) float64 {
	// 		return ys[0] - t*t + 1
	// 	},
	// }
	// tb := [2]float64{0, 2}
	// init := []float64{0}
	// nn := 2
	// tl, ylr := solver.rungekutta(fs, tb, init, nn)
	// _, yle := solver.euler(fs, tb, init, nn)
	// n := len(yle[0])
	// yt := make([]float64, n)
	// for i := 0; i < n; i++ {
	// 	t := tb[0] + (tb[1]-tb[0])/float64(nn)*float64(i)
	// 	yt[i] = t*t + 2*t + 1 - math.exp(t)
	// }
	// label := [][]string{
	// 	{"runge-kutta"},
	// 	{"euler"},
	// 	{"exact"},
	// }
	// plot.plot(tl, label, ylr, yle, [][]float64{yt})

	// A game test
	// a := []int{400, 100, 2}
	// b := []int{200, 100, 2}
	// g := game.NewGame(false, a, b)
	// g.Report()
	// g.Start()
	// g.Report()
}
