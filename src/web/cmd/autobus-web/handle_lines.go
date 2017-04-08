package main

import (
	"encoding/json"
	"io"
	"net/http"
	"web"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

type busStopID struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}

type busStopIDPayload []busStopID

func (b busStopIDPayload) toStopID() []web.StopID {
	ret := make([]web.StopID, len(b))
	for i, stop := range b {
		ret[i] = web.StopID{
			ID: stop.ID,
		}
	}
	return ret
}

type lineRoutePayload struct {
	// TODO
	Type string `json:"type"`
}

type linePayload struct {
	Name  string           `json:"name"`
	Hours []string         `json:"hours"`
	Stops busStopIDPayload `json:"stops"`
	//Route lineRoutePayload `json:"route"`
}

func (lp *linePayload) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(lp)
}

func handleGetLines(e *Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		lines := e.Backend.Lines()

		finder := make(bson.M, 1)
		query := r.URL.Query()

		if stopID := query.Get("stop_id"); stopID != "" {
			finder["stops"] = bson.M{
				"_id": stopID,
			}
		}

		all, err := lines.GetAll(finder)
		if err != nil {
			web.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		web.Response{
			OK:      true,
			Message: "Retrieved lines successfully",
			Data:    all,
		}.EncodeTo(w)
	}
}

func handleCreateLine(e *Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		lines := e.Backend.Lines()

		var payload linePayload
		if err := payload.Decode(r.Body); err != nil {
			web.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		doc := web.Line{
			Name:  payload.Name,
			Hours: payload.Hours,
			Stops: payload.Stops.toStopID(),
			// TODO: routes
		}
		if err := lines.Create(doc); err != nil {
			web.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
		web.Response{
			OK:     true,
			Status: http.StatusCreated,
		}.EncodeTo(w)
	}
}
