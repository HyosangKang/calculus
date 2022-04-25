package main

import "main/esim"

// The package simulates electric currents of
// the following electric circuits.
//  --L--- --R2--
// |      |      |+
// |      C      V
// |      |      |-
//  --R1-- ------
//   -i1->  -i2->

func main() {
	e := esim.NewEsim()
	e.Graph()
}
