package world

import "testing"

func TestLine_Invert(t *testing.T) {
	space := RealLine()
	zeroIsOwnInverse(t, space)
	inverseTwiceIsIdentity(t, space)
}

func TestCircles_Invert(t *testing.T) {
	spaces := [3]MetricSpace1D{
		Circles(1), Circles(10), Circles(100),
	}
	for _, circle := range spaces {
		zeroIsOwnInverse(t, circle)
		inverseTwiceIsIdentity(t, circle)
	}
}

func inverseTwiceIsIdentity(t *testing.T, space MetricSpace1D) {
	assertAdjacent := getAdjacencyAssertion(t, space)

	var a float64
	for range [100]any{} {
		a = randFloat(-10, 10)
		assertAdjacent(a, space.Invert(space.Invert(a)))
	}
}

func zeroIsOwnInverse(t *testing.T, space MetricSpace1D) {
	assertAdjacent := getAdjacencyAssertion(t, space)
	zero := 0.0

	assertAdjacent(space.Invert(zero), zero)
}
