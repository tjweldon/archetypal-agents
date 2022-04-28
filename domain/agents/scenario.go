package agents

import (
	"time"
	"tjweldon/archetypal-agents/domain/world"
)

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

// Distances calculates an array where the value at distances[i][j] is the distance from State.Agents[i] to State.Agents[j].
// This has the property that distances[i][j] == distances[j][i] and distances[i][i] == 0
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

// Scenario is a complete encapsulation of the simulation. It contains the current time, state, topology
// and simulation timeStep
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

func (s *Scenario) Evolve() {
	// Calculate forces
	s.CalculateForces()
	// Resolve forces
	s.IntegrateForces(s.DeltaT)
	// Resolve constraints
	s.ApplyConstraints()
	// progress simulation by Δt
	s.IntegrateVelocities(s.DeltaT)
	s.Time += s.DeltaT
}

// CalculateForces updates the force vector associated to each agent based on the current state
func (s *Scenario) CalculateForces() {

}

// IntegrateForces updates the velocities of the agents
func (s *Scenario) IntegrateForces(deltaT time.Duration) {

}

// ApplyConstraints will apply any simulation constraints to the velocities etc. before the system is evolved forward one timestep.
func (s *Scenario) ApplyConstraints() {

}

// IntegrateVelocities updates the position of each agent based on their respective velocity and the size of Δt
func (s Scenario) IntegrateVelocities(deltaT time.Duration) {

}
