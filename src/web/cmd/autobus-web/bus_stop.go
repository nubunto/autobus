package main

import (
	"web/app"

	"github.com/goadesign/goa"
	"gopkg.in/mgo.v2/bson"
)

type BusStopController struct {
	*goa.Controller
	*Env
}

type BusStop struct {
	Name string
	Loc  *GeoJSON
}

func (b *BusStop) FromPayload(p *app.BusStopPayload) {
	b.Name = p.Name
	b.Loc = &GeoJSON{
		Type: "Point",
		Coordinates: []float64{
			p.Longitude,
			p.Latitude,
		},
	}
}

func NewBusStopController(service *goa.Service, e *Env) *BusStopController {
	return &BusStopController{
		Controller: service.NewController("BusStopController"),
		Env:        e,
	}
}

func (c *BusStopController) Create(ctx *app.CreateBusStopContext) error {
	busStops := c.DB("autobus").C("stops")
	payload := ctx.Payload
	busStop := new(BusStop)
	busStop.FromPayload(payload)
	if err := busStops.Insert(busStop); err != nil {
		return err
	}
	return nil
}

func (c *BusStopController) Nearest(ctx *app.NearestBusStopContext) error {
	latitude := ctx.Latitude
	longitude := ctx.Longitude
	radius := ctx.Radius
	busStops := c.DB("autobus").C("stops")
	all := make([]*app.BusStopMedia, 0)
	if err := busStops.Find(bson.M{
		"loc": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longitude, latitude},
				},
				"$maxDistance": radius,
			},
		},
	}).Select(bson.M{
		"loc.type": 0,
	}).All(&all); err != nil {
		return err
	}
	return ctx.OK(all)
}
