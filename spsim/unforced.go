package spsim

import (
	"fmt"
	"math"
)

func UnforcedSim() ([2]float64, func(float64) float64) {
	v := handleInputUnforced()
	return [2]float64{0, v[5]}, findSolutionUnforced(v)
}

func handleInputUnforced() []float64 {
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
	return v
}

func findSolutionUnforced(v []float64) func(float64) float64 {
	m, b, k := v[0], v[1], v[2]
	d := b*b - 4*m*k
	x0, x1 := v[3], v[4]
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
	return f
}
