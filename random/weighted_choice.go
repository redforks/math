package random

import "math/rand"

// WeightedChoice select a random item according weight, returns selected item index.
func WeightedChoice(weights []float64) int {
	r := rand.Float64()
	upto := 0.0
	for k, w := range weights {
		upto += w
		if upto >= r {
			return k
		}
	}
	panic("WeightedChoice() failed") // should never happen
}
