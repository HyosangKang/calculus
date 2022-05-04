package main

import (
	"calculus/draw"
	"calculus/game"
	"calculus/solver"
	"image/color"
)

func main() {
	// minion setting: num, hp, dps
	a := []int{400, 100, 2}
	b := []int{200, 100, 15}
	c := draw.NewCanvas(600, 600)
	im := draw.NewImg(0, 0, 600, 600)

	// add numeric solution to a differential equation
	fs := []func(float64, []float64) float64{
		func(t float64, ys []float64) float64 {
			return -float64(b[2]) * ys[1] / float64(a[1])
		},
		func(t float64, ys []float64) float64 {
			return -float64(a[2]) * ys[0] / float64(b[1])
		},
	}
	var dur float64 = 40
	tb := [2]float64{0, dur}
	nn := 100
	init := []float64{float64(a[0]), float64(b[0])}
	xs, yss := solver.RungeKutta(fs, tb, init, nn)
	ei := 0
	for i := 0; i < len(xs); i++ {
		for _, ys := range yss {
			if ys[i] <= 0 {
				dur = xs[i]
				ei = i
				goto RESUME
			}
		}
	}
RESUME:
	g := draw.NewGraphArr(xs[:ei], yss[0][:ei])
	red := color.RGBA{255, 0, 0, 255}
	im.Add(g, red, "RK-A")
	g = draw.NewGraphArr(xs[:ei], yss[1][:ei])
	blue := color.RGBA{0, 0, 255, 255}
	im.Add(g, blue, "RK-B")

	// game simulation
	ga := game.NewGame(false, a, b)
	ga.Start()
	xss, yss := ga.Report(dur)
	g = draw.NewGraphArr(xss[0], yss[0])
	im.Add(g, red, "Sim-A")
	g = draw.NewGraphArr(xss[1], yss[1])
	im.Add(g, blue, "Sim-B")

	// draw results
	c.Add(im)
	c.Draw()
	c.Save("graph.png")
}
