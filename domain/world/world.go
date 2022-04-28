package world

import "math"

// Sum is the generalisation of + on real numbers. An implementation should have the following
// properties:
//  - Commutativity: 		Sum(a, b) == Sum(b, a)
//  - Associativity:		Sum(Sum(a, b), c) == Sum(a, Sum(b, c)) === Sum(a, b, c)
//  - A Zero value: 		Sum(0, a) == a
//  - Existence of Inverse: Sum(-a, a) == 0
type Sum func(scalars ...float64) float64

// Metric is the generalisation of the absolute distance between two points in a connected 1d space.
// Any implementation should have the following properties:
//  - Symmetry: 			Metric(a, b) == Metric(b, a)
//  - Positivity: 			Metric(a, b) >= 0
//  - Minimum: 				Metric(a, a) == 0
//  - Triangle Inequality: 	Metric(a, b) + Metric(b, c) >= Metric(a, c)
type Metric func(scalarA, scalarB float64) float64

// Invert is the generalisation of changing the sign of a value, i.e. finding the inverse under addition.
// An implementation should have the following properties:
//  - Inversion under Sum: 			Sum(Invert(a), a) == 0
//  - Fixed Point is 0: 			Invert(0) == 0
//  - Invert^2 === Identity: 		Invert(Invert(a)) == a
type Invert func(scalar float64) float64

// MetricSpace1D is a representation of the way distance is calculated in a given coordinate.
// It contains no state, it is an attribute of the environment in which the simulation takes place.
type MetricSpace1D struct {
	Sum    Sum
	Metric Metric
	Invert Invert
}

// RealLine returns a MetricSpace1D that behaves like the usual real numbers unbounded above and below
func RealLine() MetricSpace1D {
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

// Circles returns a MetricSpace1D that behaves like a circle with the circumference provided.
// Think of the coordinate as an angle parameter, but not normalised by radius,
// i.e. max(angle) == circumference.
func Circles(circumference float64) MetricSpace1D {
	line := RealLine()
	return MetricSpace1D{
		Sum: func(scalars ...float64) float64 {
			return math.Remainder(line.Sum(scalars...), circumference)
		},
		Metric: func(scalarA, scalarB float64) float64 {
			result := math.Remainder(line.Metric(scalarA, scalarB), circumference)
			if result > circumference/2 {
				result = circumference - result
			}
			return math.Abs(result)
		},
		Invert: func(scalar float64) float64 {
			return circumference - math.Remainder(scalar, circumference)
		},
	}
}

// VectorSpace as an interface roughly corresponds to the mathematical notion in the sense
// that the Vector objects it returns are well behaved
type VectorSpace interface {
	ZeroVector() *Vector
	NewVector(x, y float64) *Vector
}

// MetricSpace2D is a cartesian product of two MetricSpace1D
type MetricSpace2D struct {
	XCoord, YCoord MetricSpace1D
}

// GeodesicDiff returns a tuple of geodesic distances between two points (x1, y1) and (x2, y2).
// This defaults to the shortest straight path between the points. This is useful in toroidal topologies.
// The result retains the sign of the difference as well as the magnitude
func (m MetricSpace2D) GeodesicDiff(x1, x2, y1, y2 float64) (deltaX, deltaY float64) {
	deltaX = m.XCoord.Sum(m.XCoord.Invert(x1), x2)
	deltaY = m.YCoord.Sum(m.YCoord.Invert(y1), y2)
	return deltaX, deltaY
}

// Metric is the function that defines the distance between two points represented by the Vector
// pair (v1, v2)
func (m *MetricSpace2D) Metric(v1, v2 *Vector) float64 {
	deltaX, deltaY := m.GeodesicDiff(v1.X, v2.X, v1.Y, v2.Y)
	return math.Sqrt(math.Pow(deltaX, 2) + math.Pow(deltaY, 2))
}

// NewEuclideanPlane initialises a MetricSpace2D that represents euclidean geometry and non-periodic
// boundary conditions i.e. an infinite 2D plane
func NewEuclideanPlane() *MetricSpace2D {
	line := RealLine()
	return &MetricSpace2D{
		XCoord: line,
		YCoord: line,
	}
}

// NewEuclideanToroid initialises a MetricSpace2D that represents geometry and periodic boundary
// conditions on both axes i.e. a torus. This corresponds to agents leaving the screen on one
// boundary reappearing at the opposite boundary.
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

// Vector represents a point in a space and can be summed in a way that respects linearity.
// This is a mutable implementation because I don't want to get a memory out.
type Vector struct {
	metricSpace *MetricSpace2D
	X, Y        float64
}

// Accumulate is a mutating variadic sum. This is the mutable equivalent of Sum
func (v *Vector) Accumulate(vectors ...*Vector) *Vector {
	xScalars, yScalars := make([]float64, len(vectors)), make([]float64, len(vectors))
	for index, summand := range vectors {
		xScalars[index], yScalars[index] = summand.X, summand.Y
	}
	v.X = v.metricSpace.XCoord.Sum(xScalars...)
	v.Y = v.metricSpace.YCoord.Sum(yScalars...)

	return v
}

// Times is an **immutable** scalar multiple of the vector instance where the semantics are:
// u = v.Times(x) => u != v
func (v *Vector) Times(scalar float64) Vector {
	x := v.metricSpace.XCoord.Sum(scalar * v.X)
	y := v.metricSpace.YCoord.Sum(scalar * v.Y)
	return Vector{v.metricSpace, x, y}
}

// Scale is a **mutating** scalar multiple of the vector instance with the following semantics:
// u = v.Scale(x); u == v
func (v *Vector) Scale(scalar float64) *Vector {
	v.X = v.metricSpace.XCoord.Sum(scalar * v.X)
	v.Y = v.metricSpace.YCoord.Sum(scalar * v.Y)
	return v
}

// Minus is an **immutable** operation where the semantics are:
// u = v.Minus(w) => u != v
func (v *Vector) Minus(vector *Vector) Vector {
	deltaX, deltaY := v.metricSpace.GeodesicDiff(v.X, vector.X, v.Y, vector.Y)
	return Vector{X: deltaX, Y: deltaY}
}

// Subtract is a **mutable** operation where the semantics are:
// u = v.Subtract(w) => u == v
func (v *Vector) Subtract(vector *Vector) *Vector {
	v.X, v.Y = v.metricSpace.GeodesicDiff(v.X, vector.X, v.Y, vector.Y)
	return v
}

// Mag returns the length of the vector i.e. the geodesic distance from the head of
// the vector to the origin.
func (v *Vector) Mag() float64 {
	return v.metricSpace.Metric(v, v.metricSpace.ZeroVector())
}

// Dot is the scalar product of two vectors with the following semantics:
// x = u ⋅ v <=> x := u.Dot(v)
//
// Dot should have the following properties:
//  - u.Dot(v) == 0.0 				<=> u perpendicular to v
//  - u.Dot(v) == u.Mag()*v.Mag() 	<=> u parallel to v
//    -	therefore:					==> math.Sqrt(u.Dot(u)) === u.Mag()
func (v *Vector) Dot(w *Vector) float64 {
	norm := v.metricSpace.XCoord.Sum
	return norm(v.X)*norm(w.X) + norm(v.Y)*norm(w.Y)
}
