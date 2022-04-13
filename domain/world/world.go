package world

import "math"

// Sum is the generalisation of + on real numbers. An implementation should have the following characteristics:
//  - Sum(a, b) == Sum(b, a)
//  - Sum(Sum(a, b), c) == Sum(a, Sum(b, c))
//  - Sum(0, a) == a
//  - Sum(-a, a) == 0
type Sum func(scalars ...float64) float64

// Metric is the generalisation of the absolute distance between two points in a connected 1d space. Any implementation
// should have the following properties:
//  - Metric(a, b) == Metric(b, a)
//  - Metric(a, b) >= 0
//  - Metric(a, a) == 0
//  - Metric(a, b) + Metric(b, c) >= Metric(a, c)
type Metric func(scalarA, scalarB float64) float64

// Invert is the generalisation of changing the sign of a value, i.e. finding the inverse under addition.
// An implementation should have the following properties:
//  - Sum(Invert(a), a) == 0
//  - Invert(0) == 0
//  - Invert(Invert(a)) == a
type Invert func(scalar float64) float64

// MetricSpace1D is a representation of the way distance is calculated in a given coordinate.
// It contains no state, it is an attribute of the environment in which the simulation takes place.
type MetricSpace1D struct {
	Sum    Sum
	Metric Metric
	Invert Invert
}

// Line returns a MetricSpace1D that behaves like the usual real numbers unbounded above and below
func Line() MetricSpace1D {
	return MetricSpace1D{
		Sum: func(scalars ...float64) float64 {
			sum := 0.0
			for _, summand := range scalars {
				sum += summand
			}
			return sum
		},
		Metric: func(scalarA, scalarB float64) float64 {
			return math.Abs(scalarA - scalarB)
		},
		Invert: func(scalar float64) float64 {
			return -scalar
		},
	}
}

// Circles returns a MetricSpace1D that behaves like a circle with the circumference provided (think of an angle parameter)
func Circles(circumference float64) MetricSpace1D {
	line := Line()
	return MetricSpace1D{
		Sum: func(scalars ...float64) float64 {
			return math.Remainder(line.Sum(scalars...), circumference)
		},
		Metric: func(scalarA, scalarB float64) float64 {
			result := math.Remainder(line.Metric(scalarA, scalarB), circumference)
			if result > circumference/2 {
				result = circumference/2 - result
			}
			return math.Abs(result)
		},
		Invert: func(scalar float64) float64 {
			return circumference - math.Remainder(scalar, circumference)
		},
	}
}

// MetricSpace2D is a cartesian product of two MetricSpace1D
type MetricSpace2D struct {
	XCoord, YCoord MetricSpace1D
}

// GeodesicDiff returns a tuple of geodesic distances between two points (x1, y1) and (x2, y2).
// This defaults to the shortest straight path between the points. This is useful in toroidal topologies.
func (m MetricSpace2D) GeodesicDiff(x1, x2, y1, y2 float64) (deltaX, deltaY float64) {
	deltaX = m.XCoord.Sum(m.XCoord.Invert(x1), x2)
	deltaY = m.YCoord.Metric(m.YCoord.Invert(y1), y2)
	return deltaX, deltaY
}

// Metric is the function that defines the distance between two points represented by the Vector pair (v1, v2)
func (m *MetricSpace2D) Metric(v1, v2 *Vector) float64 {
	deltaX, deltaY := m.GeodesicDiff(v1.X, v2.X, v1.Y, v2.Y)
	return math.Sqrt(math.Pow(deltaX, 2) + math.Pow(deltaY, 2))
}

// NewEuclideanPlane initialises a MetricSpace2D that represents euclidean geometry and non-periodic boundary conditions i.e. an
// infinite 2D plane
func NewEuclideanPlane() *MetricSpace2D {
	line := Line()
	return &MetricSpace2D{
		XCoord: line,
		YCoord: line,
	}
}

// NewEuclideanToroid initialises a MetricSpace2D that represents geometry and periodic boundary conditions on both axes i.e. a torus.
// This corresponds to agents leaving the screen on one boundary reappearing at the opposite boundary.
func NewEuclideanToroid(w, h float64) *MetricSpace2D {
	return &MetricSpace2D{
		XCoord: Circles(w),
		YCoord: Circles(h),
	}
}

// ZeroVector is a method used to initialise a Vector in the metric space with zero magnitude
func (m *MetricSpace2D) ZeroVector() *Vector {
	return &Vector{metricSpace: m}
}

// NewVector is a method used to initialise a Vector
func (m *MetricSpace2D) NewVector(x, y float64) *Vector {
	return &Vector{metricSpace: m, X: x, Y: y}
}

// Vector represents a point in a space and can be summed in a way that respects linearity. This is a mutable implementation because
// I don't want to get a memory out.
type Vector struct {
	metricSpace *MetricSpace2D
	X, Y        float64
}

// Accumulate is a mutating variadic sum
func (v *Vector) Accumulate(vectors ...*Vector) *Vector {
	xScalars, yScalars := make([]float64, len(vectors)), make([]float64, len(vectors))
	for index, summand := range vectors {
		xScalars[index], yScalars[index] = summand.X, summand.Y
	}
	v.X = v.metricSpace.XCoord.Sum(xScalars...)
	v.Y = v.metricSpace.YCoord.Sum(yScalars...)

	return v
}

func (v *Vector) Minus(vector *Vector) Vector {
	deltaX, deltaY := v.metricSpace.GeodesicDiff(v.X, vector.X, v.Y, vector.Y)
	return Vector{X: deltaX, Y: deltaY}
}
