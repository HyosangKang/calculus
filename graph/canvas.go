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

type canvas struct {
	c        *image.Paletted
	label    []string
	w, h     int
	gs       []*graph
	xb, yb   [2]float64
	trX, trY func(float64) int
	pad      int
}

func NewCanvas(w, h int) *canvas {
	pal := []color.Color{
		color.White,
		color.Black,
	}
	return &canvas{
		c:   image.NewPaletted(image.Rect(0, 0, w, h), pal),
		w:   w,
		h:   h,
		xb:  [2]float64{0, 0},
		yb:  [2]float64{0, 0},
		pad: 10,
	}
}

func (c *canvas) Add(g *graph, co color.Color, lab string) {
	c.gs = append(c.gs, g)
	c.c.Palette = append(c.c.Palette, co)
	c.label = append(c.label, lab)
	for _, s := range g.xr {
		if s > c.xb[1] {
			c.xb[1] = s
		}
		if s < c.xb[0] {
			c.xb[0] = s
		}
	}
	for _, s := range g.yr {
		if s > c.yb[1] {
			c.yb[1] = s
		}
		if s < c.yb[0] {
			c.yb[0] = s
		}
	}
	c.trX = func(t float64) int {
		return c.pad + int((t-c.xb[0])/(c.xb[1]-c.xb[0])*float64(c.w-2*c.pad))
	}
	c.trY = func(t float64) int {
		return c.pad + int((c.yb[1]-t)/(c.yb[1]-c.yb[0])*float64(c.h-2*c.pad))
	}
}

func (c *canvas) clear() {
	for i := 0; i < c.w; i++ {
		for j := 0; j < c.h; j++ {
			c.c.Set(i, j, c.c.Palette[0])
		}
	}
}

func (c *canvas) addAxes() {
	ax := &graph{
		xr: lin(c.xb[0], c.xb[1], 100),
		yr: lin(0, 0, 100),
	}
	for i := 0; i < 100; i++ {
		c.drawline(ax, i, 1)
	}
	ax = &graph{
		xr: lin(0, 0, 100),
		yr: lin(c.yb[0], c.yb[1], 100),
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
			if math.Abs(g.xr[i]-g.xr[i+1]) >= c.xb[1]-c.xb[0] {
				continue
			}
			if math.Abs(g.yr[i]-g.yr[i+1]) >= c.yb[1]-c.yb[0] {
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
	xd := [2]int{c.trX(g.xr[n]), c.trX(g.xr[n+1])}
	yd := [2]int{c.trY(g.yr[n]), c.trY(g.yr[n+1])}
	m := findMax(xd, yd)
	for j := 0; j < m; j++ {
		xi := xd[0] + int(float64((xd[1]-xd[0])*j)/float64(m))
		yi := yd[0] + int(float64((yd[1]-yd[0])*j)/float64(m))
		c.c.Set(xi, yi, c.c.Palette[co])
	}
}

func findMax(a, b [2]int) int {
	aa := abs(a[0] - a[1])
	bb := abs(b[1] - b[0])
	if aa > bb {
		return aa
	}
	return bb
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
