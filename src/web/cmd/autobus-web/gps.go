package main

import (
	"database/sql"

	"github.com/goadesign/goa"
	"web/app"
)

// GPSController implements the GPS resource.
type GPSController struct {
	*goa.Controller
	*sql.DB
}

// NewGPSController creates a GPS controller.
func NewGPSController(service *goa.Service, db *sql.DB) *GPSController {
	return &GPSController{
		Controller: service.NewController("GPSController"),
		DB:         db,
	}
}

// Show runs the show action.
func (c *GPSController) Show(ctx *app.ShowGPSContext) error {
	var gpsData []*app.GpsMedia
	rows, err := c.DB.Query("SELECT id, time, date, longitude, latitude, status FROM gps_data")
	if err != nil {
		return err
	}
	for rows.Next() {
		var gps app.GpsMedia
		err = rows.Scan(&gps.ID, &gps.Time, &gps.Date, &gps.Longitude, &gps.Latitude, &gps.Status)
		if err != nil {
			return err
		}
		gpsData = append(gpsData, &gps)
	}
	return ctx.OK(gpsData)
}
