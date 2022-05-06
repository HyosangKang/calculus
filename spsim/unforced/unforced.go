package forced

import (
	"fmt"
	"image/gif"
	"math"
)

type unforced struct {
	v      []float64
}

func NewSpSim() *SpSim {
	s := &SpSim{
		G: &gif.GIF{
			LoopCount: 1,
		},
		N: 100,
	}
	s.handleInput()
	s.findSolution()
	s.findXY()
	return s
}

func (s *SpSim) handleInput() {
	imsg := []string{
		"\n",
		"====================================\n",
		"This program animates the solution \n",
		"of the spring motion modeled by \n",
		"the following differential equation \n",
		"on the interval [0,T].\n",
		"mx'' + bx' + kx = 0, \n",
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

func (s *SpSim) findSolution() {
	m, b, k := s.V[0], s.V[1], s.V[2]
	d := b*b - 4*m*k
	x0, x1 := s.V[3], s.V[4]
	var f func(float64) float64
	switch {
	case d > 0: // the characteristic polynomial has two real roots.
		l1, l2 := (-b+math.Sqrt(d))/(2*m), (-b-math.Sqrt(d))/(2*m)
		// x(t) = c1*exp(l1*t) + c2*exp(l2*t)
		// x(0) = x0, x'(0) = x1
		// c1 + c2 = x0
		// l1*c1 + l2*c2 = x1
		c1, c2 := (x0*l2-x1)/(l2-l1), (x1-x0*l1)/(l2-l1)
		f = func(t float64) float64 {
			return c1*math.Exp(l1*t) + c2*math.Exp(l2*t)
		}
	case d == 0:
		l1 := -b / (2 * m)
		// x(t) = c1*exp(l1*t) + c2*t*exp(l1*t)
		// c1 = x0
		// l1*c1 + c2 = x1
		c1, c2 := x0, x1-l1*x0
		f = func(t float64) float64 {
			return c1*math.Exp(l1*t) + c2*t*math.Exp(l1*t)
		}
	case d < 0:
		// l1, l2 = re+im*i , re-im*i
		re, im := -b/(2*m), math.Sqrt(-d)/(2*m)
		// x(t) = c1*exp(re*t)*cos(im*t) + c2*exp(re*t)*sin(im*t)
		// c1 = x0
		// re*c1 + im*c2 = x1
		c1, c2 := x0, (x1-x0*re)/im
		f = func(t float64) float64 {
			return c1*math.Exp(re*t)*math.Cos(im*t) + c2*math.Exp(re*t)*math.Sin(im*t)
		}
	}
	s.F = func(t float64) float64 {
		return -f(t)
	}
}

func (s *SpSim) findXY() {
	s.Xb = [2]float64{0, s.V[5]}
	ymin, ymax := math.MaxFloat64, -math.MaxFloat64
	for i := 0; i <= s.N; i++ {
		x := s.V[5] / float64(s.N) * float64(i)
		y := s.F(x)
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
	s.Yb = [2]float64{ymin, ymax}
}
