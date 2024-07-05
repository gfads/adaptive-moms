package controllers

import (
	"adaptive-moms/parameters"
	"math"
)

type HPA struct {
	PC  float64
	Out float64
	Max float64
	Min float64
}

func (c *HPA) Initialise(p parameters.AllParameters) {
	c.Max = p.Max
	c.Min = p.Min
	c.PC = p.PC
}

func (c *HPA) Update(p ...float64) float64 {
	u := 0.0

	r := p[0] // goal
	y := p[1] // plant output

	u = math.Round(c.PC * r / y)

	// control law
	if u > c.Max {
		u = c.Max
	} else if u < c.Min {
		u = c.Min
	}

	c.PC = u

	return u
}
