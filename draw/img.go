package draw

import (
	"image/color"
)

type img struct {
	graph        []graph
	color        []color.Color
	label        []string
	xb, yb       [2]float64
	x0, y0, w, h int
	pad          int
	trX, trY     func(float64) int
}

func (im *img) TrX() func(float64) int {
	return im.trX
}

func (im *img) TrY() func(float64) int {
	return im.trY
}

func (img *img) Xb() [2]float64 {
	return img.xb
}

func (img *img) Yb() [2]float64 {
	return img.yb
}

func (im *img) SetTrX(trX func(float64) int) {
	im.trX = trX
}

func (im *img) SetTrY(trY func(float64) int) {
	im.trY = trY
}

func (im *img) SetXb(xb [2]float64) {
	im.xb = xb
}

func (im *img) SetYb(yb [2]float64) {
	im.yb = yb
}

func NewImg(x0, y0, w, h int) *img {
	im := img{
		xb:  [2]float64{0, 0},
		yb:  [2]float64{0, 0},
		pad: 10,
		x0:  x0,
		y0:  y0,
		w:   w,
		h:   h,
	}
	trX := func(t float64) int {
		return im.x0 + im.pad + int((t-im.xb[0])/(im.xb[1]-im.xb[0])*float64(im.w-2*im.pad))
	}
	trY := func(t float64) int {
		return im.y0 + im.pad + int((im.yb[1]-t)/(im.yb[1]-im.yb[0])*float64(im.h-2*im.pad))
	}
	im.trX = trX
	im.trY = trY
	return &im
}

func (c *img) Add(g graph, co color.Color, lab string) {
	c.graph = append(c.graph, g)
	c.color = append(c.color, co)
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
}

func (c *img) addAxes() {
	ax := graph{
		xr: lin(c.xb[0], c.xb[1], 100),
		yr: lin(0, 0, 100),
	}
	c.Add(ax, color.Black, "")
	ax = graph{
		xr: lin(0, 0, 100),
		yr: lin(c.yb[0], c.yb[1], 100),
	}
	c.Add(ax, color.Black, "")
}

func lin(t0, t1 float64, n int) []float64 {
	v := []float64{}
	for i := 0; i < n+1; i++ {
		v = append(v, t0+(t1-t0)*float64(i)/float64(n))
	}
	return v
}
