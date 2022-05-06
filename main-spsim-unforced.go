package main

import (
	"calculus/draw"
	"calculus/spsim"
	"image/color"
	"image/gif"
	"math"
	"os"
)

// This program animates the solution of the spring motion.
func main() {
	img0 := draw.NewImg(200, 200, 400, 400)
	xb, f := spsim.UnforcedSim()
	x := func(t float64) float64 {
		return t
	}
	y := func(t float64) float64 {
		return -f(t)
	}
	g0 := draw.NewGraphFunc(x, y, xb[0], xb[1], 100)
	xr := g0.Xr()
	yr := g0.Yr()
	img0.Add(g0, color.Opaque, "")
	trY := img0.TrY()
	cir := draw.NewGraphFunc(math.Cos, math.Sin, 0, 2*math.Pi, 100)
	spx := func(t float64) float64 {
		return 10 * math.Cos(10*(t+math.Pi/4))
	}
	ceil := draw.NewGraphArr([]float64{0, 10}, []float64{0, 0})

	gi := &gif.GIF{
		LoopCount: -1,
	}
	for i := 2; i < 101; i++ {
		img := img0.Copy()
		g := draw.NewGraphArr(xr[:i], yr[:i])
		img.Add(g, color.Black, "")

		// add circle
		ciry0 := trY(yr[i-1])
		img1 := draw.NewImg(100-20, ciry0-20, 40, 40)
		img1.Add(cir, color.Black, "")
		img1.SetAxis(false)
		// add spring
		spy := func(t float64) float64 {
			return 10*math.Sin(10*(t+math.Pi/4)) - float64(ciry0)*t
		}
		sp := draw.NewGraphFunc(spx, spy, 0, 2*math.Pi, 200)
		img2 := draw.NewImg(100-20, 0, 40, ciry0)
		img2.Add(sp, color.Black, "")
		img2.SetAxis(false)

		can := draw.NewCanvas(600, 600)
		can.Add(img)
		can.Add(img1)
		can.Add(img2)
		c := can.Draw()
		gi.Image = append(gi.Image, c)
		gi.Delay = append(gi.Delay, 1)
	}
	fp, _ := os.Create("unforced-1.gif")
	defer fp.Close()
	gif.EncodeAll(fp, gi)
}
