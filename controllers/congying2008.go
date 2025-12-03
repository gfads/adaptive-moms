package controllers

import (
	"adaptive-moms/parameters"
	"math"
)

type Congying struct {
	Alfa     float64
	Alfa0    float64
	Alfa1    float64
	Alfa2    float64
	Beta0    float64
	Beta1    float64
	Beta2    float64
	Delta    float64
	Min      float64
	Max      float64
	PC       float64
	Out      float64
	SumError float64
}

func (c *Congying) Initialise(p parameters.AllParameters) {
	c.Min = p.Min
	c.Max = p.Max
	c.Alfa = p.Alfa
	c.Alfa0 = p.Alfa0
	c.
		c.Delta = p.Delta
	c.PC = p.PC
	c.SumError = 0.0
}

func (c *Congying) Update(p ...float64) float64 {
	u := 0.0

	r := p[0]
	y := p[1] // measured arrival rate
	e := r - y
	c.SumError += e

	//u=β0 fal(∫e,α0,δ) + β1 fal(e,α1,δ) + β2 fal(de,α2,δ)

	u = c.Beta0*c.fal(c.SumError, c.Alfa0, c.Delta) + c.Beta1*c.fal(e, c.Alfa1, c.Delta) + c.Beta2*c.fal(de, c.Alfa2, c.Delta)

	// final check of u
	if u < c.Min {
		u = c.Min
	}
	if u > c.Max {
		u = c.Max
	}

	return u
}

func fal(e float64, alfa float64, delta float64) float64 {
	//fal(e,α,δ)={|e|αsign(e),|e|>δeδ1−α,|e|≤δ

	r := 0.0
	if math.Abs(e) > delta {
		r = math.Pow(math.Abs(e), alfa) * sign(e)
	} else {
		r = e / (math.Pow(delta, 1-alfa))
	}
	return r
}
