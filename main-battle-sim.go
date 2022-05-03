package main

import (
	"calculus/game"
	"calculus/graph"
	"calculus/solver"
	"fmt"
	"image/color"
)

func main() {
	c := graph.NewCanvas(600, 600)

	// add numeric solution to a differential equation
	fs := []func(float64, []float64) float64{
		func(t float64, ys []float64) float64 {
			return -ys[1] / 8
		},
		func(t float64, ys []float64) float64 {
			return -.02 * ys[0]
		},
	}

	var dur float64 = 22
	tb := [2]float64{0, dur}
	init := []float64{400, 200}
	nn := 100
	xs, yss := solver.RungeKutta(fs, tb, init, nn)
	g := graph.NewGraphArr(xs, yss[0])
	c.Add(g, color.NRGBA{255, 0, 0, 255}, "RK-A")
	g = graph.NewGraphArr(xs, yss[1])
	c.Add(g, color.NRGBA{150, 50, 0, 255}, "RK-B")

	a := []int{400, 100, 2}
	b := []int{200, 80, 10}
	ga := game.NewGame(false, a, b)
	ga.Start()
	xss, yss := ga.Report(dur)
	fmt.Println(xss, yss)
	g = graph.NewGraphArr(xss[0], yss[0])
	c.Add(g, color.NRGBA{0, 0, 255, 255}, "Sim-A")
	g = graph.NewGraphArr(xs, yss[1])
	c.Add(g, color.NRGBA{0, 50, 150, 255}, "Sim-B")
	c.Draw("graph.png")
}
