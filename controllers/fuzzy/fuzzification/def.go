package fuzzification

import "math"

// Membership functions
type Triangular struct{}
type Ramp struct{}
type Gaussian struct{}
type Trapezoidal struct{}
type Pi struct{}

func (Triangular) Fuzzify(x float64, a float64, b float64, c float64) float64 {
	return math.Max(0, math.Min((x-a)/(b-a), (c-x)/(c-b)))
}

func (Ramp) Fuzzify(x float64, a float64, b float64) float64 {
	if x <= a {
		return 0
	} else if x >= b {
		return 1
	} else {
		return (x - a) / (b - a)
	}
}

func (Gaussian) Fuzzify(x float64, m float64, d float64) float64 {
	temp := -math.Pow(x-m, 2.0) / 2.0 * math.Pow(d, 2.0)
	r := math.Exp(temp)

	return r
}

func (Pi) Fuzzify(x float64, a float64, b float64, c float64, d float64) float64 {
	r := 0.0

	if x <= a {
		r = 0.0
	}
	if a <= x && x <= ((a+b)/2.0) {
		r = 2 * math.Pow((x-a)/(b-a), 2.0)
	}
	if b <= x && x <= c {
		r = 1
	}
	if c <= x && x <= ((c+d)/2.0) {
		r = 1 - 2*math.Pow((x-c)/(d-c), 2.0)
	}
	if (c+d)/2 <= x && x <= d {
		r = 2 * math.Pow((x-d)/(d-c), 2.0)
	}
	if x >= d {
		r = 0
	}

	return r
}

func (Trapezoidal) Fuzzify(x float64, a float64, b float64, c float64, d float64) float64 {
	min1 := math.Min(x-a/b-a, 1)
	r := math.Max(math.Min(min1, d-x/d-c), 0)

	return r
}
