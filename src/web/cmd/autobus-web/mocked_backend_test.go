package main

type mockedMongoBackend struct {
	lines *mockedLinesBackend
	stops *mockedStopsBackend
}

func (m *mockedMongoBackend) Lines() web.LinesBackend {
	return m.lines
}

func (m *mockedMongoBackend) Stops() web.StopsBackend {
	return m.stops
}

type mockedLinesBackend struct {
	lines []web.Line
}

type mockedStopsBackend struct {
	stops []web.BusStop
}
