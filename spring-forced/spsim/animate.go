package spsim

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"
)

func (s *SpSim) Animate() {
	s.GI = initGI()
	for i := 0; i < s.N; i++ {
		s.addFrame(i + 1)
	}
	fp, err := os.Create("graph.gif")
	if err != nil {
		panic(err)
	}
	err = gif.EncodeAll(fp, s.G)
	if err != nil {
		panic(err)
	}
	fp.Close()
	fmt.Printf("\nThe spring simulation is saved on `graph.gif`.\n\n")
}

func (s *SpSim) addFrame(count int) {
	palette := []color.Color{
		color.White,
		color.Black,
		color.NRGBA{255, 0, 0, 255},
		color.NRGBA{0, 0, 255, 255}}
	c := image.NewPaletted(image.Rect(0, 0, s.GI.Width, s.GI.Height), palette)
	s.clear(c)
	s.axes(c)

	s.drawGraph(c, count)
	s.drawSpring(c, count)

	s.G.Image = append(s.G.Image, c)
	s.G.Delay = append(s.G.Delay, s.GI.Delay)
}

func (s *SpSim) clear(c *image.Paletted) {
	for i := 0; i < s.GI.Width; i++ {
		for j := 0; j < s.GI.Height; j++ {
			c.Set(i, j, c.Palette[0])
		}
	}
}

func (s *SpSim) axes(c *image.Paletted) {
	xz, yz := s.trX(0), s.trY(0)
	for i := s.GI.Pad + s.GI.Ceil + s.GI.SpL; i <= s.GI.Height-s.GI.Pad; i++ {
		c.Set(xz, i, c.Palette[1])
	}
	for i := s.GI.SpWidth + s.GI.Pad; i <= s.GI.Width-s.GI.Pad; i++ {
		c.Set(i, yz, c.Palette[1])
	}
}

func (s *SpSim) trX(x float64) int {
	xi := s.GI.Pad + s.GI.SpWidth
	xi += int(float64(s.GI.GphWidth-2*s.GI.Pad) * (x - s.Xb[0]) / (s.Xb[1] - s.Xb[0]))
	return xi
}

func (s *SpSim) trY(y float64) int {
	yi := s.GI.Pad + s.GI.Ceil + s.GI.SpL
	yi += int(float64(s.GI.GphWidth-2*s.GI.Pad) * (s.Yb[1] - y) / (s.Yb[1] - s.Yb[0]))
	return yi
}

func (s *SpSim) drawGraph(c *image.Paletted, count int) {
	for i := 0; i < count; i++ {
		x0, x1 := s.Xb[1]/float64(s.N)*float64(i), s.Xb[1]/float64(s.N)*float64(i+1)
		y0, y1 := s.F(x0), s.F(x1)
		p0 := [2]int{s.trX(x0), s.trY(y0)}
		p1 := [2]int{s.trX(x1), s.trY(y1)}
		drawLine(c, p0, p1, 1)
		y0, y1 = s.Fc(x0), s.Fc(x1)
		p0 = [2]int{s.trX(x0), s.trY(y0)}
		p1 = [2]int{s.trX(x1), s.trY(y1)}
		drawLine(c, p0, p1, 2)
		y0, y1 = s.Fp(x0), s.Fp(x1)
		p0 = [2]int{s.trX(x0), s.trY(y0)}
		p1 = [2]int{s.trX(x1), s.trY(y1)}
		drawLine(c, p0, p1, 3)
	}
}

func drawLine(c *image.Paletted, p0, p1 [2]int, color int) {
	x, y := 0, 0
	n := findMax(p0, p1)
	for i := 0; i <= n; i++ {
		nx := p0[0] + int(float64((p1[0]-p0[0])*i)/float64(n))
		ny := p0[1] + int(float64((p1[1]-p0[1])*i)/float64(n))
		for nx == x && ny == y {
			continue
		}
		x, y := nx, ny
		c.Set(x, y, c.Palette[color])
	}
}

func findMax(p0, p1 [2]int) int {
	mx := absInt(p0[0] - p1[0])
	my := absInt(p0[1] - p1[1])
	if mx >= my {
		return mx
	}
	return my
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (s *SpSim) drawSpring(c *image.Paletted, count int) {
	s.addCeil(c)
	s.addSpring(c, count)
	s.addCircle(c, count)
}

func (s *SpSim) addCeil(c *image.Paletted) {
	for i := s.GI.Pad; i < s.GI.SpWidth-s.GI.Pad; i++ {
		c.Set(i, s.GI.Ceil, c.Palette[1])
	}
}

func (s *SpSim) addSpring(c *image.Paletted, count int) {
	for i := 0; i < 100; i++ {
		p0 := s.springXY(i, count)
		p1 := s.springXY(i+1, count)
		drawLine(c, p0, p1, 1)
	}
}

func (s *SpSim) springXY(i, count int) [2]int {
	n := s.GI.Coil[0]
	r := s.GI.Coil[1]
	d := s.GI.Coil[2]
	y := s.F(s.Xb[1] / float64(s.N) * float64(count))
	l := s.trY(y)
	l -= s.GI.Ceil + s.GI.CR + 2*d + 2*r
	st := -math.Pi / 2
	xi := int(float64(r) * math.Cos(st+(2*math.Pi*float64(n)-2*st)*float64(i)/100))
	xi += s.GI.SpWidth/2 + s.GI.Pad
	yi := int(float64(r) * math.Sin(st+(2*math.Pi*float64(n)-2*st)*float64(i)/100))
	yi += int(float64(l) / float64(100) * float64(i))
	yi += d + r + s.GI.Ceil
	return [2]int{xi, yi}
}

func (s *SpSim) addCircle(c *image.Paletted, count int) {
	y := s.F(s.Xb[1] / float64(s.N) * float64(count))
	for i := 0; i < 100; i++ {
		p0 := s.circleXY(s.trY(y), i)
		p1 := s.circleXY(s.trY(y), i+1)
		drawLine(c, p0, p1, 1)
	}
}

func (s *SpSim) circleXY(l, i int) [2]int {
	r := s.GI.CR
	xi := int(float64(r) * math.Cos(2*math.Pi*float64(i)/100))
	xi += s.GI.SpWidth/2 + s.GI.Pad
	yi := int(float64(r) * math.Sin(2*math.Pi*float64(i)/100))
	yi += l
	return [2]int{xi, yi}
}
