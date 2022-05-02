package main

import (
	"fmt"
	"math"

	"main/graph"
)

func main() {
	c := graph.NewCanvas(600, 600)
	// define the function to draw
	y := func(t float64) float64 {
		return math.Cos(t)
	}
	x := func(t float64) float64 {
		return t
	}
	// take the inputs
	// a, b := handleInput()
	g := graph.NewGraph(x, y, -2, 2, 101)
	c.Add(g)

	// add another graph
	y = func(t float64) float64 {
		return math.Sin(t) + 1
	}
	x = func(t float64) float64 {
		return t
	}
	g = graph.NewGraph(x, y, -2, 2, 101)
	c.Add(g)

	c.Draw("graph.png")
}

func handleInput() (float64, float64) {
	var a, b float64
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
		msg := fmt.Sprintf("The lower bound(%.3f) must be less than the upper bound(%.3f)\n\n", a, b)
		panic(msg)
	}
	return a, b
}
