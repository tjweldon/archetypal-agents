package world

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVector_Accumulate_EuclideanPlane(t *testing.T) {
	assert.New(t)
	// arrange
	plane := NewEuclideanPlane()
	basis := [2]*Vector{plane.NewVector(1, 0), plane.NewVector(0, 1)}
	maxPrecison := 1.0e-9

	x := randInt(1, 9)
	y := 10 - x

	terms := make([]*Vector, 10)
	for i := 1; i <= 10; i++ {
		if i <= x {
			terms[i-1] = basis[0]
		} else {
			terms[i-1] = basis[1]
		}
	}

	// action
	result := plane.ZeroVector().Accumulate(terms...)

	// assert
	assert.LessOrEqual(t, result.X-float64(x), maxPrecison)
	assert.LessOrEqual(t, result.Y-float64(y), maxPrecison)
}
