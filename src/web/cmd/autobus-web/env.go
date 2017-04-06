package main

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
)

type Env struct {
	*mgo.Session
	*zap.Logger
}

func NewEnv(options ...func(*Env) error) (*Env, error) {
	e := new(Env)
	for _, opt := range options {
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func DialDB(db string) func(*Env) error {
	return func(e *Env) error {
		if db == "" {
			return errors.New("DB address can't be blank")
		}

		session, err := mgo.Dial(db)
		if err != nil {
			return errors.Wrap(err, "error when dialing db")
		}
		e.Session = session
		return nil
	}
}

func SetLogger(l *zap.Logger) func(*Env) error {
	return func(e *Env) error {
		e.Logger = l
		return nil
	}
}

func (e *Env) ensureIndex(collection, index string) {
	e.DB("autobus").C(collection).EnsureIndexKey(index)
}
