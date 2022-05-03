package graph

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var WIDTH, HEIGHT int = 600, 600
var PAD int = 10
var EPS float64 = 1e-6
var MAX float64 = 1e6

type graph struct {
	xr, yr []float64
}

func NewGraph(fx, fy func(float64) float64, t0, t1 float64, n int) *graph {
	f := []func(float64) float64{fx, fy}
	r := make([][]float64, 2)
	for i := 0; i < n+1; i++ {
		t := t0 + (t1-t0)*float64(i)/float64(n)
		for j := 0; j < 2; j++ {
			s := f[j](t)
			if math.Abs(s) > MAX {
				r[j] = append(r[j], math.NaN())
				continue
			}
			r[j] = append(r[j], s)
		}
	}
	return &graph{
		xr: r[0],
		yr: r[1],
	}
}

type canvas struct {
	c     *image.Paletted
	label []string
	sz    [2]int
	gs    []*graph
	b     [2][2]float64
	tr    [2]func(float64) int
	pad   int
}

func NewCanvas(w, h int) *canvas {
	pal := []color.Color{
		color.White,
		color.Black,
	}
	return &canvas{
		c:   image.NewPaletted(image.Rect(0, 0, w, h), pal),
		sz:  [2]int{w, h},
		b:   [2][2]float64{{0, 0}, {0, 0}},
		pad: 10,
	}
}

func (c *canvas) Add(g *graph, co color.Color, lab string) {
	c.gs = append(c.gs, g)
	c.c.Palette = append(c.c.Palette, co)
	c.label = append(c.label, lab)
	ss := [2][]float64{g.xr, g.yr}
	for j := 0; j < 2; j++ {
		for _, s := range ss[j] {
			if s > c.b[j][1] {
				c.b[j][1] = s
			}
			if s < c.b[j][0] {
				c.b[j][0] = s
			}
		}
		j := j
		c.tr[j] = func(t float64) int {
			return c.pad + int((t-c.b[j][0])/(c.b[j][1]-c.b[j][0])*float64(c.sz[j]-2*c.pad))
		}
	}
}

func (c *canvas) clear() {
	w, h := c.c.Bounds().Max.X, c.c.Bounds().Max.Y
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			c.c.Set(i, j, c.c.Palette[0])
		}
	}
}

func (c *canvas) addAxes() {
	ax := &graph{
		xr: lin(c.b[0][0], c.b[0][1], 100),
		yr: lin(0, 0, 100),
	}
	for i := 0; i < 100; i++ {
		c.drawline(ax, i, 1)
	}
	ax = &graph{
		xr: lin(0, 0, 100),
		yr: lin(c.b[1][0], c.b[1][1], 100),
	}
	for i := 0; i < 100; i++ {
		c.drawline(ax, i, 1)
	}
}

func lin(t0, t1 float64, n int) []float64 {
	v := []float64{}
	for i := 0; i < n+1; i++ {
		v = append(v, t0+(t1-t0)*float64(i)/float64(n))
	}
	return v
}

func (c *canvas) Draw(fn string) {
	c.clear()
	c.addAxes()
	for j, g := range c.gs {
		for i := 0; i < len(g.xr)-1; i++ {
			if math.IsNaN(g.xr[i]) || math.IsNaN(g.xr[i+1]) {
				continue
			}
			if math.Abs(g.yr[i]-g.yr[i+1]) >= c.b[1][1]-c.b[1][0] {
				continue
			}
			c.drawline(g, i, j+2)
		}
	}
	c.addLabel()
	fp, _ := os.Create("graph.png")
	defer fp.Close()
	png.Encode(fp, c.c)
}

func (c *canvas) addLabel() {
	d := &font.Drawer{
		Dst:  c.c,
		Src:  image.NewUniform(color.RGBA{0, 0, 0, 255}),
		Face: basicfont.Face7x13,
	}
	for i, lab := range c.label {
		xi := c.pad + 10
		yi := c.pad + 15*(i+1)
		for j := 0; j < 30; j++ {
			c.c.Set(xi+j, yi-5, c.c.Palette[i+2])
		}
		d.Dot = fixed.Point26_6{fixed.I(xi + 40), fixed.I(yi)}
		d.DrawString(lab)
	}
}

func (c *canvas) drawline(g *graph, n, co int) {
	ss := [][]float64{g.xr[n : n+2], g.yr[n : n+2]}
	var a [2][2]int
	for j := 0; j < 2; j++ {
		for i := 0; i < 2; i++ {
			a[j][i] = c.tr[j](ss[j][i])
		}
	}
	m := findMax(a)
	for j := 0; j < m; j++ {
		xi := a[0][0] + int(float64((a[0][1]-a[0][0])*j)/float64(m))
		yi := a[1][0] + int(float64((a[1][1]-a[1][0])*j)/float64(m))
		c.c.Set(xi, yi, c.c.Palette[co])
	}
}

func findMax(a [2][2]int) int {
	ab, bb := abs(a[0][1]-a[0][0]), abs(a[1][1]-a[1][0])
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
