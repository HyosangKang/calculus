package mvf

import (
	"image/color"
	"math"
	"multivariable/vector"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func Graph2to1(f func([]float64) []float64, lim [][2]float64, filename string, N int) {
	p := plot.New()
	// drawing y-axis
	pts := plotter.XYs{
		{X: 0, Y: 0},
		{X: 1, Y: 0}}
	line, _ := plotter.NewLine(pts)
	p.Add(line)

	// drawing z-axis
	pts = plotter.XYs{
		{X: 0, Y: 0},
		{X: 0, Y: 1}}
	line, _ = plotter.NewLine(pts)
	p.Add(line)

	// drawing x-axis
	pts = plotter.XYs{
		{X: 0, Y: 0},
		{X: -1 / math.Sqrt(2), Y: -1 / math.Sqrt(2)}}
	line, _ = plotter.NewLine(pts)
	p.Add(line)
	p.HideAxes()

	for i := 0; i <= N; i++ {
		x := lim[0][0] + (lim[0][1]-lim[0][0])*float64(i)/float64(N)
		y := make([]float64, N+1)
		z := make([]float64, N+1)
		for j := 0; j <= N; j++ {
			y[j] = lim[1][0] + (lim[1][1]-lim[1][0])*float64(j)/float64(N)
			temp := f([]float64{x, y[j]})
			z[j] = temp[0]
		}
		for j := 0; j < N; j++ {
			if math.IsNaN(z[j]) || math.IsNaN(z[j+1]) {
				continue
			}
			pts = plotter.XYs{
				{X: -x/math.Sqrt(2) + y[j],
					Y: -x/math.Sqrt(2) + z[j]},
				{X: -x/math.Sqrt(2) + y[j+1],
					Y: -x/math.Sqrt(2) + z[j+1]}}
			line, _ = plotter.NewLine(pts)
			p.Add(line)
		}
	}

	for j := 0; j <= N; j++ {
		y := lim[1][0] + (lim[1][1]-lim[1][0])*float64(j)/float64(N)
		x := make([]float64, N+1)
		z := make([]float64, N+1)
		for i := 0; i <= N; i++ {
			x[i] = lim[0][0] + (lim[0][1]-lim[0][0])*float64(i)/float64(N)
			temp := f([]float64{x[i], y})
			z[i] = temp[0]
		}
		for i := 0; i < N; i++ {
			if math.IsNaN(z[i]) || math.IsNaN(z[i+1]) {
				continue
			}
			pts = plotter.XYs{
				{X: -x[i]/math.Sqrt(2) + y,
					Y: -x[i]/math.Sqrt(2) + z[i]},
				{X: -x[i+1]/math.Sqrt(2) + y,
					Y: -x[i+1]/math.Sqrt(2) + z[i+1]}}
			line, _ = plotter.NewLine(pts)
			p.Add(line)
		}
	}
	p.Save(400, 400, filename)
}

func Graph2to1Spherical(f func([]float64) []float64, lim [][2]float64, filename string, N int, t []float64) {
	e := make([]vector.Vector, 3)
	e[0] = vector.Vector([]float64{1, 0, 0})
	e[1] = vector.Vector([]float64{0, 1, 0})
	e[2] = vector.Vector([]float64{0, 0, 1})

	s := make([]vector.Vector, 3)
	for i := 0; i < 3; i++ {
		s[i] = e[i].SphericalTransTorm(t[0], t[1])
	}

	p := plot.New()
	p.HideAxes()

	r := make([]vector.Vector, 3)
	for i := 0; i < 3; i++ {
		r[i] = e[i].PlaneProjection(s[0])
		pts := plotter.XYs{
			{X: 0, Y: 0},
			{X: s[1].Inner(r[i]), Y: s[2].Inner(r[i])},
		}
		line, _ := plotter.NewLine(pts)
		line.LineStyle.Color = color.RGBA{255, 0, 0, 255}
		p.Add(line)
	}
	for i := 0; i <= N; i++ {
		x := lim[0][0] + (lim[0][1]-lim[0][0])*float64(i)/float64(N)
		y := make([]float64, N+1)
		z := make([]float64, N+1)
		for j := 0; j <= N; j++ {
			y[j] = lim[1][0] + (lim[1][1]-lim[1][0])*float64(j)/float64(N)
			temp := f([]float64{x, y[j]})
			z[j] = temp[0]
		}
		for j := 0; j < N; j++ {
			if math.IsNaN(z[j]) || math.IsNaN(z[j+1]) {
				continue
			}

			pts := make(plotter.XYs, 2)
			for k := 0; k < 2; k++ {
				v := vector.Vector([]float64{x, y[j+k], z[j+k]})
				pts[k] = plotter.XY{
					X: v.Inner(s[1]),
					Y: v.Inner(s[2]),
				}
			}
			line, _ := plotter.NewLine(pts)
			p.Add(line)
		}
	}

	for j := 0; j <= N; j++ {
		y := lim[1][0] + (lim[1][1]-lim[1][0])*float64(j)/float64(N)
		x := make([]float64, N+1)
		z := make([]float64, N+1)
		for i := 0; i <= N; i++ {
			x[i] = lim[0][0] + (lim[0][1]-lim[0][0])*float64(i)/float64(N)
			temp := f([]float64{x[i], y})
			z[i] = temp[0]
		}
		for i := 0; i < N; i++ {
			if math.IsNaN(z[i]) || math.IsNaN(z[i+1]) {
				continue
			}
			pts := make(plotter.XYs, 2)
			for k := 0; k < 2; k++ {
				v := vector.Vector([]float64{x[i+k], y, z[i+k]})
				pts[k] = plotter.XY{
					X: v.Inner(s[1]),
					Y: v.Inner(s[2]),
				}
			}
			line, _ := plotter.NewLine(pts)
			p.Add(line)
		}
	}

	p.Save(400, 400, filename)
}

func Graph1to3(c func([]float64) []float64, lim [2]float64, filename string, N int) {
	p := plot.New()
	// drawing y-axis
	pts := plotter.XYs{
		{X: 0, Y: 0},
		{X: 1, Y: 0}}
	line, _ := plotter.NewLine(pts)
	p.Add(line)

	// drawing z-axis
	pts = plotter.XYs{
		{X: 0, Y: 0},
		{X: 0, Y: 1}}
	line, _ = plotter.NewLine(pts)
	p.Add(line)

	// drawing x-axis
	pts = plotter.XYs{
		{X: 0, Y: 0},
		{X: -1 / math.Sqrt(2), Y: -1 / math.Sqrt(2)}}
	line, _ = plotter.NewLine(pts)
	p.Add(line)
	p.HideAxes()

	pts = make(plotter.XYs, N+1)
	for i := 0; i <= N; i++ {
		t := lim[0] + (lim[1]-lim[0])/float64(N)*float64(i)
		d := c([]float64{t})
		x := d[0]
		y := d[1]
		z := d[2]
		pts[i] = plotter.XY{
			X: -x/math.Sqrt(2) + y,
			Y: -x/math.Sqrt(2) + z,
		}
	}
	line, _ = plotter.NewLine(pts)
	p.Add(line)

	p.Save(400, 400, filename)
}
