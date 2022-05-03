package solver

func Euler(fs []func(float64, []float64) float64, r [2]float64, y0 []float64, nn int) ([]float64, [][]float64) {
	if len(y0) != len(fs) {
		panic("The sizes of outputs do not match!")
	}
	n := len(y0)
	h := (r[1]-r[0])/float64(nn)
	tl := []float64{r[0]}
	yl := make([][]float64, n)
	for i := 0; i < n; i++ {
		yl[i] = []float64{y0[i]}
	}
	ys := make([]float64, n)
	copy(ys, y0)
	for j := 0 ; j < nn ; j ++ {
		t := r[0] + float64(j)*h 
		k := make([]float64, n)
		tp := t
		yp := make([]float64, len(ys))
		copy(yp, ys)
		for i, f := range fs {
			k[i] = f(tp, yp)
		}
		for i := 0; i < n; i++ {
			ys[i] += h * k[i]
			yl[i] = append(yl[i], ys[i])
		}
		tl = append(tl, t+h)
	}
	return tl, yl
}
