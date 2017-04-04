package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"web"

	"github.com/dimfeld/httptreemux"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

type busStopPayload struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (b *busStopPayload) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(b)
}

type busStopResp struct {
	ID       bson.ObjectId       `json:"id"`
	Name     string              `json:"name"`
	Location busStopRespLocation `json:"location"`
}

type busStopRespLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func toBusStopResp(b []web.BusStop) []busStopResp {
	ret := make([]busStopResp, len(b))
	for i, c := range b {
		ret[i] = busStopResp{
			ID:   c.ID,
			Name: c.Name,
			Location: busStopRespLocation{
				Longitude: c.Location.Coordinates[0],
				Latitude:  c.Location.Coordinates[1],
			},
		}
	}
	return ret
}

func handleCreateStop(e *Env) httptreemux.HandlerFunc {
	return httptreemux.HandlerFunc(func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		e.Debug("creating bus stop")
		var payload busStopPayload
		if err := payload.Decode(r.Body); err != nil {
			web.ErrorResponse(w, err, http.StatusBadRequest)
			e.Debug("error decoding the payload", zap.Error(err))
			return
		}
		doc := web.BusStop{
			ID:   bson.NewObjectId(),
			Name: payload.Name,
			Location: web.BusStopLocation{
				Type:        "Point",
				Coordinates: []float64{payload.Longitude, payload.Latitude},
			},
		}
		stops := e.DB("autobus").C("stops")
		if err := stops.Insert(doc); err != nil {
			web.ErrorResponse(w, err, http.StatusInternalServerError)
			e.Debug("error inserting bus stop", zap.Error(err))
			return
		}
		web.Response{
			OK:      true,
			Message: "Created successfully",
			Status:  http.StatusCreated,
		}.EncodeTo(w)
	})
}

func parseFloatOrStreamError(f string, w http.ResponseWriter) (p float64, errorFound bool) {
	var err error
	p, err = strconv.ParseFloat(f, 64)
	if err != nil {
		web.ErrorResponse(w, err, http.StatusBadRequest)
		return 0, true
	}
	return
}

func handleGetStops(e *Env) httptreemux.HandlerFunc {
	return httptreemux.HandlerFunc(func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		query := r.URL.Query()
		radius, stop := parseFloatOrStreamError(query.Get("radius"), w)
		if stop {
			return
		}

		latitude, stop := parseFloatOrStreamError(query.Get("latitude"), w)
		if stop {
			return
		}

		longitude, stop := parseFloatOrStreamError(query.Get("longitude"), w)
		if stop {
			return
		}

		all := make([]web.BusStop, 0)
		stops := e.DB("autobus").C("stops")
		if err := stops.Find(bson.M{
			"location": bson.M{
				"$near": bson.M{
					"$geometry": bson.M{
						"type":        "Point",
						"coordinates": []float64{longitude, latitude},
					},
					"$maxDistance": radius,
				},
			},
		}).All(&all); err != nil {
			web.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		web.Response{
			OK:      true,
			Message: "Retrieved successfully",
			Data:    toBusStopResp(all),
		}.EncodeTo(w)
	})
}
