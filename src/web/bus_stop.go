package web

import (
	"gopkg.in/mgo.v2/bson"
)

type BusStop struct {
	ID       bson.ObjectId
	Name     string `json:"name"`
	Location BusStopLocation
}

type BusStopLocation struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
