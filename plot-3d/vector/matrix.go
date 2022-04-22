package vector

import (
	"log"
	"math"
)

type Matrix [][]float64

func NewMatrix(r, c int) Matrix {
	m := make([][]float64, r)
	for i := 0; i < r; i++ {
		m[i] = make([]float64, c)
	}
	return m
}

func (m Matrix) ToVector() Vector {
	if len(m[0]) != 1 {
		log.Panic("Invalid conversion from Matrix to Vectro")
	}
	v := NewVector(len(m))
	for i := 0 ; i < len(m); i ++ {
		v[i] = m[i][0]
	}
	return v
}

func (m Matrix) Mul(n Matrix) Matrix {
	if len(m[0]) != len(n) {
		log.Panic("Invalid matrix multiplication")
	}
	l := make([][]float64, len(m))
	for i := 0; i < len(m); i++ {
		l[i] = make([]float64, len(n[0]))
		for j := 0; j < len(n[0]); j++ {
			var sum float64 = 0
			for k := 0; k < len(m[0]); k++ {
				sum += m[i][k] * n[k][j]
			}
			l[i][j] = sum
		}
	}
	return l
}

func SphericalRotation(t float64, axis string) Matrix {
	m := NewMatrix(3, 3)
	switch axis {
	case "x":
	case "y":
		m[0][0] = math.Cos(t)
		m[0][2] = -math.Sin(t)
		m[1][1] = 1
		m[2][0] = math.Sin(t)
		m[2][2] = math.Cos(t)
	case "z":
		m[0][0] = math.Cos(t)
		m[0][1] = -math.Sin(t)
		m[1][0] = math.Sin(t)
		m[1][1] = math.Cos(t)
		m[2][2] = 1
	}
	return m
}
