package scatter

import (
	"math"
	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type Scatter [][2]float64

func NewSample() Scatter {
	s := make([][2]float64, 10)

	line := NewSampleLine()
	l := func(t float64) [2]float64 {
		x := line.D*math.Sin(line.Theta) + t*math.Cos(line.Theta)
		y := -line.D*math.Cos(line.Theta) + t*math.Sin(line.Theta)
		return [2]float64{x, y}
	}

	ta := make([]float64, 10)
	for i := 0; i < 10; i++ {
		ta[i] = float64(20)*rand.Float64() - float64(10)
	}

	for i := 0; i < 10; i++ {
		xy := l(ta[i])
		x := xy[0] + rand.Float64() - 0.5
		y := xy[1] + rand.Float64() - 0.5
		s[i] = [2]float64{x, y}
	}

	return s
}

func NewFromFile(filename string) Scatter {
	s := make([][2]float64, 10)
	return s
}

func (s Scatter) Regression() Line {
	var a, b, x, y float64
	for i := 0; i < len(s); i++ {
		a += s[i][0] * s[i][1]
		b += s[i][0]*s[i][0] - s[i][1]*s[i][1]
		x += s[i][0]
		y += s[i][1]
	}
	n := float64(len(s))
	u := [2]float64{2*a + 2*x*y/n, b - (x*x-y*y)/n}
	theta := (getTheta(u) + math.Pi/float64(2)) / float64(2)
	d := x/n*math.Sin(theta) - y/n*math.Cos(theta)

	return Line{
		D: d,
		Theta: theta,
	} 
}

func (s Scatter) S(line Line) float64 {
	var sum float64
	for i := 0; i < len(s); i++ {
		a := math.Sin(line.Theta)*s[i][0] - math.Cos(line.Theta)*s[i][1] - line.D
		sum += a * a
	}
	return sum
}

func (s Scatter) Graph(p *plot.Plot) {
	pts := make(plotter.XYs, len(s))
	for i := 0; i < len(s); i++ {
		pts[i].X = s[i][0]
		pts[i].Y = s[i][1]
	}
	sc, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}
	p.Add(sc)
}
