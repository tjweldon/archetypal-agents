package agents

import (
	"math"
	"tjweldon/archetypal-agents/domain/world"
	"tjweldon/archetypal-agents/utils"
)

var (
	width           = 800.0
	height          = 400.0
	neigborhoodSize = 50.0
	maxSpeed        = 10.0
)

// Agent represents an atomic interacting component of the simulation
type Agent struct {
	Index                            int
	Position, Velocity, Acceleration *world.Vector
	neighbourhood                    *Neighbourhood
}

// NewAgent initialises an agent that starts at the top left
//goland:noinspection NonAsciiCharacters
func NewAgent(space, tangentSpace *world.MetricSpace2D, neighbourhood *Neighbourhood, index int, randomise bool) (agent *Agent) {
	neighbourhood.owner = agent
	agent.neighbourhood = neighbourhood

	agent.Index = index

	agent.Acceleration = tangentSpace.ZeroVector()
	agent.Position, agent.Velocity = space.ZeroVector(), tangentSpace.ZeroVector()

	if randomise {
		agent.Position.AccumXY(utils.RandFloat(0, width), utils.RandFloat(0, height))

		// Use plane polar for initial randomisation since that's easier when a max magnitude is imposed
		r, θ := utils.RandFloat(0, maxSpeed), utils.RandFloat(0, 2*math.Pi)

		agent.Velocity.AccumXY(r*math.Cos(θ), r*math.Sin(θ))
	}
	return agent
}

// Neighbourhood provides an interface for an Agent (owner) to interact with its
// sphere of influences
type Neighbourhood struct {
	owner         *Agent
	state         *State
	displacements []*world.Vector
}

func (n *Neighbourhood) Init(s *State) *Neighbourhood {
	n.state = s
	n.displacements = make([]*world.Vector, 0, s.Population())
	return n
}

func (n *Neighbourhood) Population() int {
	return len(n.displacements)
}

func (n *Neighbourhood) CalculateDisplacements() {
	for index, displacement := range n.state.Displacements()[n.owner.Index] {
		if index == n.owner.Index || displacement.Mag() >= neigborhoodSize {
			continue
		}
		n.displacements = append(n.displacements, displacement)
	}
}

func (n *Neighbourhood) Displacements() []*world.Vector {
	return n.displacements
}

// Archetype is the information internal to each Agent that completely characterises its
// behavioral characteristics.
type Archetype struct {
	Influences    Forces
	Sensitivities Charges
}

// Define initialises and returns an Archetype. An Archetype is characterised by its
// Influence in terms of Forces and its Sensitivities in terms of Charges
func Define(sensitivities, influence Charges) Archetype {
	return Archetype{Sensitivities: sensitivities, Influences: (Forces{}).Init(influence)}
}

// DefineReciprocal can be used to express archetypes that are influential and sensitive
// in exactly equal measure
func DefineReciprocal(motives Charges) Archetype {
	return Define(motives, motives)
}

// Innocent is the base archetype. It is an implementation of the classic boids model and
// symmetrical in that its sensitivities and influence are exactly reciprocal, until we
// decide otherwise.
func Innocent() Archetype {
	return DefineReciprocal(Charges{Separation: 1, Cohesion: 1, Alignment: 1})
}
