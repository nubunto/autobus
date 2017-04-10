package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"web"

	"github.com/julienschmidt/httprouter"
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

func handleCreateStop(e *Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		stops := e.Backend.Stops()

		var payload busStopPayload
		if err := payload.Decode(r.Body); err != nil {
			web.ErrorResponse(w, err, http.StatusBadRequest)
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
		if err := stops.Create(doc); err != nil {
			web.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
		web.Response{
			OK:     true,
			Status: http.StatusCreated,
		}.EncodeTo(w)
	}
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

func handleGetStops(e *Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

		stops := e.Backend.Stops()
		all, err := stops.GetAll(bson.M{
			"location": bson.M{
				"$near": bson.M{
					"$geometry": bson.M{
						"type":        "Point",
						"coordinates": []float64{longitude, latitude},
					},
					"$maxDistance": radius,
				},
			},
		})
		if err != nil {
			web.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		web.Response{
			OK:   true,
			Data: toBusStopResp(all),
		}.EncodeTo(w)
	}
}
