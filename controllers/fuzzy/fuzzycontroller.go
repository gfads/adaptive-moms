package fuzzy

import (
	"adaptive-moms/controllers/fuzzy/deffuzification"
	"adaptive-moms/controllers/fuzzy/fuzzification"
	"adaptive-moms/parameters"
	"adaptive-moms/shared"
	"fmt"
	"os"
)

// fuzzification
const EXTREMELYPOSITIVE = "EP"
const LARGEPOSITIVE = "LP" // Large Positive
const SMALLPOSITIVE = "SP" // Small Positive
const ZERO = "ZE"          // Zero
const SMALLNEGATIVE = "SN" // Smal Negative
const LARGENEGATIVE = "LN" // Large Negative
const EXTREMELYNEGATIVE = "EN"

// deffuzification
const LARGEINCREASE = "LI"  // Large Positive
const MEDIUMINCREASE = "MI" // Small Positive
const SMALLINCREASE = "SI"  // Small Positive
const MAINTAIN = "MAINTAIN" // Zero
const SMALLDECREASE = "SD"  // Small Negative
const MEDIUMDECREASE = "MD" // Small Positive
const LARGEDECREASE = "LD"  // Large Negative

// Membership functions
const TRIANGULAR = "Triangular"
const GAUSSIAN = "Gaussian"
const PI = "Pi"
const RAMP = "Ramp"
const TRAPEZOIDAL = "Trapezoidal"

type FuzzyController struct {
	MembershipFunction    string
	DeffuzificationMethod string
	PC                    float64
	Max                   float64
	Min                   float64
	Out                   float64
}

func (c *FuzzyController) Initialise(p parameters.AllParameters) {
	c.DeffuzificationMethod = p.DeffuzificationMethod
	c.MembershipFunction = p.MembershipFunction
	c.PC = p.PC
	c.Min = p.Min
	c.Max = p.Max
	c.Out = 0.0
}

func (c *FuzzyController) Update(p ...float64) float64 {
	goal := p[0]
	rate := p[1]
	pc := p[2]

	e := goal - rate

	// 1. Fuzzification
	fuzzifiedSetError := fuzzyInput(e, c.MembershipFunction)

	// 2. apply rules
	output := applyRules(fuzzifiedSetError)

	// 3. Deffuzification
	f := deffuzification.Centroid{}
	u := f.Deffuzify(output)

	// Check the interval of the PC
	c.Out = pc + u
	if c.Out > c.Max {
		c.Out = c.Max
	} else if c.Out < c.Min {
		c.Out = c.Min
	}

	return c.Out
}

func applyRules(e map[string]float64) shared.OutputX {
	o := shared.OutputX{}

	// Rule 1:  IF e = EXTREMELYPOSITIVE THEN output = LARGEINCREASE
	o.Mx = append(o.Mx, e[EXTREMELYPOSITIVE])
	o.Out = append(o.Out, getMaxOutput(LARGEINCREASE))

	// Rule 2:  IF error LARGEPOSITIVE THEN output = MEDIUMINCREASE
	o.Mx = append(o.Mx, e[LARGEPOSITIVE])
	o.Out = append(o.Out, getMaxOutput(MEDIUMINCREASE))

	// Rule 3:  IF e = SMALLPOSITIVE THEN output = SMALLINCREASE
	o.Mx = append(o.Mx, e[SMALLPOSITIVE]) // saida = +1 s
	o.Out = append(o.Out, getMaxOutput(SMALLINCREASE))

	// Rule 4:  IF e = ZE THEN output = MAINTAIN
	o.Mx = append(o.Mx, e[ZERO]) // saida = 0
	o.Out = append(o.Out, getMaxOutput(MAINTAIN))

	// Rule 5:  IF e = SMALLNEGATIVE THEN output = SMALLDECREASE
	o.Mx = append(o.Mx, e[SMALLNEGATIVE]) // saida = -1 s
	o.Out = append(o.Out, getMaxOutput(SMALLDECREASE))

	// Rule 6:  IF e = LARGENEGATIVE THEN output = MEDIUMDECREASE
	o.Mx = append(o.Mx, e[LARGENEGATIVE])
	o.Out = append(o.Out, getMaxOutput(MEDIUMDECREASE))

	// Rule 7:  IF e = EXTREMELYNEGATIVE THEN output = LARGEDECREASE
	o.Mx = append(o.Mx, e[EXTREMELYNEGATIVE])
	o.Out = append(o.Out, getMaxOutput(LARGEDECREASE))

	//fmt.Printf("[%.2f %.2f %.2f %.2f %.2f %.2f %.2f]\n", o.Mx[0], o.Mx[1], o.Mx[2], o.Mx[3], o.Mx[4], o.Mx[5], o.Mx[6])
	//fmt.Printf("[%.2f %.2f %.2f %.2f %.2f %.2f %.2f]\n", o.Out[0], o.Out[1], o.Out[2], o.Out[3], o.Out[4], o.Out[5], o.Out[6])
	return o
}
func fuzzyInput(x float64, mf string) map[string]float64 {
	r := map[string]float64{}

	switch mf {
	case TRIANGULAR:
		f := fuzzification.Triangular{}
		r[EXTREMELYPOSITIVE] = f.Fuzzify(x, 1250, 5000, 10000)
		r[LARGEPOSITIVE] = f.Fuzzify(x, 500, 1250, 2000)          //500,750,1000
		r[SMALLPOSITIVE] = f.Fuzzify(x, 0, 625, 1250)             // 0, 500,1000
		r[ZERO] = f.Fuzzify(x, -500, 0, 500)                      // -500,0,500
		r[SMALLNEGATIVE] = f.Fuzzify(x, -1250, -625, 0)           //-1000,-500,0
		r[LARGENEGATIVE] = f.Fuzzify(x, -2000, -1250, -500)       // -1000,-750,-500
		r[EXTREMELYNEGATIVE] = f.Fuzzify(x, -1250, -5000, -10000) // -1000,-750,-500
	case GAUSSIAN:
		f := fuzzification.Gaussian{}
		r[EXTREMELYPOSITIVE] = f.Fuzzify(x, 3000.0, 0.01)
		r[LARGEPOSITIVE] = f.Fuzzify(x, 1500.0, 0.01)      //500,750,1000
		r[SMALLPOSITIVE] = f.Fuzzify(x, 500.0, 0.01)       // 0, 500,1000
		r[ZERO] = f.Fuzzify(x, 0.0, 0.1)                   // -500,0,500
		r[SMALLNEGATIVE] = f.Fuzzify(x, -500.0, 0.01)      //-1000,-500,0
		r[LARGENEGATIVE] = f.Fuzzify(x, -1500.0, 0.01)     // -1000,-750,-500
		r[EXTREMELYNEGATIVE] = f.Fuzzify(x, -3000.0, 0.01) // -1000,-750,-500
	case PI:
		f := fuzzification.Pi{}
		r[EXTREMELYPOSITIVE] = f.Fuzzify(x, 1250, 2500, 5000, 10000)
		r[LARGEPOSITIVE] = f.Fuzzify(x, 500, 250, 1750, 2000)            //500,750,1000
		r[SMALLPOSITIVE] = f.Fuzzify(x, 0, 250, 1000, 1250)              // 0, 500,1000
		r[ZERO] = f.Fuzzify(x, -500, -250, 250, 500)                     // -500,0,500
		r[SMALLNEGATIVE] = f.Fuzzify(x, -1250, -1000, -250, 0)           //-1000,-500,0
		r[LARGENEGATIVE] = f.Fuzzify(x, -2000, -1750, -250, -500)        // -1000,-750,-500
		r[EXTREMELYNEGATIVE] = f.Fuzzify(x, -10000, -5000, -2500, -1250) // -1000,-750,-500
	default:
		shared.ErrorHandler(shared.GetFunction(), "Error: Membership function invalid!")
	}

	/*
		fmt.Printf("Error = %.2f FuzzifiedError [%.2f %.2f %.2f %.2f %.2f %.2f %.2f]\n", x,
			r[EXTREMELYNEGATIVE],
			r[LARGENEGATIVE],
			r[SMALLNEGATIVE],
			r[ZERO],
			r[SMALLPOSITIVE],
			r[LARGEPOSITIVE],
			r[EXTREMELYPOSITIVE])
	*/
	return r
}
func fuzzyOutput(n float64, mf string) map[string]float64 {
	r := map[string]float64{}

	switch mf {

	case GAUSSIAN:
		f := fuzzification.Gaussian{}
		r[LARGEINCREASE] = f.Fuzzify(n, 3.0, 0.01)  // original = 2
		r[MEDIUMINCREASE] = f.Fuzzify(n, 2.0, 0.01) // original = 2
		r[SMALLINCREASE] = f.Fuzzify(n, 1.0, 0.01)  // original = 1
		r[MAINTAIN] = f.Fuzzify(n, 0.0, 0.01)
		r[SMALLDECREASE] = f.Fuzzify(n, -1.0, 0.01)  // original=-1
		r[MEDIUMDECREASE] = f.Fuzzify(n, -2.0, 0.01) // original=-1
		r[LARGEDECREASE] = f.Fuzzify(n, -3.0, 0.01)  // original= -2
	case TRIANGULAR:
		f := fuzzification.Triangular{}
		r[LARGEINCREASE] = f.Fuzzify(n, 1.0, 2.0, 3.0)
		r[SMALLINCREASE] = f.Fuzzify(n, 0.0, 1.0, 2.0)
		r[MAINTAIN] = f.Fuzzify(n, 0.5, 0.0, -0.5)
		r[SMALLDECREASE] = f.Fuzzify(n, -2.0, -1.0, 0.0)
		r[LARGEDECREASE] = f.Fuzzify(n, -3.0, -2.0, -1.0)
	case TRAPEZOIDAL:
		f := fuzzification.Trapezoidal{}
		r[LARGEINCREASE] = f.Fuzzify(n, 2.0, 2.5, 3.5, 4.0)
		r[SMALLINCREASE] = f.Fuzzify(n, 1.0, 1.5, 2.5, 3.0)
		r[MAINTAIN] = f.Fuzzify(n, -1.0, -0.5, 0.5, 1.0)
		r[SMALLDECREASE] = f.Fuzzify(n, -3.0, -2.5, -1.5, -1.0)
		r[LARGEDECREASE] = f.Fuzzify(n, -4.0, -3.5, -2.5, -3.0)
	default:
		fmt.Println("Error: Membership function invalid!")
		os.Exit(0)
	}
	return r
}
func getMaxOutput(s string) float64 {
	r := 0.0
	max := -10000.0

	for i := -3.0; i <= 3.0; i += 0.5 { // TODO
		v := fuzzyOutput(i, GAUSSIAN)
		if v[s] > max {
			max = v[s]
			r = i
		}
	}
	return r
}
