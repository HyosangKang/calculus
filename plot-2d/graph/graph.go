package graph

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

var WIDTH, HEIGHT int = 600, 600
var PAD int = 10
var EPS float64 = 1e-6
var MAX float64 = 1e6

type Graph struct {
	Xr, Yr []float64
	Xb, Yb [2]float64
}

func New(f func(float64) float64, xb [2]float64, n int) *Graph {
	xr := make([]float64, n+1)
	yr := make([]float64, n+1)
	ymin, ymax := math.MaxFloat64, -math.MaxFloat64
	for i := 0; i < n+1; i++ {
		x := xb[0] + (xb[1]-xb[0])*float64(i)/float64(n)
		y := f(x)
		if y > MAX || y < -MAX {
			xr[i], yr[i] = math.NaN(), math.NaN()
			continue
		}
		xr[i], yr[i] = x, y
		if y > ymax {
			ymax = y
		}
		if y < ymin {
			ymin = y
		}
	}
	if ymin > 0 {
		ymin = 0
	}
	if ymax < 0 {
		ymax = 0
	}

	return &Graph{
		Xr: xr,
		Yr: yr,
		Xb: [2]float64{xr[0], xr[n]},
		Yb: [2]float64{ymin, ymax},
	}
}

func (g *Graph) Draw(filename string) {
	var pal = []color.Color{
		color.White,
		color.Black,
	}
	img := image.NewPaletted(image.Rect(0, 0, WIDTH, HEIGHT), pal)
	for i := 0; i < WIDTH; i++ {
		for j := 0; j < HEIGHT; j++ {
			img.Set(i, j, img.Palette[0])
		}
	}

	x0 := PAD + int((-g.Xb[0])/(g.Xb[1]-g.Xb[0])*float64(WIDTH-2*PAD))
	y0 := PAD + int(g.Yb[1]/(g.Yb[1]-g.Yb[0])*float64(HEIGHT-2*PAD))
	for i := PAD; i < WIDTH-PAD; i++ {
		img.Set(i, y0, img.Palette[1])
	}
	for i := PAD; i < HEIGHT-PAD; i++ {
		img.Set(x0, i, img.Palette[1])
	}

	for i := 0; i < len(g.Xr)-1; i++ {
		if math.IsNaN(g.Xr[i]) || math.IsNaN(g.Xr[i+1]) {
			continue
		}
		if math.Abs(g.Yr[i]-g.Yr[i+1]) >= g.Yb[1]-g.Yb[0] {
			continue
		}
		g.drawline(img, i)
	}
	fp, _ := os.Create("graph.png")
	defer fp.Close()
	png.Encode(fp, img)
}

func (g *Graph) drawline(c *image.Paletted, i int) {
	a := make([]int, 2)
	b := make([]int, 2)
	for j := 0; j < 2; j++ {
		a[j] = PAD + int((g.Xr[i+j]-g.Xb[0])/(g.Xb[1]-g.Xb[0])*float64(WIDTH-2*PAD))
		b[j] = PAD + int((g.Yb[1]-g.Yr[i+j])/(g.Yb[1]-g.Yb[0])*float64(HEIGHT-2*PAD))
	}
	m := findMax(a, b)
	var ai, bi int
	for j := 0; j < m; j++ {
		ain := a[0] + int(float64((a[1]-a[0])*j)/float64(m))
		bin := b[0] + int(float64((b[1]-b[0])*j)/float64(m))
		if ai == ain && bi == bin {
			continue
		}
		ai, bi = ain, bin
		c.Set(ai, bi, c.Palette[1])
	}
}

func findMax(a, b []int) int {
	ab, bb := abs(a[1]-a[0]), abs(b[1]-b[0])
	if ab > bb {
		return ab
	}
	return bb
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// func trX(x float64, xb [2]float64) int {
// 	return PAD + int(float64(WIDTH-2*PAD)*(x-xb[0])/(xb[1]-xb[0]))
// }

// func trY(y float64, yb [2]float64) int {
// 	return HEIGHT - PAD - int(float64(HEIGHT-2*PAD)*(y-yb[0])/(yb[1]-yb[0]))
// }
