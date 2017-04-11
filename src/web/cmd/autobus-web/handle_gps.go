package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"web"
)

func handleGetGPSTransient(e *Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		gps := e.Backend.GPSTransient()

		all, err := gps.GetAll(nil)
		if err != nil {
			web.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		web.OK(w, all)
	}
}
