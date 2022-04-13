package world

import (
	"fmt"
	"math"
	"testing"
)

const MaxDegrees = 10.0

func pathLength(circumference, circuits float64) float64 {
	return circumference * circuits
}

func TestDistanceBetweenEndsIsSmall(t *testing.T) {
	circle := Circles(1.0)
	oneDegree := 1.0 / MaxDegrees
	minusOneDegree := (MaxDegrees - 1.0) / MaxDegrees
	minDelta := 10e-9

	geodesic := circle.Metric(oneDegree, minusOneDegree)

	if math.Abs(geodesic-2*oneDegree) > minDelta {
		t.Fatalf("The distance between +1 degree and -1 degree was off by %e in %e", geodesic-2*oneDegree, geodesic)
	}
}

func TestMaximumDistanceIsHalfCircumference(t *testing.T) {
	cirumference := math.SqrtE
	circle := Circles(cirumference)

	dists := func(n int) float64 {
		return circle.Metric(pathLength(cirumference, float64(n)/MaxDegrees), 0)
	}

	formattedContext := func(angle, offset int) string {
		return fmt.Sprintf(
			"{Angle: %d, Distance: %.4f}",
			angle-offset%MaxDegrees, dists(angle-offset),
		)
	}

	for i := range [MaxDegrees * 2]any{} {
		if dists(i) > cirumference/2 {
			t.Fatalf("The distance became greater than half the circumference at %s. Preceding distances: %s, %s",
				formattedContext(i, 0), formattedContext(i, 1), formattedContext(i, 2),
			)
		}
	}
}

func BenchmarkCirclesDistanceCalculation(b *testing.B) {
	cirumference := math.SqrtE
	circle := Circles(cirumference)

	dists := func(n int) float64 {
		return circle.Metric(pathLength(cirumference, float64(n)/MaxDegrees), 0)
	}

	for i := 0; i < b.N; i++ {
		_ = dists(i)
	}
}
