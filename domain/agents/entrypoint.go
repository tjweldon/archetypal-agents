package agents

import "fmt"

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
