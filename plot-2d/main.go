package main

import (
	"fmt"
	"math"

	"main/graph"
)

func main() {
	// define the function to draw
	f := func(x float64) float64 {
		return math.Sin(x) / (x * (x + 1))
	}

	// take the inputs
	var a, b float64 // domain interval [a,b]
	fmt.Printf("\nDrawing the graph of sin(x)/(x(x+1)) over the interval [a,b]\n")
	fmt.Printf("Input the lower bound(a): ")
	_, err := fmt.Scanf("%f\n", &a)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Input the upper bound(b): ")
	_, err = fmt.Scanf("%f\n", &b)
	if err != nil {
		panic(err)
	}
	if a >= b {
		fmt.Printf("The lower bound(%.3f) must be less than the upper bound(%.3f)\n\n", a, b)
		return
	}

	g := graph.New(f, [2]float64{a, b})
	g.Draw("graph.png")
}
