package scatter

import (
	"math"
	"math/rand"
)

func (s Scatter) GradientDescent() Line {
	line := Line{
		D:     rand.Float64(),
		Theta: rand.Float64()}
	delta := 1e-4

	S := s.S(line)
	for {
		v := Gradient(s.S, line)
		norm := getNorm(v)
		for i := 0; i < 2; i++ {
			v[i] = v[i] / norm * delta
		}
		line.D -= v[0]
		line.Theta -= v[1]
		tS := s.S(line)
		if math.Abs(S-tS) < 1e-4 {
			break
		}
		S = tS
	}
	return line
}

func Gradient(S func(Line) float64, l Line) [2]float64 {
	f := func(x [2]float64) float64 {
		return S(Line{D: x[0], Theta: x[1]})
	}
	x := [2]float64{l.D, l.Theta}
	var v [2]float64
	for i := 0; i < 2; i++ {
		v[i] = PartialDerivative(x, f, i)
	}
	return v
}

func PartialDerivative(x [2]float64, f func([2]float64) float64, i int) float64 {
	delta := 1e-1
	var xp [2]float64
	for j := 0; j < 2; j++ {
		if j == i {
			xp[j] = x[j] + delta
		} else {
			xp[j] = x[j]
		}
	}
	r := (f(xp) - f(x)) / delta
	for {
		delta = delta / 10.0
		for j := 0; j < 2; j++ {
			if j == i {
				xp[j] = x[j] + delta
			} else {
				xp[j] = x[j]
			}
		}
		rr := (f(xp) - f(x)) / delta
		if math.Abs(r-rr) < 1e-3 {
			return r
		} else {
			r = rr
		}
	}
}
