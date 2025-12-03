package controllers

import (
	"adaptive-moms/parameters"
	"math"
	"time"
)

/***************************************/
/**********  Basic PID  ****************/
/***************************************/

type BasicPID struct {
	Kp                float64
	Ki                float64
	Kd                float64
	Direction         float64
	PreviousError     float64
	SumPreviousErrors float64
	Out               float64
	DeltaTime         time.Duration
	Max               float64
	Min               float64
}

func (c *BasicPID) Initialise(p parameters.AllParameters) {
	c.Kp = p.Kp
	c.Ki = p.Ki
	c.Kd = p.Kd
	c.SumPreviousErrors = 0.0
	c.PreviousError = 0.0
	c.Direction = p.Direction
	c.DeltaTime = time.Duration(p.DeltaTime)
	c.Max = p.Max
	c.Min = p.Min
}

func (c *BasicPID) Update(p ...float64) float64 {

	r := p[0] // goal
	y := p[1] // plant output

	// compute error
	err := c.Direction * (r - y)

	// Proportional
	proportional := c.Kp * err

	// Integrator (page 49)
	integrator := (c.SumPreviousErrors + err) * c.Ki * c.DeltaTime.Seconds()

	// Differentiator (page 49)
	differentiator := c.Kd * (err - c.PreviousError) / c.DeltaTime.Seconds()

	// control law
	c.Out = proportional + integrator + differentiator

	//println("Controller:: Out:: ", c.Out)
	if c.Out > c.Max {
		c.Out = c.Max
	} else if c.Out < c.Min {
		c.Out = c.Min
	}

	c.PreviousError = err
	c.SumPreviousErrors += err

	return c.Out
}

/***************************************/
/********     Deadzone PID  ************/
/***************************************/

type DeadzonePID struct {
	Kp                float64
	Ki                float64
	Kd                float64
	Direction         float64
	PreviousError     float64
	SumPreviousErrors float64
	Out               float64
	DeltaTime         time.Duration
	Max               float64
	Min               float64
	Deadzone          float64
}

func (c *DeadzonePID) Initialise(p parameters.AllParameters) {

	c.Kp = p.Kp
	c.Ki = p.Ki
	c.Kd = p.Kd
	c.SumPreviousErrors = 0.0
	c.PreviousError = 0.0
	c.Direction = p.Direction
	c.DeltaTime = time.Duration(p.DeltaTime)
	c.Max = p.Max
	c.Min = p.Min
	c.Deadzone = p.Deadzone
}

func (c *DeadzonePID) Update(p ...float64) float64 {

	r := p[0] // goal
	y := p[1] // plant output

	// errors
	err := c.Direction * (r - y)

	if math.Abs(err) > c.Deadzone { // outside deadzone
		// Proportional
		proportional := c.Kp * err

		// Integrator (David page 49)
		integrator := (c.SumPreviousErrors + err) * c.Ki * c.DeltaTime.Seconds()

		// Differentiator (David page 49)
		differentiator := c.Kd * (err - c.PreviousError) / c.DeltaTime.Seconds()

		// pid output
		c.Out = proportional + integrator + differentiator
	} else { // inside deadzone
		c.Out = c.Out // No action
	}

	if c.Out > c.Max {
		c.Out = c.Max
	} else if c.Out < c.Min {
		c.Out = c.Min
	}

	c.PreviousError = err
	c.SumPreviousErrors += err

	return c.Out
}

/***************************************/
/******** Error Square Full ************/
/***************************************/

type ErrorSquareFull struct {
	Kp                float64
	Ki                float64
	Kd                float64
	Direction         float64
	PreviousError     float64
	SumPreviousErrors float64
	Out               float64
	DeltaTime         time.Duration
	Max               float64
	Min               float64
}

func (c *ErrorSquareFull) Initialise(p parameters.AllParameters) {
	c.Kp = p.Kp
	c.Ki = p.Ki
	c.Kd = p.Kd
	c.SumPreviousErrors = 0.0
	c.PreviousError = 0.0
	c.Direction = p.Direction
	c.DeltaTime = time.Duration(p.DeltaTime)
	c.Max = p.Max
	c.Min = p.Min
}

func (c *ErrorSquareFull) Update(p ...float64) float64 {

	r := p[0] // goal
	y := p[1] // plant output

	// errors
	err := c.Direction * (r - y)

	// Proportional
	proportional := c.Kp * err

	// Integrator (David page 49)
	integrator := (c.SumPreviousErrors + err) * c.Ki * c.DeltaTime.Seconds()

	// Differentiator (David page 49)
	differentiator := c.Kd * (err - c.PreviousError) / c.DeltaTime.Seconds()

	// pid output Page 109
	c.Out = math.Abs(err) * (proportional + integrator + differentiator)

	if c.Out > c.Max {
		c.Out = c.Max
	} else if c.Out < c.Min {
		c.Out = c.Min
	}

	c.PreviousError = err
	c.SumPreviousErrors += err

	return c.Out
}

/***********************************************/
/******** Error Square Proportional ************/
/***********************************************/

type ErrorSquareProportional struct {
	Kp                float64
	Ki                float64
	Kd                float64
	Direction         float64
	PreviousError     float64
	SumPreviousErrors float64
	Out               float64
	DeltaTime         time.Duration
	Max               float64
	Min               float64
}

func (c *ErrorSquareProportional) Initialise(p parameters.AllParameters) {
	c.Kp = p.Kp
	c.Ki = p.Ki
	c.Kd = p.Kd
	c.SumPreviousErrors = 0.0
	c.PreviousError = 0.0
	c.Direction = p.Direction
	c.DeltaTime = time.Duration(p.DeltaTime)
	c.Max = p.Max
	c.Min = p.Min
}

func (c *ErrorSquareProportional) Update(p ...float64) float64 {

	r := p[0] // goal
	y := p[1] // plant output

	// errors
	err := c.Direction * (r - y)

	// Proportional
	proportional := c.Kp * err

	// Integrator (David page 49)
	integrator := (c.SumPreviousErrors + err) * c.Ki * c.DeltaTime.Seconds()

	// Differentiator (David page 49)
	differentiator := c.Kd * (err - c.PreviousError) / c.DeltaTime.Seconds()

	// pid output Page 109
	c.Out = math.Abs(err)*proportional + integrator + differentiator

	if c.Out > c.Max {
		c.Out = c.Max
	} else if c.Out < c.Min {
		c.Out = c.Min
	}

	c.PreviousError = err
	c.SumPreviousErrors += err

	return c.Out
}

/***********************************************/
/************ Incremental PID ******************/
/***********************************************/

type IncrementalPID struct {
	Kp                    float64
	Ki                    float64
	Kd                    float64
	Direction             float64
	PreviousError         float64
	PreviousPreviousError float64
	SumPreviousErrors     float64
	Out                   float64
	DeltaTime             time.Duration
	Max                   float64
	Min                   float64
}

func (c *IncrementalPID) Initialise(p parameters.AllParameters) {
	c.Kp = p.Kp
	c.Ki = p.Ki
	c.Kd = p.Kd
	c.SumPreviousErrors = 0.0
	c.PreviousError = 0.0
	c.PreviousPreviousError = 0.0
	c.Direction = p.Direction
	c.DeltaTime = time.Duration(p.DeltaTime)
	c.Max = p.Max
	c.Min = p.Min
}

func (c *IncrementalPID) Update(p ...float64) float64 {
	r := p[0] // goal
	y := p[1] // plant output

	// errors
	err := c.Direction * (r - y)

	// Delta of the new PC
	deltaU := c.Kp*(err-c.PreviousError) + c.Ki*err*c.DeltaTime.Seconds() + c.Kd*(err-2*c.PreviousError+c.PreviousPreviousError)/c.DeltaTime.Seconds()

	// pid output
	c.Out = c.Out + deltaU // see page 106 why add an integrator

	if c.Out > c.Max {
		c.Out = c.Max
	} else if c.Out < c.Min {
		c.Out = c.Min
	}

	c.PreviousPreviousError = c.PreviousError
	c.PreviousError = err
	c.SumPreviousErrors += err

	return c.Out
}

/***********************************************/
/************ Setpoint Weighting ***************/
/***********************************************/

type SetPointWeighting struct {
	Kp                    float64
	Ki                    float64
	Kd                    float64
	Direction             float64
	PreviousError         float64
	PreviousPreviousError float64
	SumPreviousErrors     float64
	Out                   float64
	DeltaTime             time.Duration
	Max                   float64
	Min                   float64
	Alfa                  float64
	Beta                  float64
	Integrator            float64
}

func (c *SetPointWeighting) Initialise(p parameters.AllParameters) {
	c.Kp = p.Kp
	c.Ki = p.Ki
	c.Kd = p.Kd
	c.SumPreviousErrors = 0.0
	c.PreviousError = 0.0
	c.PreviousPreviousError = 0.0
	c.Direction = p.Direction
	c.DeltaTime = time.Duration(p.DeltaTime)
	c.Max = p.Max
	c.Min = p.Min
	c.Alfa = p.Alfa
	c.Beta = p.Beta
	c.Integrator = 0.0
}

func (c *SetPointWeighting) Update(p ...float64) float64 {

	r := p[0] // goal
	y := p[1] // plant output

	// errors
	err := c.Direction * (r - y)

	// Proportional
	proportional := c.Kp * c.Direction * (c.Alfa*r - y)

	// Integrator (page 49)
	c.Integrator += c.DeltaTime.Seconds() * err
	integrator := c.Integrator * c.Ki

	// Differentiator (page 108)
	differentiator := c.Kd * ((1-c.Beta)*r - y - c.PreviousError) / c.DeltaTime.Seconds()

	// control law
	c.Out = proportional + integrator + differentiator

	if c.Out > c.Max {
		c.Out = c.Max
	} else if c.Out < c.Min {
		c.Out = c.Min
	}

	c.PreviousError = err
	c.SumPreviousErrors += err

	return c.Out
}

/***********************************************/
/************ Smoothing PID ********************/
/***********************************************/

type SmoothingPID struct {
	Kp                     float64
	Ki                     float64
	Kd                     float64
	Direction              float64
	PreviousError          float64
	PreviousPreviousError  float64
	SumPreviousErrors      float64
	Out                    float64
	DeltaTime              time.Duration
	Max                    float64
	Min                    float64
	Alfa                   float64
	PreviousDifferentiator float64
}

func (c *SmoothingPID) Initialise(p parameters.AllParameters) {
	c.Kp = p.Kp
	c.Ki = p.Ki
	c.Kd = p.Kd
	c.SumPreviousErrors = 0.0
	c.PreviousError = 0.0
	c.PreviousPreviousError = 0.0
	c.PreviousDifferentiator = 0.0
	c.Direction = p.Direction
	c.DeltaTime = time.Duration(p.DeltaTime)
	c.Max = p.Max
	c.Min = p.Min
	c.Alfa = p.Alfa
}

func (c *SmoothingPID) Update(p ...float64) float64 {

	r := p[0] // goal
	y := p[1] // plant output

	// errors
	err := c.Direction * (r - y)

	// Proportional
	proportional := c.Kp * err

	// Integrator (David page 49)
	integrator := (c.SumPreviousErrors + err) * c.Ki * c.DeltaTime.Seconds()

	// smoothing the derivative term (page 104)
	differentiator := c.Kd * (c.Alfa*(err-c.PreviousError)/c.DeltaTime.Seconds() + (1-c.Alfa)*c.PreviousDifferentiator)
	c.PreviousDifferentiator = differentiator

	// pid output
	c.Out = proportional + integrator + differentiator

	if c.Out > c.Max {
		c.Out = c.Max
	} else if c.Out < c.Min {
		c.Out = c.Min
	}

	c.PreviousError = err
	c.SumPreviousErrors += err

	return c.Out
}
