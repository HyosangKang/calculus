package main

import (
	"calculus/spsim/unforced"
)

// This program animates the solution of the spring motion.
func main() {
	s := unforced.NewSpSim()
	s.Animate("unforced.gif")
}
