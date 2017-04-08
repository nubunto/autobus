package main

func TestGetBusStops(t *testing.T) {
	env := &Env{
		Backend: &mockedMongoBackend{
			lines: &mockedLinesBackend{},
			stops: &mockedStopsBackend{},
		},
	}
	mux := httprouter.New()
	mux.GET("/", handleGetBusStops(env))
	server := httptest.NewServer()
}
