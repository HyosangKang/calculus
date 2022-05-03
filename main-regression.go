package main

import (
	"main/scatter"

	"fmt"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	data := scatter.NewSample()
	line := data.Regression()
	p := plot.New()
	data.Graph(p)
	line.Graph(p)
	p.Save(400, 400, "regression-exact.png")
	fmt.Printf("Exact: %.3f\n", data.S(line))
	fmt.Println(line)

	line = data.GradientDescent()
	p = plot.New()
	data.Graph(p)
	line.Graph(p)
	p.Save(400, 400, "regression-gradient.png")
	fmt.Printf("Gradient descent: %.3f\n", data.S(line))
	fmt.Println(line)
}
