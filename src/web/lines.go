package web

import "gopkg.in/mgo.v2/bson"

type Line struct {
	ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name  string        `json:"name"`
	Hours []string      `json:"hours"`
	Stops []StopID      `json:"stops"`
	Route LineRoute     `json:"routes"`
}

type StopID struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`
}

type LineRoute struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}
