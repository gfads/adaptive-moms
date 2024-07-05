package controllers

import "adaptive-moms/parameters"

type ASTAR struct {
	Min            float64
	Max            float64
	HysteresisBand float64
	PC             float64
	PreviousRate   float64
	PreviousOut    float64
	Out            float64
}

func (c *ASTAR) Initialise(p parameters.AllParameters) {
	c.Min = p.Min
	c.Max = p.Max
	c.HysteresisBand = p.HysteresisBand
	c.PC = p.PC
}

func (c *ASTAR) Update(p ...float64) float64 {
	u := 0.0

	r := p[0]
	y := p[1] // measured arrival rate

	if y < (r - c.HysteresisBand) { // The system is bellow the goal
		if y > c.PreviousRate {
			u = c.PreviousOut + 1
		} else {
			u = c.PreviousOut * 2
		}
	} else if y > (r + c.HysteresisBand) { // The system is above the goal
		if y < c.PreviousRate {
			u = c.PreviousOut - 1
		} else {
			u = c.PreviousOut / 2
		}
	} else { // The system is at Optimum state, no action required
		u = c.PreviousOut
	}

	// final check of rnew
	if u < c.Min {
		u = c.Min
	}
	if u > c.Max {
		u = c.Max
	}

	c.PreviousOut = u
	c.PreviousRate = y

	return u
}
