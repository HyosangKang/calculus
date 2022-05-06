package main

import (
	"calculus/draw"
	"image/color"
	"image/gif"
	"math"
	"os"
)

func main() {
	// define the function to draw
	x := func(t float64) float64 {
		return t
	}
	f1 := func(t float64) float64 {
		return math.Cos(t)
	}
	// add another graph
	f2 := func(t float64) float64 {
		return math.Sin(t)
	}
	g1 := draw.NewGraphFunc(x, f1, -2, 2, 100)
	xr := g1.Xr()
	y1r := g1.Yr()
	g2 := draw.NewGraphFunc(x, f2, -2, 2, 100)
	y2r := g2.Yr()
	im := draw.NewImg(0, 0, 600, 600)
	im.Add(g1, color.Opaque, "")
	im.Add(g2, color.Opaque, "")
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}

	g := &gif.GIF{
		LoopCount: -1,
	}
	for i := 1; i <= 100; i++ {
		g1p := draw.NewGraphArr(xr[:i], y1r[:i])
		g2p := draw.NewGraphArr(xr[:i], y2r[:i])
		ima := im.Copy()
		ima.Add(g1p, red, "cos")
		ima.Add(g2p, blue, "sin")
		can := draw.NewCanvas(600, 600)
		can.Add(ima)
		c := can.Draw()
		g.Image = append(g.Image, c)
		g.Delay = append(g.Delay, 1)
	}
	fp, _ := os.Create("graph.gif")
	defer fp.Close()
	gif.EncodeAll(fp, g)
}
