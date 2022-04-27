package world

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tjweldon/archetypal-agents/utils"
)

func getAdjacencyAssertion(t *testing.T, space MetricSpace1D) func(coordA, coordB float64) {
	a := assert.New(t)
	return func(coordA, coordB float64) {
		separation := space.Metric(coordA, coordB)
		a.LessOrEqual(separation, MaxPrecision)
	}
}

func TestLine_Sum(t *testing.T) {
	space := RealLine()
	isCommutative(t, space)
	isAssociative(t, space)
	hasZeroElement(t, space)
	inverseExists(t, space)
}

func TestCircles_Sum(t *testing.T) {
	spaces := [3]MetricSpace1D{
		Circles(1), Circles(10), Circles(100),
	}
	for _, circle := range spaces {
		isCommutative(t, circle)
		isAssociative(t, circle)
		hasZeroElement(t, circle)
		inverseExists(t, circle)
	}
}

func isCommutative(t *testing.T, space MetricSpace1D) {
	assertAdjacent := getAdjacencyAssertion(t, space)
	getSummands := func() (alpha float64, beta float64) {
		alpha, beta = utils.RandFloat(-10, 10), utils.RandFloat(-10, 10)
		return alpha, beta
	}

	var a, b float64
	for range [100]any{} {
		a, b = getSummands()

		// Assert Sum(a, b) == Sum(b, a)
		assertAdjacent(space.Sum(a, b), space.Sum(b, a))
	}
}

func isAssociative(t *testing.T, space MetricSpace1D) {
	assertAdjacent := getAdjacencyAssertion(t, space)
	getSummands := func() (summands []float64) {
		summands = []float64{
			utils.RandFloat(-10, 10), utils.RandFloat(-10, 10), utils.RandFloat(-10, 10),
		}
		return summands
	}

	var s, sums []float64
	for range [100]any{} {
		s = getSummands()
		sums = []float64{
			space.Sum(s[0], space.Sum(s[1], s[2])), // Right fold
			space.Sum(space.Sum(s[0], s[1]), s[2]), // Left fold
			space.Sum(s...),                        // Variadic
		}

		// Assert Sum(a, Sum(b, c)) == Sum(Sum(a, b), c)
		assertAdjacent(sums[0], sums[1])

		// Assert Sum(a, Sum(b, c)) == Sum(Sum(a, b), c) === Sum(a, b, c)
		assertAdjacent(sums[0], sums[2])
		assertAdjacent(sums[1], sums[2])
	}
}

func hasZeroElement(t *testing.T, space MetricSpace1D) {
	assertAdjacent := getAdjacencyAssertion(t, space)
	var zero, a float64

	for range [100]any{} {
		a = utils.RandFloat(-10, 10)

		assertAdjacent(space.Sum(zero, a), a)
	}
}

func inverseExists(t *testing.T, space MetricSpace1D) {
	assertAdjacent := getAdjacencyAssertion(t, space)

	var a float64

	for range [100]any{} {
		a = utils.RandFloat(-10, 10)

		assertAdjacent(space.Sum(space.Invert(a), a), 0.0)
	}
}
