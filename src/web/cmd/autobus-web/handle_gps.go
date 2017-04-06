package main

import (
	"domain"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"web"
)

type gpsResp struct {
	MessageHead string       `json:"message_head"`
	ID          string       `json:"id"`
	Type        string       `json:"type"`
	Valid       bool         `json:"valid"`
	Loc         *gpsLocation `json:"location"`
	DateTime    time.Time    `json:"date_time"`
	Speed       float64      `json:"speed"`
	Direction   int64        `json:"direction"`
	Status      string       `json:"status"`
}

func (g *gpsResp) UnmarshalDomain(m domain.GPSMessage) {
	g.MessageHead = m.MessageHead
	g.ID = m.ID
	g.Type = m.Type
	g.Valid = m.Valid
	g.Loc = new(gpsLocation)
	g.Loc.Latitude = m.Loc.Coordinates[0]
	g.Loc.Latitude = m.Loc.Coordinates[1]
	g.DateTime = m.DateTime
	g.Speed = m.Speed
	g.Direction = m.Direction
	g.Status = m.Status
}

type gpsLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func handleGetGPS(e *Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		transient := e.DB("autobus").C("gps_data_transient")
		all := make([]*gpsResp, 0)
		raw := make([]domain.GPSMessage, 0)
		if err := transient.Find(nil).All(&raw); err != nil {
			web.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
		for i, m := range raw {
			g := new(gpsResp)
			g.UnmarshalDomain(m)
			all[i] = g
		}
		web.Response{
			OK:      true,
			Message: "OK",
			Data:    all,
		}.EncodeTo(w)
	}
}
