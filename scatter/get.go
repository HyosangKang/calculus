package scatter

import "math"

func getTheta(u [2]float64) float64 {
	if isZero(u[0]) {
		if u[1] > 0 {
			return 0.5 * math.Pi
		} else {
			return -0.5 * math.Pi
		}
	} else {
		if u[0] > 0 {
			return math.Atan(u[1] / u[0])
		} else {
			return math.Atan(u[1]/u[0]) + math.Pi
		}
	}
}

func isZero(f float64) bool {
	return math.Abs(f) < 1e-10
}

func getNorm(v [2]float64) float64 {
	var n float64
	for i := 0; i < 2; i++ {
		n += v[i] * v[i]
	}
	return math.Sqrt(n)
}

func getMaxAbs(v [2]float64) float64 {
	var s float64
	for i := 0; i < 2; i++ {
		t := math.Abs(v[i])
		if s < t {
			s = t
		}
	}
	return s
}
