package agents

import (
	"fmt"
	"time"
	"tjweldon/archetypal-agents/domain/world"
)

var (
	width  float64 = 800
	height float64 = 400
)

type Agent struct {
	Position, Velocity *world.Vector
}

func NewAgent(positions, velocities *world.MetricSpace2D) *Agent {
	return &Agent{positions.ZeroVector(), velocities.ZeroVector()}
}

// State represents a static snapshot of the simulation at a given time
type State struct {
	CoordinateSystem *world.MetricSpace2D
	Agents           []*Agent
}

func NewState(positions, velocities *world.MetricSpace2D) *State {
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
	positions, velocities *world.MetricSpace2D
	state                 *State
}

func InitialiseScenario(timeStep time.Duration) Scenario {
	toroid, euclideanPlane := world.NewEuclideanToroid(width, height), world.NewEuclideanPlane()
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
