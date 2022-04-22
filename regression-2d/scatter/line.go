package scatter

import (
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type Line struct {
	D, Theta float64
}

func NewSampleLine() Line {
	return Line{
		D:     1,
		Theta: float64(5) * math.Pi / float64(4),
	}
}

func (line Line) Graph(p *plot.Plot) {
	t := [2]float64{-10, 10}
	l := func(t float64) [2]float64 {
		x := line.D*math.Sin(line.Theta) + t*math.Cos(line.Theta)
		y := -line.D*math.Cos(line.Theta) + t*math.Sin(line.Theta)
		return [2]float64{x, y}
	}
	pts := plotter.XYs{
		{X: l(t[0])[0], Y: l(t[0])[1]},
		{X: l(t[1])[0], Y: l(t[1])[1]},
	}
	ld, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}
	p.Add(ld)
}
