package esim

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

type Esim struct {
	L, R1, R2, C                       float64
	I                                  []func(float64) float64
	Xb                                 [2]float64
	TrX                                func(float64) int
	Yb                                 [][2]float64
	TrY                                []func(float64) int
	G                                  *gif.GIF
	Delay, N                           int
	Width, GphWidth, Height, GphHeight int
	Pad, Lpad, Tpad                    int
}

func NewEsim() *Esim {
	e := &Esim{
		Width:     600,
		GphWidth:  570,
		Height:    600,
		GphHeight: 190,
		Delay:     1,
		N:         100,
		Pad:       10,
		Lpad:      30,
		Tpad:      30,
	}
	e.handleInput()
	return e
}

func (e *Esim) handleInput() {
	msg := []string{
		"==========================================\n",
		"The program simulates electric currents on\n",
		"the following electric circuits for [0, T]\n",
		"  --L--- --R2--\n",
		" |      |      |+\n",
		" |      C      V\n",
		" |      |      |-\n",
		"  --R1-- ------\n",
		"   -i1->  -i2->\n",
		"The voltage(V) is given by sin(t)\n",
		"and the initial conditions are\n",
		"i1(0)=0, i1'(0)=0, i2(0)=0, i2'(0)=0\n",
		"=========================================\n",
	}
	for _, m := range msg {
		fmt.Printf("%s", m)
	}
	val := make([]float64, 5)
	msg = []string{
		"L",
		"R1",
		"R2",
		"C",
		"T",
	}
	for i, m := range msg {
		fmt.Printf("Enter the value of %s: ", m)
		fmt.Scanf("%f\n", &val[i])
	}
	for i, v := range val {
		if v <= 0 {
			panic(fmt.Sprintf("Invalid value of %s", msg[i]))
		}
	}
	e.L = val[0]
	e.R1 = val[1]
	e.R2 = val[2]
	e.C = val[3]
	e.Xb = [2]float64{0, val[4]}
	e.G = &gif.GIF{
		LoopCount: -1,
	}
}

func (e *Esim) Graph(fn string) {
	e.Solve()
	e.FindBounds()
	e.FindTr()
	for i := 0; i < e.N; i++ {
		e.AddFrame(i + 1)
	}
	e.Save(fn)
}

func (e *Esim) Solve() {
	L, R1, R2, C := e.L, e.R1, e.R2, e.C
	// Initial setup
	// -L*i1' - R1*i1 - 1/C*int(i1-i2) = 0
	// sin(t) - R2*i2 - 1/C*int(i2-i1) = 0

	// Simplied
	// L*i1'' + R1*i1' + 1/C*(i1-i2) = 0
	// R2*i2' + 1/C*(i2-i1) = cos(t)

	// After Laplace transform
	// L*s^2*I1 + R1*s*I1 + 1/C*(I1-I2) = 0
	// R2*s*I2 + 1/C*(I2-I1) = s/(s^2+1)

	// I1 = (a*s+b)/(s^2+1)
	//      + (c*s+d)/(C*L*R2*s^2 + (L+C*R1*R2)*s + R1+R2)
	// c = -C*L*R2*a
	// d = 1 - (R1+R2)*b
	// (R1+R2-C*L*R2)*a + (L+C*R1*R2)*b = 0
	// (L+C*R1*R2)*a + (C*L*R2-R1-R2)*b = -1
	a1, b1 := Cramer([2][2]float64{
		{R1 + R2 - C*L*R2, L + C*R1*R2},
		{L + C*R1*R2, C*L*R2 - R1 - R2}}, [2]float64{0, -1})
	c1, d1 := -C*L*R2*a1, 1-(R1+R2)*b1
	// I1 = a * s/(s^2+1) + b * 1/(s^2+1)
	// 		+ (c/C*L*R2) * ((s+alpha)/(s+alpha)^2+-beta^2)
	// 		+ gamma/(C*L*R2) * 1/((s+alpha)^2+-beta^2)
	// alpha = (L+C*R1*R2)/(2*C*L*R2)
	// +-beta^2 = (R1+R2)/(C*L*R2) - (L+C*R1*R2)^2/(4*(C*L*R2)^2)
	// gamma = d - c*(L+C*R1*R2)/(2*C*L*R2)
	alpha1 := (L + C*R1*R2) / (2 * C * L * R2)
	gamma1 := d1 - c1*(L+C*R1*R2)/(2*C*L*R2)
	bb1 := (R1+R2)/(C*L*R2) - (L+C*R1*R2)*(L+C*R1*R2)/(4*(C*L*R2)*(C*L*R2))
	beta1 := math.Sqrt(math.Abs(bb1))
	var f1, f2 func(float64) float64
	if bb1 < 0 {
		f1, f2 = math.Cosh, math.Sinh
	} else {
		f1, f2 = math.Cos, math.Sin
	}
	e.I = append(e.I, func(t float64) float64 {
		return a1*math.Cos(t) + b1*math.Sin(t) + c1/(C*L*R2)*math.Exp(-alpha1*t)*f1(beta1*t) + gamma1/(C*L*R2)*math.Exp(-alpha1*t)*f2(beta1*t)
	})
	// I2 = a * (s/(s^2+1)) + b * (1/(s^2+1))
	// 		+ (c/L*R2) * (s+alpha)/((s+alpha)^2+-beta^2)
	// 		+ (gamma/L*R2) * 1/(s+alpha)^2+-beta^2
	a2, b2 := Cramer([2][2]float64{
		{L/C + R1*R2, L*R2 - (R1+R2)/C},
		{(R1+R2)/C - L*R2, L/C + R1*R2}}, [2]float64{L - 1/C, R1})
	c2, d2 := -L*R2*a2, 1/C-(R1+R2)/C*b2
	// alpha = (L/C+R1*R2)/(2*L*R2)
	// +-beta^2 = -L*R2*(L/C+R1*R2)^2/(4*(L*R2)^2) + (R1+R2)/C
	// gamma = (d - c*alpha)/beta
	alpha2 := (L/C + R1*R2) / (2 * L * R2)
	bb2 := -(L/C+R1*R2)*(L/C+R1*R2)/(4*(L*R2)*(L*R2)) + (R1+R2)/C/(L*R2)
	beta2 := math.Sqrt(math.Abs(bb2))
	gamma2 := (d2 - c2*alpha2) / beta2

	if bb2 < 0 {
		f1, f2 = math.Cosh, math.Sinh
	} else {
		f1, f2 = math.Cos, math.Sin
	}
	e.I = append(e.I, func(t float64) float64 {
		return a2*math.Cos(t) + b2*math.Sin(t) + c2/(L*R2)*math.Exp(-alpha2*t)*f1(beta2*t) + gamma2/(L*R2)*math.Exp(-alpha2*t)*f2(beta2*t)
	})
}

func Cramer(m [2][2]float64, b [2]float64) (float64, float64) {
	d := m[0][0]*m[1][1] - m[0][1]*m[1][0]
	x0 := (b[0]*m[1][1] - b[1]*m[0][1]) / d
	x1 := (m[0][0]*b[1] - b[0]*m[1][0]) / d
	return x0, x1
}

func (e *Esim) FindBounds() {
	for _, f := range e.I {
		ymin, ymax := math.MaxFloat64, -math.MaxFloat64
		for j := 0; j < e.N; j++ {
			t := e.Xb[1] / float64(e.N) * float64(j)
			y := f(t)
			if y > ymax {
				ymax = y
			}
			if y < ymin {
				ymin = y
			}
		}
		if ymax < 0 {
			ymax = 0
		}
		if ymin > 0 {
			ymin = 0
		}
		e.Yb = append(e.Yb, [2]float64{ymin, ymax})
	}
}

func (e *Esim) FindTr() {
	e.TrX = func(x float64) int {
		return e.Lpad + e.Pad + int(float64(e.GphWidth-2*e.Pad)*(x-e.Xb[0])/(e.Xb[1]-e.Xb[0]))
	}
	e.TrY = append(e.TrY, func(y float64) int {
		return e.Tpad + e.Pad + e.GphHeight + int(float64(e.GphHeight-2*e.Pad)*(e.Yb[0][1]-y)/(e.Yb[0][1]-e.Yb[0][0]))
	})
	e.TrY = append(e.TrY, func(y float64) int {
		return e.Tpad + e.Pad + e.GphHeight*2 + int(float64(e.GphHeight-2*e.Pad)*(e.Yb[1][1]-y)/(e.Yb[1][1]-e.Yb[1][0]))
	})
}

func (e *Esim) AddFrame(count int) {
	p := []color.Color{color.White, color.Black}
	c := image.NewPaletted(image.Rect(0, 0, e.Width, e.Height), p)
	clear(c)
	e.addAxes(c)
	e.addGraph(c, count)
	e.addLabel(c)
	e.G.Image = append(e.G.Image, c)
	e.G.Delay = append(e.G.Delay, e.Delay)
}

func clear(c *image.Paletted) {
	width, height := c.Bounds().Max.X, c.Bounds().Max.Y
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			c.Set(i, j, c.Palette[0])
		}
	}
}

func (e *Esim) addAxes(c *image.Paletted) {
	for i, trY := range e.TrY {
		for j := e.TrX(e.Xb[0]); j <= e.TrX(e.Xb[1]); j++ {
			c.Set(j, trY(0), c.Palette[1])
		}
		for j := trY(e.Yb[i][1]); j <= trY(e.Yb[i][0]); j++ {
			c.Set(e.TrX(0), j, c.Palette[1])
		}
	}
}

func (e *Esim) addGraph(c *image.Paletted, count int) {
	for i, f := range e.I {
		for j := 0; j < count; j++ {
			x0, x1 := e.Xb[1]*float64(j)/float64(e.N), e.Xb[1]*float64(j+1)/float64(e.N)
			y0, y1 := f(x0), f(x1)
			drawLine(c, [2]int{e.TrX(x0), e.TrY[i](y0)}, [2]int{e.TrX(x1), e.TrY[i](y1)})
		}
	}
}

func drawLine(c *image.Paletted, p0, p1 [2]int) {
	x, y := 0, 0
	n := findMax(p0, p1)
	for i := 0; i <= n; i++ {
		nx := p0[0] + int(float64((p1[0]-p0[0])*i)/float64(n))
		ny := p0[1] + int(float64((p1[1]-p0[1])*i)/float64(n))
		for nx == x && ny == y {
			continue
		}
		x, y := nx, ny
		c.Set(x, y, c.Palette[1])
	}
}

func (e *Esim) addLabel(c *image.Paletted) {
	d := &font.Drawer{
		Dst:  c,
		Src:  image.NewUniform(color.RGBA{0, 0, 0, 255}),
		Face: basicfont.Face7x13,
	}
	d.Dot = fixed.Point26_6{fixed.I(15), fixed.I(15)}
	d.DrawString("Electric circuit Simulator v1.0 ")

	texts := []string{
		"  ---L--- --R2--",
		" |      |      |+",
		" |      C      V sin(t)",
		" |      |      |-",
		"  --R1-- ------",
		"   -i1->  -i2->",
	}
	for i, t := range texts {
		d.Dot = fixed.Point26_6{fixed.I(e.Width / 2), fixed.I(45 + 20*i)}
		d.DrawString(t)
	}

	d.Dot = fixed.Point26_6{fixed.I(15), fixed.I(e.Tpad + e.GphHeight + 15)}
	d.DrawString("I1")
	d.Dot = fixed.Point26_6{fixed.I(15), fixed.I(e.Tpad + e.GphHeight*2 + 15)}
	d.DrawString("I2")
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

func (e *Esim) Save(fn string) {
	fp, _ := os.Create(fn)
	gif.EncodeAll(fp, e.G)
	fp.Close()
}
