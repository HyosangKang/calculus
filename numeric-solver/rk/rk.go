package runge_kutta

func RungeKutta(fs []func(float64, []float64) float64, r [2]float64, y0 []float64, h float64) ([]float64, [][]float64) {
	if len(y0) != len(fs) {
		panic("The sizes of outputs do not match!")
	}
	tl := []float64{r[0]}
	yl := [][]float64{y0}
	ys := y0
	for t := r[0]; t <= r[1]-h; t += h {
		k := make([][]float64, 4)
		for j := 0; j < 4; j++ {
			k[j] = make([]float64, len(y0))
			tp := t
			yp := make([]float64, len(ys))
			switch j {
			case 0:
				for l, y := range ys {
					yp[l] = y
				}
			case 1:
				tp += h / 2
				for l, y := range ys {
					yp[l] = y + h*k[0][l]/2
				}
			case 2:
				tp += h / 2
				for l, y := range ys {
					yp[l] = y + h*k[1][l]/2
				}
			case 3:
				tp += h
				for l, y := range ys {
					yp[l] = y + h*k[2][l]
				}
			}
			for i, f := range fs {
				k[j][i] = f(tp, yp)
			}
		}
		ny := []float64{}
		for i, y := range ys {
			ys[i] = y + h*(k[0][i]+2*k[1][i]+2*k[2][i]+k[3][i])/6
			ny = append(ny, ys[i])
		}
		tl = append(tl, t+h)
		yl = append(yl, ny)
	}
	return tl, yl
}
