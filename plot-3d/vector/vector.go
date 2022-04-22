package vector

import (
	"log"
	"math"
)

type Vector []float64

func NewVector(l int) Vector {
	return make([]float64, l)
}

func (v Vector) Norm() float64 {
	var sum float64 = 0
	for _, e := range v {
		sum += e * e
	}
	return math.Sqrt(sum)
}

func (v1 Vector) Inner(v2 Vector) float64 {
	if len(v1) != len(v2) {
		log.Panic("Invalid inner product")
	}
	var sum float64 = 0
	for i := 0; i < len(v1); i++ {
		sum += v1[i] * v2[i]
	}
	return sum
}

func (v Vector) Scalar(c float64) Vector {
	w := make([]float64, len(v))
	for i, e := range v {
		w[i] = c * e
	}
	return w
}

func (v1 Vector) Projection(v2 Vector) Vector {
	if len(v1) != len(v2) {
		log.Panic("Invalid projection")
	}
	w := make([]float64, len(v1))
	r := v1.Inner(v2)
	n1 := v1.Norm()
	n2 := v2.Norm()
	for i := 0; i < len(v1); i++ {
		w[i] = v1[i] - r*v1[i]/n1/n2
	}
	return w
}

func (v Vector) PlaneProjection(n Vector) Vector {
	c := v.Inner(n) / n.Norm()
	w := make([]float64, len(v))
	for i := 0; i < len(v); i++ {
		w[i] = v[i] - c*n[i]
	}
	return w
}

func (v Vector) SphericalTransTorm(t, p float64) Vector {
	r1 := SphericalRotation(t, "y")
	r2 := SphericalRotation(p, "z")
	return r2.Mul(r1).Mul(v.ToMatrix()).ToVector()
}

func (v Vector) ToMatrix() Matrix {
	m := NewMatrix(len(v), 1)
	for i := 0; i < len(v); i++ {
		m[i][0] = v[i]
	}
	return m
}
