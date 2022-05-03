package main

import "fmt"

var eps = 1e-6

func main() {
	fmt.Println("1. exp, 2. cos, 3. sin")
	fmt.Printf("Choose the function(#): ")

	var n int
	fmt.Scanf("%d\n", &n)
	fmt.Printf("Input x: ")

	var x float64
	fmt.Scanf("%f\n", &x)

	switch n {
	case 1:
		fmt.Printf("The value of exp(%.1f) is %.3f\n", x, Exp(x))
	case 2:
		fmt.Printf("The value of cos(%.1f) is %.3f\n", x, Cos(x))
	case 3:
		fmt.Printf("The value of sin(%.1f) is %.3f\n", x, Sin(x))
	}
}

func Exp(x float64) float64 {
	var y, t float64 = 1, 1
	for i := 1; ; i++ {
		t *= x / float64(i)
		y += t
		if t < eps {
			break
		}
	}
	return y
}

func Cos(x float64) float64 {
	var y, t float64 = 1, 1
	for i := 1; ; i++ {
		t *= x * x / float64(2*i*(2*i-1))
		if i%2 == 0 {
			y += t
		} else {
			y -= t
		}
		if t < eps {
			break
		}
	}
	return y
}

func Sin(x float64) float64 {
	var y, t float64 = x, x
	for i := 1; ; i++ {
		t *= x * x / float64(2*i*(2*i+1))
		if i%2 == 0 {
			y += t
		} else {
			y -= t
		}
		if t < eps {
			break
		}
	}
	return y
}
