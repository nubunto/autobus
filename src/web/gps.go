package main

import (
	"github.com/goadesign/goa"
	"web/app"
)

// GPSController implements the GPS resource.
type GPSController struct {
	*goa.Controller
}

// NewGPSController creates a GPS controller.
func NewGPSController(service *goa.Service) *GPSController {
	return &GPSController{Controller: service.NewController("GPSController")}
}

// Show runs the show action.
func (c *GPSController) Show(ctx *app.ShowGPSContext) error {
	// GPSController_Show: start_implement

	// Put your logic here

	// GPSController_Show: end_implement
	return nil
}
