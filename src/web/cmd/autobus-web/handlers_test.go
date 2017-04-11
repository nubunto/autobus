package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"web"

	"github.com/julienschmidt/httprouter"
)

func newEnvWithMockedMongoBackend() *Env {
	return &Env{
		Backend: newMockedMongoBackend(),
	}
}

func newMockedMongoBackend() *mockedMongoBackend {
	return &mockedMongoBackend{
		lines: &mockedLinesBackend{
			lines: make([]web.Line, 0),
		},
		stops: &mockedStopsBackend{
			stops: make([]web.BusStop, 0),
		},
		gps: &mockedGPSBackend{
			gps: make([]web.GPSData, 0),
		},
		gpsTransient: &mockedGPSBackend{
			gps: make([]web.GPSData, 0),
		},
	}
}

type hooks struct {
	beforeHandler func(b web.Backend)
	afterHandler  func(t *testing.T, b web.Backend)
}

func TestHandlers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		method       string
		requestPath  string
		registerPath string
		query        string
		env          *Env
		hooks        *hooks
		handler      func(*Env) httprouter.Handle
		payload      io.Reader
		assert       func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:         "GetLines",
			method:       "GET",
			registerPath: "/lines",
			requestPath:  "/lines",
			env:          newEnvWithMockedMongoBackend(),
			handler:      handleGetLines,
			assert: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Error("should have status code 200 ok, has:", rec.Code)
				}
				body := rec.Body.Bytes()
				var response struct {
					OK      bool       `json:"ok"`
					Message string     `json:"message"`
					Data    []web.Line `json:"data"`
				}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Error("should be valid json:", rec.Body.String())
				}
				if !response.OK {
					t.Error("should be a valid response:", response)
				}
				if len(response.Data) != 0 {
					t.Error("should have 0 lines from mocked backend:", response.Data)
				}
			},
		},

		{
			name:         "GetLinesWithStopID",
			method:       "GET",
			registerPath: "/lines/:stopID",
			requestPath:  "/lines/58ebed69183add0001d82019",
			env:          newEnvWithMockedMongoBackend(),
			hooks: &hooks{
				beforeHandler: func(b web.Backend) {
					b.Lines().Create(web.Line{
						Stops: []web.StopID{
							{ID: bson.ObjectIdHex("58ebed69183add0001d82019")},
						},
					})
				},
			},
			handler: handleGetLinesWithStopID,
			assert: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Error("should have status code 200 ok, has:", rec.Code)
				}
				body := rec.Body.Bytes()
				var response struct {
					OK      bool       `json:"ok"`
					Message string     `json:"message"`
					Data    []web.Line `json:"data"`
				}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Error("should be valid json:", err)
				}
				if !response.OK {
					t.Error("should be a valid response:", response)
				}
				if len(response.Data) != 1 {
					t.Error("should have 1 line associated with the given ID", response.Data)
				}
			},
		},

		{
			name:         "CreateLine",
			method:       "POST",
			registerPath: "/lines",
			requestPath:  "/lines",
			env:          newEnvWithMockedMongoBackend(),
			payload: strings.NewReader(`
				{
					"name": "244",
					"hours": ["08:00", "09:00"],
					"stops": [
						{"id": "58e6ab56d8959f2403cc4eda"},
						{"id": "b58e6ab56d8959f2403cc4ed"}
					]
				}
			`),
			handler: handleCreateLine,
			assert: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusCreated {
					t.Errorf("should have status code %d (%s), has %d (%s)",
						http.StatusCreated, http.StatusText(http.StatusCreated),
						rec.Code, http.StatusText(rec.Code))
				}
				body := rec.Body.Bytes()
				var response struct {
					OK      bool       `json:"ok"`
					Message string     `json:"message"`
					Data    []web.Line `json:"data"`
				}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Error("should be valid json:", err)
				}
				if !response.OK {
					t.Error("should be a valid response:", response)
				}
			},
			hooks: &hooks{
				afterHandler: func(t *testing.T, b web.Backend) {
					lines, err := b.Lines().GetAll(nil)
					if err != nil {
						t.Error("should retrieve all lines successfully from mock")
					}
					if len(lines) != 1 {
						t.Error("after handler, it should have created 1 line, created", len(lines))
					}
				},
			},
		},

		{
			name:         "CreateBusStop",
			method:       "POST",
			registerPath: "/stops",
			requestPath:  "/stops",
			env:          newEnvWithMockedMongoBackend(),
			payload: strings.NewReader(`
				{
					"name": "Sherlock Holmes stop",
					"address": "221b Baker Street",
					"latitude": -24.98329,
					"longitude": 87.87393
				}
			`),
			handler: handleCreateStop,
			assert: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusCreated {
					t.Errorf("should have status code %d (%s), has %d (%s)",
						http.StatusCreated, http.StatusText(http.StatusCreated),
						rec.Code, http.StatusText(rec.Code))
				}
				body := rec.Body.Bytes()
				var response struct {
					OK      bool          `json:"ok"`
					Message string        `json:"message"`
					Data    []web.BusStop `json:"data"`
				}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Error("should be valid json:", err)
				}
				if !response.OK {
					t.Error("should be a valid response:", response)
				}
			},
			hooks: &hooks{
				afterHandler: func(t *testing.T, b web.Backend) {
					stops, err := b.Stops().GetAll(nil)
					if err != nil {
						t.Error("should retrieve all stops successfully from mock")
					}
					if len(stops) != 1 {
						t.Error("after handler, it should have created 1 new stop, created", len(stops))
					}
				},
			},
		},

		{
			name:         "GetBusStop",
			method:       "GET",
			registerPath: "/stops",
			requestPath:  "/stops",
			query:        "latitude=1&longitude=2&radius=10",
			env:          newEnvWithMockedMongoBackend(),
			handler:      handleGetStops,
			assert: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Errorf("should have status code %d (%s), has %d (%s)",
						http.StatusOK, http.StatusText(http.StatusOK),
						rec.Code, http.StatusText(rec.Code))
				}
				body := rec.Body.Bytes()
				var response struct {
					OK      bool          `json:"ok"`
					Message string        `json:"message"`
					Data    []web.BusStop `json:"data"`
				}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Error("should be valid json:", err)
				}
				if !response.OK {
					t.Error("should be a valid response:", response)
				}
				if len(response.Data) != 0 {
					t.Error("should have 0 stops", response.Data)
				}
			},
		},

		{
			name:         "GetBusStopNoParams",
			method:       "GET",
			registerPath: "/stops",
			requestPath:  "/stops",
			env:          newEnvWithMockedMongoBackend(),
			handler:      handleGetStops,
			assert: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("should have status code %d (%s), has %d (%s)",
						http.StatusBadRequest, http.StatusText(http.StatusBadRequest),
						rec.Code, http.StatusText(rec.Code))
				}
				body := rec.Body.Bytes()
				var response struct {
					OK      bool   `json:"ok"`
					Message string `json:"message"`
				}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Error("should be valid json:", err)
				}
				if response.OK {
					t.Error("should be a invalid response:", response)
				}
			},
		},

		{
			name:         "GetGPSData",
			method:       "GET",
			registerPath: "/live",
			requestPath:  "/live",
			env:          newEnvWithMockedMongoBackend(),
			handler:      handleGetGPSTransient,
			assert: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Errorf("should have status code %d (%s), has %d (%s)",
						http.StatusOK, http.StatusText(http.StatusOK),
						rec.Code, http.StatusText(rec.Code))
				}
				body := rec.Body.Bytes()
				var response struct {
					OK      bool   `json:"ok"`
					Message string `json:"message"`
				}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Error("should be valid json:", err)
				}
				if !response.OK {
					t.Error("should be a valid response:", response)
				}
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			mux := httprouter.New()
			mux.Handle(test.method, test.registerPath, test.handler(test.env))

			req := httptest.NewRequest(test.method, "http://example.com"+test.requestPath+"?"+test.query, test.payload)
			w := httptest.NewRecorder()

			var h hooks
			if test.hooks != nil {
				h = *test.hooks
			}
			if h.beforeHandler != nil {
				h.beforeHandler(test.env.Backend)
			}
			mux.ServeHTTP(w, req)
			if h.afterHandler != nil {
				h.afterHandler(tt, test.env.Backend)
			}

			test.assert(tt, w)
		})
	}
}
