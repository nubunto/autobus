package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	}
}

func TestHandlers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		method  string
		path    string
		query   string
		env     *Env
		handler func(*Env) httprouter.Handle
		payload io.Reader
		assert  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:    "GetLines",
			method:  "GET",
			path:    "/lines",
			env:     newEnvWithMockedMongoBackend(),
			payload: nil,
			handler: handleGetLines,
			assert: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Error("should have status code 200 ok, has:", rec.Code)
				}
			},
		},

		{
			name:   "CreateLine",
			method: "POST",
			path:   "/lines",
			env:    newEnvWithMockedMongoBackend(),
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
			},
		},

		{
			name:   "CreateBusStop",
			method: "POST",
			path:   "/stops",
			env:    newEnvWithMockedMongoBackend(),
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
			},
		},

		{
			name:    "GetBusStop",
			method:  "GET",
			path:    "/stops",
			query:   "latitude=1&longitude=2&radius=10",
			env:     newEnvWithMockedMongoBackend(),
			handler: handleGetStops,
			assert: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Errorf("should have status code %d (%s), has %d (%s)",
						http.StatusOK, http.StatusText(http.StatusOK),
						rec.Code, http.StatusText(rec.Code))
				}
				var response web.Response
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Error("response should be valid JSON, instead got error:", err)
				}
				if !response.OK {
					t.Error("response should be ok")
				}
			},
		},

		{
			name:    "GetBusStopNoParams",
			method:  "GET",
			path:    "/stops",
			env:     newEnvWithMockedMongoBackend(),
			handler: handleGetStops,
			assert: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("should have status code %d (%s), has %d (%s)",
						http.StatusBadRequest, http.StatusText(http.StatusBadRequest),
						rec.Code, http.StatusText(rec.Code))
				}
				var response web.Response
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Error("response should be valid JSON, instead got error:", err)
				}
				if response.OK {
					t.Error("response should be not be ok")
				}
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			mux := httprouter.New()
			mux.Handle(test.method, test.path, test.handler(test.env))

			req := httptest.NewRequest(test.method, "http://example.com"+test.path+"?"+test.query, test.payload)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)
			test.assert(tt, w)
		})
	}
}
