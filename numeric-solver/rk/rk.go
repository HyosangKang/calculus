package rk

func RungeKutta(fs []func(float64, []float64) float64, r [2]float64, y0 []float64, h float64) ([]float64, [][]float64) {
	if len(y0) != len(fs) {
		panic("The sizes of outputs do not match!")
	}
	n := len(y0)
	tl := []float64{r[0]}
	yl := make([][]float64, n)
	for i := 0; i < n; i++ {
		yl[i] = []float64{y0[i]}
	}
	ys := make([]float64, n)
	copy(ys, y0)
	for t := r[0]; t <= r[1]-h; t += h {
		k := make([][]float64, 4)
		for j := 0; j < 4; j++ {
			k[j] = make([]float64, len(y0))
			tp := t
			yp := make([]float64, len(ys))
			copy(yp, ys)
			switch j {
			case 1:
				tp += h / 2
				for i := 0; i < n; i++ {
					yp[i] += h * k[0][i] / 2
				}
			case 2:
				tp += h / 2
				for i := 0; i < n; i++ {
					yp[i] += h * k[1][i] / 2
				}
			case 3:
				tp += h
				for i := 0; i < n; i++ {
					yp[i] += h * k[2][i]
				}
			}
			for i, f := range fs {
				k[j][i] = f(tp, yp)
			}
		}
		for i := 0; i < n; i++ {
			ys[i] += h * (k[0][i] + 2*k[1][i] + 2*k[2][i] + k[3][i]) / 6
			yl[i] = append(yl[i], ys[i])
		}
		tl = append(tl, t+h)
	}
	return tl, yl
}
