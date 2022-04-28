package agents

import (
	"math"
	"tjweldon/archetypal-agents/domain/world"
	"tjweldon/archetypal-agents/utils"
)

var (
	width    = 800.0
	height   = 400.0
	maxSpeed = 10.0
)

// Agent represents an atomic interacting component of the simulation
type Agent struct {
	Position, Velocity, Acceleration *world.Vector
}

// NewAgent initialises an agent that starts at the top left
func NewAgent(positions, tangentSpace *world.MetricSpace2D, randomise bool) *Agent {
	var position, velocity, acceleration *world.Vector
	acceleration = tangentSpace.ZeroVector()
	if !randomise {
		position, velocity = positions.ZeroVector(), tangentSpace.ZeroVector()
	} else {
		position = positions.NewVector(utils.RandFloat(0, width), utils.RandFloat(0, height))

		// Use plane polar for initial randomisation since that's easier when a max magnitude is imposed
		r, theta := utils.RandFloat(0, maxSpeed), utils.RandFloat(0, 2*math.Pi)
		velocity = tangentSpace.NewVector(r*math.Cos(theta), r*math.Sin(theta))
	}
	return &Agent{position, velocity, acceleration}
}

// Archetype is the information internal to each Agent that completely characterises its
// behavioral characteristics.
type Archetype struct {
	Influences    Forces
	Sensitivities Charges
}

// Define initialises and returns an Archetype. An Archetype is characterised by its Influence in terms of Forces
// and its Sensitivities in terms of Charges
func Define(sensitivities, influence Charges) Archetype {
	return Archetype{Sensitivities: sensitivities, Influences: (Forces{}).Init(influence)}
}

// DefineReciprocal can be used to express archetypes that are influential and sensitive in exactly equal measure
func DefineReciprocal(motives Charges) Archetype {
	return Define(motives, motives)
}

// Innocent is the base archetype. It is an implementation of the classic boids model and
// symmetrical in that its sensitivities and influence are exactly reciprocal, until we decide otherwise.
var Innocent = DefineReciprocal(Charges{CollisionAvoidance: 1, Cohesion: 1, Alignment: 1})
