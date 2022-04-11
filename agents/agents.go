package agents

import (
	"fmt"
	"math"
	"time"
)

var (
	width  float64 = 800
	height float64 = 400
)

type Sum func(scalars ...float64) float64

type Diff func(scalarA, scalarB float64) float64

// MetricSpace1D is a representation of the way distance is calculated in a given coordinate.
// It contains no state, it is an attribute of the environment in which the simulation takes place.
type MetricSpace1D struct {
	Sum    Sum
	Metric Diff
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
			result := line.Metric(scalarA, scalarB)
			if result > circumference/2 {
				result = math.Abs(circumference - result)
			}
			return result
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
	deltaX = m.XCoord.Metric(x1, x2)
	deltaY = m.YCoord.Metric(y1, y2)
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

// NewVector is a method used to initialise a Vector in the metric space m
func (m *MetricSpace2D) NewVector() *Vector {
	return &Vector{metricSpace: m}
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

type Agent struct {
	Position, Velocity *Vector
}

func NewAgent(positions, velocities *MetricSpace2D) *Agent {
	return &Agent{positions.NewVector(), velocities.NewVector()}
}

// State represents a static snapshot of the simulation at a given time
type State struct {
	CoordinateSystem *MetricSpace2D
	Agents           []*Agent
}

func NewState(positions, velocities *MetricSpace2D) *State {
	agents := make([]*Agent, 100)
	for index, _ := range agents {
		agents[index] = NewAgent(positions, velocities)
	}
	return &State{Agents: agents, CoordinateSystem: positions}
}

// Distances calculates an array where the value at distances[i][j] is the distance from State.Agents[i] to State.Agents[j]. This has the property that
// distances[i][j] == distances[j][i] and distances[i][i] == 0
func (s State) Distances() (distances [][]float64) {
	distances = make([][]float64, s.Population())
	for index, _ := range distances {
		distances[index] = make([]float64, s.Population())
	}

	for i := 0; i < len(distances); i++ {
		for j := 0; j < i; j++ {
			distances[i][j] = s.CoordinateSystem.Metric(s.Agents[i].Position, s.Agents[j].Position)
			distances[j][i] = distances[i][j]
		}
	}

	return distances
}

// Population returns the number of agents
func (s State) Population() int {
	population := len(s.Agents)
	return population
}

// Scenario is a complete encapsulation of the simulation. It contains the current time, state, topology and simulation timeStep
type Scenario struct {
	Time, DeltaT          time.Duration
	positions, velocities *MetricSpace2D
	state                 *State
}

func InitialiseScenario(timeStep time.Duration) Scenario {
	toroid, euclideanPlane := NewEuclideanToroid(width, height), NewEuclideanPlane()
	state := NewState(toroid, euclideanPlane)
	return Scenario{positions: toroid, velocities: euclideanPlane, state: state, DeltaT: timeStep}
}

// LPFloat serialises as a number represented to a fixed number
// decimal places eg. 1.00
type LPFloat struct {
	Value  float64 // the actual value
	Digits int     // the number of digits used in json
}

// MarshalJSON serialises the LPFloat Type
func (l LPFloat) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.*f", l.Digits, l.Value)
	return []byte(s), nil
}

type Coords struct {
	X LPFloat `json:"x"`
	Y LPFloat `json:"y"`
}

type Frame []Coords

// GetFrameAt Retrieves the frame data at time t (in seconds)
func GetFrameAt(t float64) Frame {
	return Frame{}
}

// GetNextFrame is a generator function for Frame instances
func (s Scenario) GetNextFrame() (frame Frame) {
	return Frame{}
}
