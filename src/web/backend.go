package web

import (
	stderrors "errors"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
)

type Backend interface {
	// Lines return the backend for the Lines structs
	// Mainly, a interface that knows how to CRUD lines data
	Lines() LinesBackend

	// Stops return the backend for BusStop structs
	// Mainly, a interface that knows how to CRUD stops data
	Stops() StopsBackend

	// gps data (read-only)
	GPS() GPSBackend

	// gps transient data (read-only)
	GPSTransient() GPSBackend

	// Close releases resources held by the backend.
	Close() error
}

var ErrNotAllowed = stderrors.New("not allowed")

type LinesBackend interface {
	GetAll(finder interface{}) ([]Line, error)
	GetOne(finder interface{}) (*Line, error)
	Create(interface{}) error
	Update(selector, update interface{}) error
	Delete(id interface{}) error
	Close() error
}

type StopsBackend interface {
	GetAll(finder interface{}) ([]BusStop, error)
	GetOne(finder interface{}) (*BusStop, error)
	Create(interface{}) error
	Update(id, fields interface{}) error
	Delete(selector interface{}) error
	Close() error
}

type GPSBackend interface {
	GetAll(finder interface{}) ([]GPSData, error)
	GetOne(finder interface{}) (*GPSData, error)
	// Create is a noop. GPS data comes from somewhere else.
	Create(interface{}) error
	// Update is a noop.
	Update(id, fields interface{}) error
	// Delete is a noop
	Delete(selector interface{}) error

	Close() error
}

func NewMongoBackend(s *mgo.Session) Backend {
	return &mongoBackend{
		Session: s.Copy(),
		LinesBackend: &mongoLinesBackend{
			Session: s.Copy(),
		},
		StopsBackend: &mongoStopsBackend{
			Session: s.Copy(),
		},
		GPSBackend: &mongoGPSBackend{
			Session: s.Copy(),
		},
		GPSTransientBackend: &mongoGPSTransientBackend{
			Session: s.Copy(),
		},
	}
}

type mongoBackend struct {
	LinesBackend
	StopsBackend
	GPSBackend
	GPSTransientBackend GPSBackend
	*mgo.Session
}

type mongoLinesBackend struct {
	*mgo.Session
}

type mongoStopsBackend struct {
	*mgo.Session
}

type mongoGPSBackend struct {
	*mgo.Session
}

type mongoGPSTransientBackend struct {
	*mgo.Session
}

func (mb *mongoBackend) Lines() LinesBackend {
	return mb.LinesBackend
}

func (mb *mongoBackend) Stops() StopsBackend {
	return mb.StopsBackend
}

func (mb *mongoBackend) GPS() GPSBackend {
	return mb.GPSBackend
}

func (mb *mongoBackend) GPSTransient() GPSBackend {
	return mb.GPSTransientBackend
}

func (mb *mongoBackend) Close() error {
	mb.LinesBackend.Close()
	mb.StopsBackend.Close()
	mb.GPSBackend.Close()
	mb.GPSTransientBackend.Close()
	mb.Session.Close()
	return nil
}

func (ml *mongoLinesBackend) GetAll(finder interface{}) ([]Line, error) {
	s := ml.Copy()
	defer s.Close()

	c := s.DB("autobus").C("lines")
	var all []Line
	if err := c.Find(finder).All(&all); err != nil {
		return nil, errors.Wrap(err, "error retrieving list of lines")
	}
	return all, nil
}

func (ml *mongoLinesBackend) GetOne(finder interface{}) (*Line, error) {
	s := ml.Copy()
	defer s.Close()

	c := s.DB("autobus").C("lines")
	var one Line
	if err := c.Find(finder).One(&one); err != nil {
		return nil, errors.Wrap(err, "error retrieving single line")
	}
	return &one, nil
}

func (ml *mongoLinesBackend) Create(doc interface{}) error {
	s := ml.Copy()
	defer s.Close()

	c := s.DB("autobus").C("lines")
	if err := c.Insert(doc); err != nil {
		return errors.Wrap(err, "error creating new line")
	}
	return nil
}

func (ml *mongoLinesBackend) Update(selector, update interface{}) error {
	s := ml.Copy()
	defer s.Close()

	c := s.DB("autobus").C("lines")
	if err := c.Update(selector, update); err != nil {
		return errors.Wrap(err, "error updating line")
	}
	return nil
}

func (ml *mongoLinesBackend) Delete(selector interface{}) error {
	s := ml.Copy()
	defer s.Close()

	c := s.DB("autobus").C("lines")
	if err := c.Remove(selector); err != nil {
		return errors.Wrap(err, "error removing line")
	}
	return nil
}

func (ml *mongoLinesBackend) Close() error {
	ml.Session.Close()
	return nil
}

func (ms *mongoStopsBackend) GetAll(finder interface{}) ([]BusStop, error) {
	s := ms.Copy()
	defer s.Close()

	var all []BusStop
	c := s.DB("autobus").C("stops")
	if err := c.Find(finder).All(&all); err != nil {
		return nil, errors.Wrap(err, "error retrieving list of stops")
	}
	return all, nil
}

func (ms *mongoStopsBackend) GetOne(selector interface{}) (*BusStop, error) {
	s := ms.Copy()
	defer s.Close()

	var one BusStop
	c := s.DB("autobus").C("stops")
	if err := c.Find(selector).One(&one); err != nil {
		return nil, errors.Wrap(err, "error retrieving single stop")
	}
	return &one, nil
}

func (ms *mongoStopsBackend) Create(doc interface{}) error {
	s := ms.Copy()
	defer s.Close()

	c := s.DB("autobus").C("stops")
	if err := c.Insert(doc); err != nil {
		return errors.Wrap(err, "error creating stop")
	}
	return nil
}

func (ms *mongoStopsBackend) Update(selector, update interface{}) error {
	s := ms.Copy()
	defer s.Close()

	c := s.DB("autobus").C("stops")
	if err := c.Update(selector, update); err != nil {
		return errors.Wrap(err, "error updating stop")
	}
	return nil
}

func (ms *mongoStopsBackend) Delete(selector interface{}) error {
	s := ms.Copy()
	defer s.Close()

	c := s.DB("autobus").C("stops")
	if err := c.Remove(selector); err != nil {
		return errors.Wrap(err, "error removing stop")
	}
	return nil
}

func (ms *mongoStopsBackend) Close() error {
	ms.Session.Close()
	return nil
}

func (mg *mongoGPSBackend) GetAll(selector interface{}) ([]GPSData, error) {
	s := mg.Copy()
	defer s.Close()

	c := s.DB("autobus").C("gps_data")
	var all []GPSData
	if err := c.Find(selector).All(&all); err != nil {
		return nil, errors.Wrap(err, "error retrieving gps data")
	}
	return all, nil
}

func (mg *mongoGPSBackend) GetOne(selector interface{}) (*GPSData, error) {
	return nil, ErrNotAllowed
}

func (mg *mongoGPSBackend) Create(selector interface{}) error {
	return ErrNotAllowed
}

func (mg *mongoGPSBackend) Update(selector, update interface{}) error {
	return ErrNotAllowed
}

func (mg *mongoGPSBackend) Delete(selector interface{}) error {
	return ErrNotAllowed
}

func (mg *mongoGPSBackend) Close() error {
	mg.Session.Close()
	return nil
}

func (mg *mongoGPSTransientBackend) GetAll(selector interface{}) ([]GPSData, error) {
	s := mg.Copy()
	defer s.Close()

	c := s.DB("autobus").C("gps_data_transient")
	var all []GPSData
	if err := c.Find(selector).All(&all); err != nil {
		return nil, errors.Wrap(err, "error retrieving transient gps data")
	}
	return all, nil
}

func (mg *mongoGPSTransientBackend) GetOne(selector interface{}) (*GPSData, error) {
	return nil, ErrNotAllowed
}

func (mg *mongoGPSTransientBackend) Create(selector interface{}) error {
	return ErrNotAllowed
}

func (mg *mongoGPSTransientBackend) Update(selector, update interface{}) error {
	return ErrNotAllowed
}

func (mg *mongoGPSTransientBackend) Delete(selector interface{}) error {
	return ErrNotAllowed
}

func (mg *mongoGPSTransientBackend) Close() error {
	mg.Session.Close()
	return nil
}
