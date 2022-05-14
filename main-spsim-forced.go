package main

import "calculus/spsim/forced"

// This program animates the solution of the spring motion.
func main() {
	s := forced.NewSpSim()
	s.Animate("forced.gif")
}
