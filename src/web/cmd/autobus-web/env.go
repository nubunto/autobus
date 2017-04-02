package main

import (
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
)

type Env struct {
	*mgo.Session
}

type EnvOption func(*Env) error

func NewEnv(options ...EnvOption) (*Env, error) {
	e := new(Env)
	for _, opt := range options {
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func DialDB(db string) EnvOption {
	return func(e *Env) error {
		if db == "" {
			return errors.New("DB address can't be blank")
		}
		session, err := mgo.Dial(db)
		if err != nil {
			return errors.Wrap(err, "can't reach mongodb @ "+db)
		}
		e.Session = session
		return nil
	}
}

func (e *Env) ensure2dsphereIndex(collection string) {
	index := mgo.Index{
		Key: []string{"$2dsphere:loc"},
	}
	c := e.Session.DB("autobus").C(collection)
	c.EnsureIndex(index)
}
