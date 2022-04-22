package mvf

import (
	"fmt"
	"log"
)

type mvf struct {
	formula  func([]float64) []float64
	domain   [][2]float64
	ind, dep int
}

func NewMvf(f func([]float64) []float64, dom [][2]float64, in, dn int) *mvf {
	if len(dom) != in {
		log.Panic("Invalid domain")
	}
	return &mvf{
		formula: f,
		domain:  dom,
		ind:     in,
		dep:     dn,
	}
}

func (f *mvf) Graph(filename string, N int) {
	if f.ind == 2 && f.dep == 1 {
		Graph2to1(f.formula, f.domain, filename, N)
	} else if f.ind == 1 && f.dep == 3 {
		Graph1to3(f.formula, f.domain[0], filename, N)
	} else {
		fmt.Println("Cannot graph the function")
	}
}

func (f *mvf) GraphSpherical(filename string, N int, t, p float64) {
	if f.ind == 2 && f.dep == 1 {
		Graph2to1Spherical(f.formula, f.domain, filename, N, []float64{t, p})
	} else {
		fmt.Println("Cannot graph the function")
	}
}
