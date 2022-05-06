package unforced

import (
	"fmt"
	"image/gif"
	"math"
)

type forced struct {
	V      []float64
	G      *gif.GIF
	GI     Ginfo
	F      func(float64) float64 // spring motion solution
	Fc     func(float64) float64 // complementary solution
	Fp     func(float64) float64 // complementary solution
	N      int                   // number of intervals
	Xb, Yb [2]float64
}

func NewSpSim() *forced {
	s := &forced{
		G: &gif.GIF{
			LoopCount: -1,
		},
		N: 100,
	}
	s.handleInput()
	s.findSolution()
	s.findXY()
	return s
}

func (s *forced) handleInput() {
	imsg := []string{
		"\n",
		"====================================\n",
		"This program animates the solution \n",
		"of the spring motion modeled by \n",
		"the following differential equation \n",
		"on the interval [0,T].\n",
		"mx'' + bx' + kx = sin(t), \n",
		"x(0) = x0, x'(0) = x1\n",
		"====================================\n",
		"\n",
	}
	for _, m := range imsg {
		fmt.Printf("%s", m)
	}

	msg := []string{
		"m: ",
		"b: ",
		"k: ",
		"x0: ",
		"x1: ",
		"T: ",
	}

	v := make([]float64, 6)

	for i, m := range msg {
		fmt.Printf("Enter the value of %s", m)
		fmt.Scanf("%f\n", &v[i])
	}
	s.V = v
}

func (s *forced) findSolution() {
	m, b, k := s.V[0], s.V[1], s.V[2]
	d := b*b - 4*m*k
	x0, x1 := s.V[3], s.V[4]
	var fc func(float64) float64
	var fp func(float64) float64
	var A, B float64
	var c1, c2 float64
	if b == 0 && m == k { // method of undetermined coefficients fails
		// xp(t) = A*t*cos(t) + B*t*sin(t)
		// xp'(t) = A*cos(t) - A*t*sin(t) + B*sin(t) + B*t*cos(t)
		// xp''(t) = -2*A*sin(t) - A*t*cos(t) + 2*B*cos(t) - B*t*sin(t)
		// (-m*A + B*b + k*A)*t*cos(t) + (2*B*m + b*A)*cos(t)
		// + (-B*m - b*A + k*B)*t*sin(t) + (-2*m*A + B*b)*sin(t) = sin(t)
		// b*A + 2*m*B = 0
		// -2*m*A + b*B = 1
		A = -2 * m / (b*b + 4*m*m)
		B = b / (b*b + 4*m*m)
		fp = func(t float64) float64 {
			return A*t*math.Cos(t) + B*t*math.Sin(t)
		}
	} else {
		// xp(t) = A*cos(t) + B*sin(t)
		// (-A*m + B*b + A*k)*cos(t) + (-B*m - A*b + B*k)*sin(t) = sin(t)
		// (-m+k)*A + b*B = 0
		// -b*A + (k-m)*B = 1
		A = (-b) / ((-m+k)*(k-m) + b*b)
		B = (-m + k) / ((-m+k)*(k-m) + b*b)
		fp = func(t float64) float64 {
			return A*math.Cos(t) + B*math.Sin(t)
		}
	}
	switch {
	case d > 0: // the characteristic polynomial has two real roots.
		l1, l2 := (-b+math.Sqrt(d))/(2*m), (-b-math.Sqrt(d))/(2*m)
		// x(t) = c1*exp(l1*t) + c2*exp(l2*t) + A*cos(t) + B*sin(t)
		// x(0) = x0, x'(0) = x1
		// c1 + c2 + A = x0
		// l1*c1 + l2*c2 + B = x1
		c1, c2 = ((x0-A)*l2-(x1-B))/(l2-l1), ((x1-B)-(x0-A)*l1)/(l2-l1)
		fc = func(t float64) float64 {
			return c1*math.Exp(l1*t) + c2*math.Exp(l2*t)
		}
	case d == 0:
		l1 := -b / (2 * m)
		// x(t) = c1*exp(l1*t) + c2*t*exp(l1*t) + A*cos(t) + B*sin(t)
		// c1 + A = x0
		// l1*c1 + c2 + B = x1
		c1, c2 = (x0 - A), (x1-B)-l1*(x0-A)
		fc = func(t float64) float64 {
			return c1*math.Exp(l1*t) + c2*t*math.Exp(l1*t)
		}
	case d < 0:
		// l1, l2 = re+im*i , re-im*i
		re, im := -b/(2*m), math.Sqrt(-d)/(2*m)
		if b == 0 && m == k {
			// x(t) + c1*cos(t) + c2*sin(t) + A*t*cos(t) + B*t*sin(t)
			// c1 = x0
			// c2 + A = x1
			c1, c2 = x0, x1-A
		} else {
			// x(t) = c1*exp(re*t)*cos(im*t) + c2*exp(re*t)*sin(im*t) + A*cos(t) + B*sin(t)
			// c1 + A = x0
			// re*c1 + im*c2 + B = x1
			c1, c2 = x0-A, (x1-B-(x0-A)*re)/im
		}
		fc = func(t float64) float64 {
			return c1*math.Exp(re*t)*math.Cos(im*t) + c2*math.Exp(re*t)*math.Sin(im*t)
		}
	}
	s.Fc = func(t float64) float64 {
		return -fc(t)
	}
	s.Fp = func(t float64) float64 {
		return -fp(t)
	}
	s.F = func(t float64) float64 {
		return -fp(t) - fc(t)
	}
}

func (s *forced) findXY() {
	s.Xb = [2]float64{0, s.V[5]}
	ymin, ymax := math.MaxFloat64, -math.MaxFloat64
	for i := 0; i <= s.N; i++ {
		x := s.V[5] / float64(s.N) * float64(i)
		y := s.F(x)
		yc := s.Fc(x)
		if y < ymin {
			ymin = y
		}
		if yc < ymin {
			ymin = yc
		}
		if y > ymax {
			ymax = y
		}
		if yc > ymax {
			ymax = yc
		}
	}
	if ymin > 0 {
		ymin = 0
	}
	if ymax < 0 {
		ymax = 0
	}
	s.Yb = [2]float64{ymin, ymax}
}
