package draw

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
	c   *image.Paletted
	img []*img
}

func NewCanvas(w, h int) *canvas {
	pal := []color.Color{
		color.White,
		color.Black,
	}
	return &canvas{
		c: image.NewPaletted(image.Rect(0, 0, w, h), pal),
	}
}

func (c *canvas) Add(img *img) {
	c.img = append(c.img, img)
	for _, co := range img.color {
		c.c.Palette = append(c.c.Palette, co)
	}
}

func (c *canvas) Draw() {
	for _, im := range c.img {
		im.addAxes()
		for j, g := range im.graph {
			for i := 0; i < len(g.xr)-1; i++ {
				if math.IsNaN(g.xr[i]) || math.IsNaN(g.xr[i+1]) {
					continue
				}
				if math.Abs(g.xr[i]-g.xr[i+1]) >= im.xb[1]-im.xb[0] {
					continue
				}
				if math.Abs(g.yr[i]-g.yr[i+1]) >= im.yb[1]-im.yb[0] {
					continue
				}
				c.drawline(im, j, i)
			}
			c.addLabel(im)
		}
	}
}

func (c *canvas) Save(fn string) {
	fp, _ := os.Create(fn)
	defer fp.Close()
	png.Encode(fp, c.c)
}

func (c *canvas) drawline(im *img, gn, n int) {
	g := im.graph[gn]
	xd := [2]int{im.trX(g.xr[n]), im.trX(g.xr[n+1])}
	yd := [2]int{im.trY(g.yr[n]), im.trY(g.yr[n+1])}
	m := findMax(xd, yd)
	for j := 0; j < m; j++ {
		xi := xd[0] + int(float64((xd[1]-xd[0])*j)/float64(m))
		yi := yd[0] + int(float64((yd[1]-yd[0])*j)/float64(m))
		c.c.Set(xi, yi, im.color[gn])
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

func (c *canvas) addLabel(im *img) {
	d := &font.Drawer{
		Dst:  c.c,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
	}
	xi, yi := im.x0+im.pad+15, im.y0+im.pad+15
	for i, lab := range im.label {
		if len(lab) == 0 {
			continue
		}
		for j := 0; j < 30; j++ {
			c.c.Set(xi+j, yi-5, im.color[i])
		}
		d.Dot = fixed.Point26_6{fixed.I(xi + 40), fixed.I(yi)}
		d.DrawString(lab)
		yi += 15
	}
}
