package plot

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func Plot(t []float64, ys [][]float64) {
	width, height := 600, 600
	pad := 10
	tb := [2]float64{t[0], t[len(t)-1]}
	yb := findYb(ys)
	pal := []color.Color{
		color.White,
		color.Black,
	}
	c := image.NewPaletted(image.Rect(0, 0, width, height), pal)
	clear(c)
	axes(c, tb, yb, pad)
	a := make([]int, 2)
	b := make([]int, 2)
	for _, y := range ys {
		for i := 0; i < len(t)-1; i++ {
			for j := 0; j < 2; j++ {
				a[j] = pad + int((t[i+j]-tb[0])/(tb[1]-tb[0])*float64(width-2*pad))
				b[j] = pad + int((yb[1]-y[i+j])/(yb[1]-yb[0])*float64(height-2*pad))
			}
			n := findMax(a, b)
			var ai, bi int
			for j := 0; j < n; j++ {
				ain := a[0] + int(float64((a[1]-a[0])*j)/float64(n))
				bin := b[0] + int(float64((b[1]-b[0])*j)/float64(n))
				if ai == ain && bi == bin {
					continue
				}
				ai, bi = ain, bin
				c.Set(ai, bi, c.Palette[1])
			}
		}
	}
	fp, _ := os.Create("graph.png")
	defer fp.Close()
	png.Encode(fp, c)
}

func findYb(ys [][]float64) [2]float64 {
	ymin, ymax := math.MaxFloat64, -math.MaxFloat64
	for _, y := range ys {
		for _, yy := range y {
			if yy > ymax {
				ymax = yy
			}
			if yy < ymin {
				ymin = yy
			}
		}
	}
	if ymin > 0 {
		ymin = 0
	}
	if ymax < 0 {
		ymax = 0
	}
	return [2]float64{ymin, ymax}
}

func clear(c *image.Paletted) {
	w, h := c.Bounds().Max.X, c.Bounds().Max.Y
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			c.Set(i, j, c.Palette[0])
		}
	}
}

func axes(c *image.Paletted, tb, yb [2]float64, p int) {
	w, h := c.Bounds().Max.X, c.Bounds().Max.Y
	t0 := p + int((-tb[0])/(tb[1]-tb[0])*float64(w-2*p))
	y0 := p + int(yb[1]/(yb[1]-yb[0])*float64(h-2*p))
	for i := p; i < w-p; i++ {
		c.Set(i, y0, c.Palette[1])
	}
	for i := p; i < h-p; i++ {
		c.Set(t0, i, c.Palette[1])
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
