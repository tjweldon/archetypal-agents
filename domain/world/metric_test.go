package world

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tjweldon/archetypal-agents/utils"
)

const MaxPrecision = 1.0e-9

type checker func(m Metric, t *testing.T) float64

type metricDefinitionChecks struct {
	Symmetry, Positivity, Minimum, TriangleInequality checker
}

var metricsMustSatisfy = metricDefinitionChecks{
	Symmetry: func(m Metric, t *testing.T) float64 {
		a, b := utils.RandFloat(-1.0, 1.0), utils.RandFloat(-1.0, 1.0)
		return m(a, b) - m(b, a)
	},
	Positivity: func(m Metric, t *testing.T) float64 {
		a, b := utils.RandFloat(-1.0, 1.0), utils.RandFloat(-1.0, 1.0)
		return m(a, b)
	},
	Minimum: func(m Metric, t *testing.T) float64 {
		a := utils.RandFloat(-1.0, 1.0)
		return m(a, a)
	},
	TriangleInequality: func(m Metric, t *testing.T) float64 {
		f := func() float64 { return utils.RandFloat(-1.0, 1.0) }
		a, b, c := f(), f(), f()
		a2b, b2c, a2c := m(a, b), m(b, c), m(a, c)

		return a2b + b2c - a2c
	},
}

func TestLine(t *testing.T) {
	assert.New(t)
	m := RealLine().Metric

	for i := 0; i < 100; i++ {
		assert.GreaterOrEqual(t, metricsMustSatisfy.Positivity(m, t), 0.0, "Positivity")
		assert.Equal(t, 0.0, metricsMustSatisfy.Symmetry(m, t), "Symmetry")
		assert.Equal(t, 0.0, metricsMustSatisfy.Minimum(m, t), "Minimum")
		assert.GreaterOrEqual(t, metricsMustSatisfy.TriangleInequality(m, t), -MaxPrecision, "TriangleInequality")
	}

}

func TestCircleSmall(t *testing.T) {
	assert.New(t)
	m := Circles(1.0).Metric

	for i := 0; i < 100; i++ {
		assert.GreaterOrEqual(t, metricsMustSatisfy.Positivity(m, t), 0.0, "Positivity")
		assert.Equal(t, 0.0, metricsMustSatisfy.Symmetry(m, t), "Symmetry")
		assert.Equal(t, 0.0, metricsMustSatisfy.Minimum(m, t), "Minimum")
		assert.GreaterOrEqual(t, metricsMustSatisfy.TriangleInequality(m, t), -MaxPrecision, "TriangleInequality")
	}
}

func TestCircleBig(t *testing.T) {
	assert.New(t)
	m := Circles(100.0).Metric

	for i := 0; i < 100; i++ {
		assert.GreaterOrEqual(t, metricsMustSatisfy.Positivity(m, t), 0.0, "Positivity")
		assert.Equal(t, 0.0, metricsMustSatisfy.Symmetry(m, t), "Symmetry")
		assert.Equal(t, 0.0, metricsMustSatisfy.Minimum(m, t), "Minimum")
		assert.GreaterOrEqual(t, metricsMustSatisfy.TriangleInequality(m, t), -MaxPrecision, "TriangleInequality")
	}
}
