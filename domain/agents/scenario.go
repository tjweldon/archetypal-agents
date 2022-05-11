package agents

import (
	"time"
	"tjweldon/archetypal-agents/domain/world"
)

// State represents a static (and informationally complete) snapshot of the simulation at a given time
type State struct {
	CoordinateSystem *world.MetricSpace2D
	Agents           []*Agent

	// displacements & distances have pre-allocated memory for this value to be read from instead of caculating and
	// allocating it every time
	displacements [][]*world.Vector
	distances     [][]float64
}

// NewState initialises a new State struct with the position and velocity vector spaces expressed as
// a pair of world.MetricSpace2D
func NewState(positions, velocities *world.MetricSpace2D) (state *State) {
	state.CoordinateSystem = positions

	population := state.Population()
	state.displacements = make([][]*world.Vector, population)
	state.distances = make([][]float64, population)

	for index := range state.displacements {
		state.displacements[index] = make([]*world.Vector, population)
		state.distances[index] = make([]float64, population)
	}

	state.Agents = make([]*Agent, 100)
	for index := range state.Agents {
		state.Agents[index] = NewAgent(positions, velocities, &Neighbourhood{state: state}, index, true)
	}

	return state
}

// CalculateDistances calculates an array where the value at distances[i][j] is the distance from State.Agents[i] to State.Agents[j].
// This has the property that distances[i][j] == distances[j][i] and distances[i][i] == 0
func (s *State) CalculateDistances() [][]float64 {
	for i := 0; i < s.Population(); i++ {
		for j := 0; j < i; j++ {
			s.distances[i][j] = s.displacements[i][j].Mag()
			s.distances[j][i] = s.distances[i][j]
		}
	}

	return s.distances
}

// Distances returns a 2D array of size population^2 of the distance from agent[i] to agent[j]
func (s *State) Distances() [][]float64 {
	return s.distances
}

// Population returns the number of agents
func (s *State) Population() int {
	population := len(s.Agents)
	return population
}

func (s State) CalculateDisplacements() [][]*world.Vector {
	for i := 0; i < len(s.displacements); i++ {
		for j := 0; j < i; j++ {
			s.displacements[i][j].Scale(0). // displacement[i, j] = 0
				Accumulate(s.Agents[i].Position). //      + pos[i]
				Subtract(s.Agents[j].Position) //      - pos[j]

			s.displacements[j][i].Scale(0). // displacement[j, i] = 0
				Subtract(s.displacements[j][i]) //      - displacement[i, j]
		}
	}

	return s.displacements
}

func (s State) Displacements() [][]*world.Vector {
	return s.displacements
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
	s.CacheWarmup()
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

func (s Scenario) CacheWarmup() {
	s.state.CalculateDisplacements()
	s.state.CalculateDistances()
	for _, agent := range s.state.Agents {
		agent.neighbourhood.CalculateDisplacements()
	}
}
