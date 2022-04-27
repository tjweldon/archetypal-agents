package agents

import (
	"fmt"
	"math"
	"time"
	"tjweldon/archetypal-agents/domain/world"
	"tjweldon/archetypal-agents/utils"
)

var (
	width    float64 = 800
	height   float64 = 400
	maxSpeed float64 = 10.0
)

// Agent represents an atomic interacting component of the simulation
type Agent struct {
	Position, Velocity *world.Vector
}

// NewAgent initialises an agent that starts at the top left
func NewAgent(positions, velocities *world.MetricSpace2D, randomise bool) *Agent {
	var position, velocity *world.Vector
	if !randomise {
		position = positions.ZeroVector()
		velocity = velocities.ZeroVector()
	} else {
		position = positions.NewVector(utils.RandFloat(0, width), utils.RandFloat(0, height))

		// Use plane polar for initial randomisation since that's easier when a max magnitude is imposed
		r, theta := utils.RandFloat(0, maxSpeed), utils.RandFloat(0, 2*math.Pi)
		velocity = velocities.NewVector(r*math.Cos(theta), r*math.Sin(theta))
	}
	return &Agent{position, velocity}
}

// State represents a static (and informationally complete) snapshot of the simulation at a given time
type State struct {
	CoordinateSystem *world.MetricSpace2D
	Agents           []*Agent
}

// NewState initialises a new State struct with the position and velocity vector spaces expressed as
// a pair of world.MetricSpace2D
func NewState(positions, velocities *world.MetricSpace2D) *State {
	agents := make([]*Agent, 100)
	for index := range agents {
		agents[index] = NewAgent(positions, velocities, true)
	}
	return &State{Agents: agents, CoordinateSystem: positions}
}

// Distances calculates an array where the value at distances[i][j] is the distance from State.Agents[i] to State.Agents[j]. This has the property that
// distances[i][j] == distances[j][i] and distances[i][i] == 0
func (s State) Distances() (distances [][]float64) {
	distances = make([][]float64, s.Population())
	for index := range distances {
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

// InitialiseScenario sets up the simulation scenario.
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

// Coords are a part of the socket API, probably shouldn't be defined here
type Coords struct {
	X LPFloat `json:"x"`
	Y LPFloat `json:"y"`
}

// Frame (see comment on Coords)
type Frame []Coords

// GetFrameAt Retrieves the frame data at time t (in seconds)
func GetFrameAt(t float64) Frame {
	return Frame{}
}

// GetNextFrame is a generator function for Frame instances
func (s Scenario) GetNextFrame() (frame Frame) {
	return Frame{}
}
