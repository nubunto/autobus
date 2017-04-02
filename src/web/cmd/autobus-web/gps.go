package main

import (
	"time"
	"web/app"

	"github.com/goadesign/goa"
	bson "gopkg.in/mgo.v2/bson"
)

// GPSController implements the GPS resource.
type GPSController struct {
	*goa.Controller
	*Env
}

type Message struct {
	MessageHead string
	ID          string
	Type        string
	Valid       bool
	Latitude    float64
	Longitude   float64
	DateTime    time.Time
	Speed       float64
	Direction   int64
	Status      string
}

type MessageData []Message

func (md MessageData) toGpsMedia() []*app.GpsMedia {
	gpsData := make([]*app.GpsMedia, len(md))
	for i, m := range md {
		p := new(app.GpsMedia)
		direction := int(m.Direction)
		p.Head = &m.MessageHead
		p.ID = &m.ID
		p.Type = &m.Type
		p.Valid = &m.Valid
		p.Latitude = &m.Latitude
		p.Longitude = &m.Longitude
		p.DateTime = &m.DateTime
		p.Speed = &m.Speed
		p.Direction = &direction
		p.Status = &m.Status
		gpsData[i] = p
	}
	return gpsData
}

// NewGPSController creates a GPS controller.
func NewGPSController(service *goa.Service, env *Env) *GPSController {
	return &GPSController{
		Controller: service.NewController("GPSController"),
		Env:        env,
	}
}

// Show runs the show action.
func (c *GPSController) Show(ctx *app.ShowGPSContext) error {
	transient := c.DB("autobus").C("gps_data_transient")
	var gpsData MessageData
	if err := transient.Find(bson.M{}).All(&gpsData); err != nil {
		return err
	}
	return ctx.OK(gpsData.toGpsMedia())
}
