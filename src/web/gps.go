package web

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type GPS struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	MessageHead string        `json:"head" bson:"head"`
	GPSID       string        `json:"gps_id" bson:"gps_id"`
	Type        string        `json:"type"`
	Valid       bool          `json:"valid"`
	Loc         *Location     `json:"location"`
	DateTime    time.Time     `json:"date_time"`
	Speed       float64       `json:"speed"`
	Direction   int64         `json:"direction"`
	Status      string        `json:"status"`
}

type Location struct {
	Type        string `json:"type"`
	Coordinates string `json:"coordinates"`
}
