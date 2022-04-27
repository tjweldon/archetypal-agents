package utils

import (
	"math"
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandFloat(a, b float64) float64 {
	lower := math.Min(a, b)
	return random.Float64()*math.Abs(a-b) + lower
}

func RandInt(a, b int) int {
	min, max := a, b
	if b < a {
		min, max = b, a
	}
	diff := max - min + 1
	return random.Intn(diff) + min
}
