package graph

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var WIDTH, HEIGHT int = 600, 600
var PAD int = 10
var EPS float64 = 1e-6

type Graph struct {
	Canvas *image.Paletted
	Func   func(float64) float64
	Xrange [2]float64
}

func New(f func(float64) float64, xr [2]float64) *Graph {
	var palette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff}, // black
		color.RGBA{0xff, 0xff, 0xff, 0xff}} // white
	img := image.NewPaletted(image.Rect(0, 0, WIDTH, HEIGHT), palette)
	for i := 0; i < WIDTH; i++ {
		for j := 0; j < HEIGHT; j++ {
			img.Set(i, j, img.Palette[1])
		}
	}
	return &Graph{
		Canvas: img,
		Func:   f,
		Xrange: xr,
	}
}

func (g *Graph) Draw(filename string) {
	xss, yss, yb := split(g.Func, g.Xrange)
	for i := 0; i < len(xss); i++ {
		drawCurve(g.Canvas, xss[i], yss[i], g.Xrange, yb)
	}
	fp, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	png.Encode(fp, g.Canvas)
	fp.Close()
}

func split(f func(float64) float64, xb [2]float64) ([][]float64, [][]float64, [2]float64) {
	xss := [][]float64{}
	yss := [][]float64{}
	xs := []float64{}
	ys := []float64{}
	for x := xb[0]; x <= xb[1]; x += 1e-3 {
		y := f(x)
		if y != y || y > 1e6 || y < -1e6 {
			if len(xs) > 0 {
				xss = append(xss, xs)
				yss = append(yss, ys)
			}
			xs = []float64{}
			ys = []float64{}
			continue
		}
		xs = append(xs, x)
		ys = append(ys, y)
	}
	xss = append(xss, xs)
	yss = append(yss, ys)
	return xss, yss, bounds(yss)
}

func bounds(yss [][]float64) [2]float64 {
	var c, d float64 = 1e6, -1e6
	for _, ys := range yss {
		for _, y := range ys {
			if y < c {
				c = y
			}
			if y > d {
				d = y
			}
		}
	}
	return [2]float64{c, d}
}

func drawCurve(c *image.Paletted, xs, ys []float64, xb, yb [2]float64) {
	for i := 0; i < len(xs)-1; i++ {
		ps := []int{trX(xs[i], xb), trY(ys[i], yb)}
		pe := []int{trX(xs[i+1], xb), trY(ys[i+1], yb)}
		drawLine(c, ps, pe, xb, yb)
	}
}

func drawLine(c *image.Paletted, ps, pe []int, xb, yb [2]float64) {
	for i := 0; i < 1000; i++ {
		x := ps[0] + int(float64(pe[0]-ps[0])/float64(1000)*float64(i))
		y := ps[1] + int(float64(pe[1]-ps[1])/float64(1000)*float64(i))
		c.Set(x, y, c.Palette[0])
	}
}

func trX(x float64, xb [2]float64) int {
	return PAD + int(float64(WIDTH-2*PAD)*(x-xb[0])/(xb[1]-xb[0]))
}

func trY(y float64, yb [2]float64) int {
	return HEIGHT - PAD - int(float64(HEIGHT-2*PAD)*(y-yb[0])/(yb[1]-yb[0]))
}
