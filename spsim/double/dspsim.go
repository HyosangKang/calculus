package double

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type DspSim struct {
	Val []float64
	F   [2]func(float64) float64 // spring motion solution
	Xb  [2]float64
	Yb  [2][2]float64
	TrX func(float64) int
	TrY [2]func(float64) int
	N   int // number of intervals
	GI  GI
}

type GI struct {
	G         *gif.GIF
	Delay     int // Delay number (100th of seconds)
	Height    int // height (pixel) of image
	GphHeight int
	Width     int // width (pixel) of image
	SprWidth  int // spring canvas width (pixel) -- on the left
	GphWidth  int // graph canvas width (pixel) -- on the right
	Pad       int // padding size (pixel)
	Ceil      int // ceiling height (pixel)
	SpL       int // spring minimum length (pixel)
	Mr        int // circle radius (pixel)
	CoilN     int // coil number (int)
	CoilR     int // coil radius (pixel)
	CoilM     int // coil margin (pixel)
}

func NewGI() GI {
	return GI{
		G: &gif.GIF{
			LoopCount: -1,
		},
		Delay:     1,
		Height:    600,
		GphHeight: 175,
		Width:     600,
		SprWidth:  200,
		GphWidth:  400,
		Pad:       10,
		Ceil:      50,
		SpL:       100,
		Mr:        10,
		CoilN:     10,
		CoilR:     10,
		CoilM:     5,
	}
}

func NewDspSim() *DspSim {
	d := &DspSim{
		GI: NewGI(),
		N:  100,
	}
	d.handleInput()
	d.Solve()
	d.findBounds()
	d.findTr()
	return d
}

func (d *DspSim) handleInput() {
	msg := []string{
		"====================================\n",
		"  @       This program simulates\n",
		"  @ k1    the double spring motion\n",
		"  @       modeled by two displacement\n",
		"  m1      functions x1(t) and x2(t)\n",
		"  @       over the interval [0, T]\n",
		"  @ k2    Intial values are\n",
		"  @       x1(0)=x0, x1'(0)=x1\n",
		"  m2      x2(0)=x2, x2'(0)=x3\n",
		"=====================================\n",
	}
	for _, m := range msg {
		fmt.Printf("%s", m)
	}
	msg = []string{
		"m1: ",
		"m2: ",
		"k1: ",
		"k2: ",
		"x0: ",
		"x1: ",
		"x2: ",
		"x3: ",
		"T: ",
	}
	val := make([]float64, 9)
	for i, m := range msg {
		fmt.Printf("Enter the value of %s", m)
		fmt.Scanf("%f\n", &val[i])
	}
	d.Val = val
}

func (d *DspSim) Solve() {
	m1, m2, k1, k2 := d.Val[0], d.Val[1], d.Val[2], d.Val[3]
	x0, x1, x2, x3 := d.Val[4], d.Val[5], d.Val[6], d.Val[7]
	// Differential equation
	// m1*x1'' = -k1*x1 + k2*(x2-x1)
	// m2*x2'' = -k2*(x2-x1)
	// Laplace transform
	// (m1*s^2 + (k1+k2))*X1 - k2*X2 = m1*x0*s + m1*x1
	// -k2*X1 + (m2*s^2+k2)*X2 = m2*x2*s + m2*x3
	// a^2, b^2 are the (real) roots of t^2 + (k2/m2 + (k1+k2)/m1)*t + k2^2/(m1*m2)
	aa, bb := roots([2]float64{k2/m2 + (k1+k2)/m1, k2 * k2 / (m1 * m2)})
	a, b := math.Sqrt(-aa), math.Sqrt(-bb)
	// X1 = A*(s/(s^2+a^2) + B*(1/(s^2+a^2) + C*(s/(s^2+b^2) + D*(1/(s^2+b^2)
	// A + C = x0
	// b^2*A + a^2*C = k2*(x2/m1+x0/m2)
	// B + D = x1
	// b^2*B + a^2*D = k2(x3/m1+x1/m2)
	// X2 = A*(s/(s^2+a^2) + B*(1/(s^2+a^2) + C*(s/(s^2+b^2) + D*(1/(s^2+b^2)
	// A + C = x2
	// b^2*A + a^2*C = (k1+k2)/m1*x2 + k2/m2*x0
	// B + D = x3
	// b^2*B + a^2*D = (k1+k2)/m1*x3 + k2/m2*x1
	A0, C0 := cramer([2][2]float64{
		{1, 1},
		{-bb, -aa}}, [2]float64{x0, k2 * (x2/m1 + x0/m2)})
	B0, D0 := cramer([2][2]float64{
		{1, 1},
		{-bb, -aa}}, [2]float64{x1, k2 * (x3/m1 + x1/m2)})
	d.F[0] = func(t float64) float64 {
		return A0*math.Cos(a*t) + B0*math.Sin(a*t) + C0*math.Cos(b*t) + D0*math.Sin(b*t)
	}
	A1, C1 := cramer([2][2]float64{
		{1, 1},
		{-bb, -aa}}, [2]float64{x2, (k1+k2)/m1*x2 + k2/m2*x0})
	B1, D1 := cramer([2][2]float64{
		{1, 1},
		{-bb, -aa}}, [2]float64{x3, (k1+k2)/m1*x3 + k2/m2*x1})
	d.F[1] = func(t float64) float64 {
		return A1*math.Cos(a*t) + B1*math.Sin(a*t) + C1*math.Cos(b*t) + D1*math.Sin(b*t)
	}
}

func cramer(m [2][2]float64, b [2]float64) (float64, float64) {
	d := m[0][0]*m[1][1] - m[0][1]*m[1][0]
	x0 := (b[0]*m[1][1] - b[1]*m[0][1]) / d
	x1 := (m[0][0]*b[1] - b[0]*m[1][0]) / d
	return x0, x1
}

func roots(l [2]float64) (float64, float64) {
	if l[0]*l[0]-4*l[1] < 0 {
		panic("Cannot compute real roots.")
	}
	return (-l[0] - math.Sqrt(l[0]*l[0]-4*l[1])) / 2, (-l[0] + math.Sqrt(l[0]*l[0]-4*l[1])) / 2
}

func (d *DspSim) findBounds() {
	d.Xb = [2]float64{0, d.Val[8]}
	for k := 0; k < 2; k++ {
		ymin, ymax := math.MaxFloat64, -math.MaxFloat64
		for i := 0; i <= d.N; i++ {
			x := (d.Xb[1] - d.Xb[0]) / float64(d.N) * float64(i)
			y := -d.F[k](x)
			if y < ymin {
				ymin = y
			}
			if y > ymax {
				ymax = y
			}
		}
		if ymin > 0 {
			ymin = 0
		}
		if ymax < 0 {
			ymax = 0
		}
		d.Yb[k] = [2]float64{ymin, ymax}
	}
}

func (d *DspSim) findTr() {
	d.TrX = func(x float64) int {
		return d.GI.Pad + d.GI.SprWidth + int(float64(d.GI.GphWidth-2*d.GI.Pad)*(x-d.Xb[0])/(d.Xb[1]-d.Xb[0]))
	}
	d.TrY[0] = func(y float64) int {
		return d.GI.Pad + d.GI.Ceil + d.GI.SpL + int(float64(d.GI.GphHeight-2*d.GI.Pad)*(d.Yb[0][1]-y)/(d.Yb[0][1]-d.Yb[0][0]))
	}
	d.TrY[1] = func(y float64) int {
		return d.GI.Pad + d.GI.Ceil + 2*d.GI.SpL + d.GI.GphHeight + int(float64(d.GI.GphHeight-2*d.GI.Pad)*(d.Yb[1][1]-y)/(d.Yb[1][1]-d.Yb[1][0]))
	}
}

func (d *DspSim) Animate() {
	for i := 0; i < d.N; i++ {
		d.addFrame(i + 1)
	}
	fp, err := os.Create("graph.gif")
	if err != nil {
		panic(err)
	}
	err = gif.EncodeAll(fp, d.GI.G)
	if err != nil {
		panic(err)
	}
	fp.Close()
}

func (d *DspSim) addFrame(count int) {
	p := []color.Color{
		color.White,
		color.Black,
		color.NRGBA{255, 0, 0, 255},
		color.NRGBA{0, 0, 255, 255}}
	c := image.NewPaletted(image.Rect(0, 0, d.GI.Width, d.GI.Height), p)
	clear(c)
	d.axes(c)

	d.drawGraph(c, count)
	d.drawSpring(c, count)
	d.addLabel(c)

	d.GI.G.Image = append(d.GI.G.Image, c)
	d.GI.G.Delay = append(d.GI.G.Delay, d.GI.Delay)
}

func (s *DspSim) addLabel(c *image.Paletted) {
	d := &font.Drawer{
		Dst:  c,
		Src:  image.NewUniform(color.RGBA{0, 0, 0, 255}),
		Face: basicfont.Face7x13,
	}
	label := "Double Spring Simulator v1.0 "
	d.Dot = fixed.Point26_6{fixed.I(30), fixed.I(30)}
	d.DrawString(label)

	label = fmt.Sprintf("m1 = %.1f(kg), m2 = %.1f(kg), k1 = %.1f(N/m), k2 = %.1f", s.Val[0], s.Val[1], s.Val[2], s.Val[3])
	d.Dot = fixed.Point26_6{fixed.I(s.GI.SprWidth + 30), fixed.I(60)}
	d.DrawString(label)

	for j := 0; j < 2; j++ {
		for i := 0; i < 2; i++ {
			label = fmt.Sprintf("%.1f", s.Yb[j][i])
			d.Dot = fixed.Point26_6{fixed.I(s.GI.SprWidth + 20), fixed.I(s.TrY[j](s.Yb[j][i]))}
			d.DrawString(label)
		}
		label = fmt.Sprintf("%.1f", s.Xb[1])
		d.Dot = fixed.Point26_6{fixed.I(s.GI.Width - 50), fixed.I(s.TrY[j](0) - 5)}
		d.DrawString(label)
	}
}

func clear(c *image.Paletted) {
	width, height := c.Bounds().Max.X, c.Bounds().Max.Y
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			c.Set(i, j, c.Palette[0])
		}
	}
}

func (d *DspSim) axes(c *image.Paletted) {
	for j := 0; j < 2; j++ {
		for i := d.TrY[j](d.Yb[j][1]); i <= d.TrY[j](d.Yb[j][0]); i++ {
			c.Set(d.TrX(0), i, c.Palette[1])
		}
		for i := d.TrX(d.Xb[0]); i <= d.TrX(d.Xb[1]); i++ {
			c.Set(i, d.TrY[j](0), c.Palette[1])
		}

	}
}

func (d *DspSim) drawGraph(c *image.Paletted, count int) {
	for i := 0; i < count; i++ {
		x0, x1 := d.Xb[1]/float64(d.N)*float64(i), d.Xb[1]/float64(d.N)*float64(i+1)
		for j := 0; j < 2; j++ {
			y0, y1 := d.F[j](x0), d.F[j](x1)
			p0 := [2]int{d.TrX(x0), d.TrY[j](y0)}
			p1 := [2]int{d.TrX(x1), d.TrY[j](y1)}
			drawLine(c, p0, p1, j+2)
		}
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

func (d *DspSim) drawSpring(c *image.Paletted, count int) {
	d.addCeil(c)
	d.addSpring(c, count)
	d.addMass(c, count)
}

func (d *DspSim) addCeil(c *image.Paletted) {
	for i := d.GI.Pad; i < d.GI.SprWidth-d.GI.Pad; i++ {
		c.Set(i, d.GI.Ceil, c.Palette[1])
	}
}

func (d *DspSim) addSpring(c *image.Paletted, count int) {
	for i := 0; i < 100; i++ {
		for j := 0; j < 2; j++ {
			p0 := d.springXY(j, i, count)
			p1 := d.springXY(j, i+1, count)
			drawLine(c, p0, p1, 1)
		}
	}
}

func (d *DspSim) springXY(j, i, count int) [2]int {
	n := d.GI.CoilN
	r := d.GI.CoilR
	m := d.GI.CoilM
	y0 := d.F[0](d.Xb[1] / float64(d.N) * float64(count))
	y1 := d.F[1](d.Xb[1] / float64(d.N) * float64(count))
	st := -math.Pi / 2
	xi := int(float64(r) * math.Cos(st+(2*math.Pi*float64(n)-2*st)*float64(i)/100))
	xi += d.GI.Pad + d.GI.SprWidth/2
	var yi int
	switch j {
	case 0:
		l := d.TrY[0](y0) - d.GI.Ceil - d.GI.Mr - 2*r - 2*m
		yi = int(float64(r) * math.Sin(st+(2*math.Pi*float64(n)-2*st)*float64(i)/100))
		yi += int(float64(l) / float64(100) * float64(i))
		yi += m + r + d.GI.Ceil
	case 1:
		l := d.TrY[1](y1) - d.TrY[0](y0) - 2*d.GI.Mr - 2*r - 2*m
		yi = int(float64(r) * math.Sin(st+(2*math.Pi*float64(n)-2*st)*float64(i)/100))
		yi += d.TrY[0](y0) + d.GI.Mr + m + r + int(float64(l)/float64(100)*float64(i))
	}
	return [2]int{xi, yi}
}

func (d *DspSim) addMass(c *image.Paletted, count int) {
	for j := 0; j < 2; j++ {
		y := d.F[j](d.Xb[1] / float64(d.N) * float64(count))
		for i := 0; i < 100; i++ {
			p0 := d.circleXY(d.TrY[j](y), i)
			p1 := d.circleXY(d.TrY[j](y), i+1)
			drawLine(c, p0, p1, 1)
		}
	}
}

func (d *DspSim) circleXY(l, i int) [2]int {
	r := d.GI.Mr
	xi := int(float64(r) * math.Cos(2*math.Pi*float64(i)/100))
	xi += d.GI.SprWidth/2 + d.GI.Pad
	yi := int(float64(r) * math.Sin(2*math.Pi*float64(i)/100))
	yi += l
	return [2]int{xi, yi}
}
