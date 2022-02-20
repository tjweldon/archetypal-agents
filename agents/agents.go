package agents

import (
	"fmt"
	"math"
)

var (
	width  int = 800
	height int = 400
)

type trajectory func(t float64) Coords

func straightPath(vX float64, vY float64) trajectory {
	traj := func(t float64) Coords {
		x := int(t*vX) % width
		y := int(t*vY) % height
		return Coords{
			X: LPFloat{Value: float64(x)}, Y: LPFloat{Value: float64(y)},
		}
	}

	return traj
}

var Line = straightPath(30.0, 45.0)

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
	return Frame{
		Coords{
			X: LPFloat{Value: 100*math.Cos(2*math.Pi*t) + 200},
			Y: LPFloat{Value: 100*math.Sin(2*math.Pi*t) + 200},
		},
		Line(t),
		Line(t + 0.5),
		Line(t - 0.5),
	}
}
