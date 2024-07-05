package deffuzification

import "adaptive-moms/shared"

// Deffuzification methods
type Centroid struct{}

func (Centroid) Deffuzify(output shared.OutputX) float64 {

	numerator := 0.0
	denominator := 0.0

	for i := 0; i < len(output.Mx); i++ {
		numerator = numerator + output.Mx[i]*output.Out[i]
		denominator = denominator + output.Mx[i]
	}
	u := 0.0
	if denominator == 0 {
		u = 1
	} else {
		u = numerator / denominator
	}
	return u
}
