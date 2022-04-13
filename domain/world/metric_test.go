package world

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

const AcceptibleError = 1.0e-9

func randFloat(a, b float64) float64 {
	lower := math.Min(a, b)
	return random.Float64()*math.Abs(a-b) + lower
}

func randInt(a, b int) int {
	min, max := a, b
	if b < a {
		min, max = b, a
	}
	diff := max - min + 1
	return random.Intn(diff) + min
}

type checker func(m Metric, t *testing.T) float64

type metricDefinitionChecks struct {
	Symmetry, Positivity, Minimum, TriangleInequality checker
}

var metricsMustSatisfy = metricDefinitionChecks{
	Symmetry: func(m Metric, t *testing.T) float64 {
		a, b := randFloat(-1.0, 1.0), randFloat(-1.0, 1.0)
		return m(a, b) - m(b, a)
	},
	Positivity: func(m Metric, t *testing.T) float64 {
		a, b := randFloat(-1.0, 1.0), randFloat(-1.0, 1.0)
		return m(a, b)
	},
	Minimum: func(m Metric, t *testing.T) float64 {
		a := randFloat(-1.0, 1.0)
		return m(a, a)
	},
	TriangleInequality: func(m Metric, t *testing.T) float64 {
		f := func() float64 { return randFloat(-1.0, 1.0) }
		a, b, c := f(), f(), f()
		a2b, b2c, a2c := m(a, b), m(b, c), m(a, c)

		return a2b + b2c - a2c
	},
}

func TestLine(t *testing.T) {
	assert.New(t)
	m := Line().Metric

	for i := 0; i < 100; i++ {
		assert.GreaterOrEqual(t, metricsMustSatisfy.Positivity(m, t), 0.0, "Positivity")
		assert.Equal(t, 0.0, metricsMustSatisfy.Symmetry(m, t), "Symmetry")
		assert.Equal(t, 0.0, metricsMustSatisfy.Minimum(m, t), "Minimum")
		assert.GreaterOrEqual(t, metricsMustSatisfy.TriangleInequality(m, t), -AcceptibleError, "TriangleInequality")
	}

}

func TestCircleSmall(t *testing.T) {
	assert.New(t)
	m := Circles(1.0).Metric

	for i := 0; i < 100; i++ {
		assert.GreaterOrEqual(t, metricsMustSatisfy.Positivity(m, t), 0.0, "Positivity")
		assert.Equal(t, 0.0, metricsMustSatisfy.Symmetry(m, t), "Symmetry")
		assert.Equal(t, 0.0, metricsMustSatisfy.Minimum(m, t), "Minimum")
		assert.GreaterOrEqual(t, metricsMustSatisfy.TriangleInequality(m, t), -AcceptibleError, "TriangleInequality")
	}
}

func TestCircleBig(t *testing.T) {
	assert.New(t)
	m := Circles(100.0).Metric

	for i := 0; i < 100; i++ {
		assert.GreaterOrEqual(t, metricsMustSatisfy.Positivity(m, t), 0.0, "Positivity")
		assert.Equal(t, 0.0, metricsMustSatisfy.Symmetry(m, t), "Symmetry")
		assert.Equal(t, 0.0, metricsMustSatisfy.Minimum(m, t), "Minimum")
		assert.GreaterOrEqual(t, metricsMustSatisfy.TriangleInequality(m, t), -AcceptibleError, "TriangleInequality")
	}
}
