package agents

import "tjweldon/archetypal-agents/domain/world"

// Field is the label for one of the forces present in the "physics" of the simulation.
// Forces and Charges are labelled by Field
type Field int

const (
	Separation Field = iota
	Cohesion
	Alignment
	fieldCount
	// fieldCount has to stay at the bottom. Is a convenience that allows
	// Iteration over the enum
)

// BasisCharges Returns a new Charges with the charge q with respect to the Field, and 0
// with respect to any other Field.
func (f Field) BasisCharges(q Charge) Charges {
	if f == fieldCount {
		return nil
	}
	charges := (Charges{}).Init()
	charges[f] = q
	return charges
}

type Charge float64

// Charges is an abstraction of the degree to which each force is felt
// Similarly it can also be used to represent the degree to which each force is exerted
type Charges map[Field]Charge

// Transform is a function signature that allows matrix-like transformations of Charges
type Transform func(field Field, f Charge) Charge

// Scalar is basically a special case of a Transform that applies the same scalar transform
// to each charge
type Scalar func(m Charge) Charge

// Init initialises Charges at 0 for every Field
func (c Charges) Init() Charges {
	for field := range [fieldCount]any{} {
		c[Field(field)] = 0
	}

	return c
}

// Apply applies the Transform passed to the receiving Charges
func (c Charges) Apply(mapping Transform) Charges {
	for field, charge := range c {
		c[field] = mapping(field, charge)
	}

	return c
}

// Add as a bare function is the transform that adds c to whichever Charges Apply it:
//	a, b := Alignment.BasisCharges(1), Cohesion.BasisCharges(1)
// 	a: Charges{Separation: 0, Alignment: 1, Cohesion: 0}
// 	b: Charges{Separation: 0, Alignment: 0, Cohesion: 1}
//
// 	c := b.Apply(a.Add())
// 	c: Charges{Separation: 0, Alignment: 1, Cohesion: 1}
//
func (c Charges) Add() Transform {
	return func(field Field, f Charge) Charge {
		return f + c[field]
	}
}

// Action is the signature of a function that represents the way a Force acts to contribute to
// the acceleration of a body
type Action func(acc *world.Vector, n *Neighbourhood, source Charge)

// actionMap is the type alias of a registry of Action functions indexed by the relevant field
type actionMap map[Field]Action

// noAction is the zero value of the Action function type
func noAction(_ *world.Vector, _ *Neighbourhood, _ Charge) {}

// Actions are the actual mapping of a Field to an Action that is used to resolve the motion
var Actions = actionMap{
	Separation: func(acc *world.Vector, n *Neighbourhood, q Charge) {
		n.Displacements()
		acc.Accumulate()
	},
	Alignment: noAction,
	Cohesion:  noAction,
}

// Force is wrapper that provides a way to actually accumulate the action of forces
type Force struct {
	source Charge
	field  Field
	action Action
}

func (f Force) Init(source Charges, field Field) Force {
	f.source = source[field]
	f.field = field
	f.action = Actions[field]
	return f
}

// Apply applies the receiver Force to the acc world.Vector supplied based on the displacement and charge
func (f Force) Apply(acc *world.Vector, neighbourhood *Neighbourhood, sensetivity Charge) {
	f.action(acc, neighbourhood, f.source*sensetivity)
}

// Forces can be thought of as the action of a set of source Charges on another subject body with its own charges
type Forces map[Field]Force

// Init takes the source Charges and returns Forces with each Force applying its action in proportion to the source charge
func (f Forces) Init(source Charges) (forces Forces) {
	var field Field
	for fIndex := range [fieldCount]any{} {
		field = Field(fIndex)
		f[field] = (Force{}).Init(source, field)
	}

	return f
}

// Act is similar to the method of the same name on the Force struct, however it accumulates the action
// of all the forces, scaled by their respective Charges.
func (f Forces) Act(acc *world.Vector, sensitivity Charges, neighbourhood *Neighbourhood) {
	for field, force := range f {
		force.Apply(acc, neighbourhood, sensitivity[field])
	}
}
