package world

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVector_Accumulate_IsLinear_EuclideanPlane(t *testing.T) {
	assert.New(t)
	// Arrange: generate terms for a sum of vectors n*X + (10 - n)*Y where X and Y are
	// basis vectors
	plane := NewEuclideanPlane()
	basis := [2]*Vector{plane.NewVector(1, 0), plane.NewVector(0, 1)}
	x := randInt(1, 9)
	y := 10 - x

	terms := generateTerms(x, basis)

	// Action: Accumulate with the terms generated
	result := plane.ZeroVector().Accumulate(terms...)

	// assert: That result ⋅ X == n and result ⋅ Y == 10 - n
	assert.LessOrEqual(t, plane.XCoord.Metric(result.X, float64(x)), MaxPrecision)
	assert.LessOrEqual(t, plane.YCoord.Metric(result.Y, float64(y)), MaxPrecision)
}

func TestVector_Accumulate_IsLinear_Toroid(t *testing.T) {
	assert.New(t)

	// Arrange: generate terms for a sum of vectors n*X + (10 - n)*Y where X and Y are
	// basis vectors
	plane := NewEuclideanToroid(5, 2)
	basis := [2]*Vector{plane.NewVector(1, 0), plane.NewVector(0, 1)}

	x := randInt(1, 9)
	y := 10 - x

	terms := generateTerms(x, basis)

	// Action: Accumulate with the terms generated
	result := plane.ZeroVector().Accumulate(terms...)

	// assert: That result ⋅ X == n modulo Width and result ⋅ Y == 10 - n modulo Height
	assert.LessOrEqual(t, plane.XCoord.Metric(result.X, float64(x)), MaxPrecision)
	assert.LessOrEqual(t, plane.YCoord.Metric(result.Y, float64(y)), MaxPrecision)
}

func generateTerms(x int, basis [2]*Vector) []*Vector {
	terms := make([]*Vector, 10)
	for i := 1; i <= 10; i++ {
		if i <= x {
			terms[i-1] = basis[0]
		} else {
			terms[i-1] = basis[1]
		}
	}
	return terms
}
