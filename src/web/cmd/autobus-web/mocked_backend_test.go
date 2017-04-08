package main

import (
	"errors"
	"web"

	"gopkg.in/mgo.v2/bson"
)

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

func (m *mockedMongoBackend) Close() error {
	// noop
	return nil
}

type mockedLinesBackend struct {
	lines []web.Line
}

type mockedStopsBackend struct {
	stops []web.BusStop
}

func (ml *mockedLinesBackend) GetAll(selector interface{}) ([]web.Line, error) {
	return ml.lines, nil
}

func (ml *mockedLinesBackend) GetOne(selector interface{}) (*web.Line, error) {
	if len(ml.lines) == 0 {
		return nil, errors.New("empty lines collection")
	}
	index, ok := selector.(int)
	if ok {
		return &ml.lines[index], nil
	}
	b, ok := selector.(bson.M)
	if ok {
		id := b["id"].(bson.ObjectId)
		for _, l := range ml.lines {
			if l.ID.Hex() == id.Hex() {
				return &l, nil
			}
		}
	}
	return nil, errors.New("not found")
}

func (ml *mockedLinesBackend) Create(doc interface{}) error {
	line := doc.(web.Line)
	ml.lines = append(ml.lines, line)
	return nil
}

func (ml *mockedLinesBackend) Update(selector, update interface{}) error {
	index := selector.(int)
	doc := update.(web.Line)
	ml.lines[index] = doc
	return nil
}

func (ml *mockedLinesBackend) Delete(selector interface{}) error {
	index := selector.(int)
	ml.lines = append(ml.lines[:index], ml.lines[index+1:]...)
	return nil
}

func (ml *mockedLinesBackend) Close() error {
	// noop
	return nil
}

func (ml *mockedStopsBackend) GetAll(selector interface{}) ([]web.BusStop, error) {
	return ml.stops, nil
}

func (ml *mockedStopsBackend) GetOne(selector interface{}) (*web.BusStop, error) {
	if len(ml.stops) == 0 {
		return nil, errors.New("empty lines collection")
	}
	index, ok := selector.(int)
	if ok {
		return &ml.stops[index], nil
	}
	b, ok := selector.(bson.M)
	if ok {
		id := b["id"].(bson.ObjectId)
		for _, l := range ml.stops {
			if l.ID.Hex() == id.Hex() {
				return &l, nil
			}
		}
	}
	return nil, errors.New("not found")
}

func (ml *mockedStopsBackend) Create(doc interface{}) error {
	stop := doc.(web.BusStop)
	ml.stops = append(ml.stops, stop)
	return nil
}

func (ml *mockedStopsBackend) Update(selector, update interface{}) error {
	index := selector.(int)
	doc := update.(web.BusStop)
	ml.stops[index] = doc
	return nil
}

func (ml *mockedStopsBackend) Delete(selector interface{}) error {
	index := selector.(int)
	ml.stops = append(ml.stops[:index], ml.stops[index+1:]...)
	return nil
}

func (ml *mockedStopsBackend) Close() error {
	// noop
	return nil
}
