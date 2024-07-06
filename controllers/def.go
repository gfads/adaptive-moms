package controllers

import (
	"adaptive-moms/controllers/fuzzy"
	"adaptive-moms/parameters"
	"adaptive-moms/shared"
)

var ControllerSet = map[string]Controller{
	shared.BasicPID:                new(BasicPID),
	shared.DeadzonePID:             new(DeadzonePID),
	shared.HPA:                     new(HPA),
	shared.ASTAR:                   new(ASTAR),
	shared.ErrorSquareFull:         new(ErrorSquareFull),
	shared.ErrorSquareProportional: new(ErrorSquareProportional),
	shared.IncrementalPID:          new(IncrementalPID),
	shared.SetPointWeighting:       new(SetPointWeighting),
	shared.SmoothingPID:            new(SmoothingPID),
	shared.FuzzyController:         new(fuzzy.FuzzyController),
}

type Controller interface {
	Initialise(parameters parameters.AllParameters)
	Update(...float64) float64
}

func NewController(t string) Controller {

	_, ok := ControllerSet[t]
	if !ok {
		shared.ErrorHandler(shared.GetFunction(), "Controller type does not exist ["+t+"]")
	}
	return ControllerSet[t]
}
