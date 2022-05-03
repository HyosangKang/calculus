package main

import (
	"main/spsim"
)

// This program animates the solution of the spring motion.
func main() {
	s := spsim.NewSpSim()
	s.Animate()
}
