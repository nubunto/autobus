package main

import (
	"go.uber.org/zap"
	"web"
)

type Env struct {
	web.Backend
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

func SetBackend(backend web.Backend) func(*Env) error {
	return func(e *Env) error {
		e.Backend = backend
		return nil
	}
}

func SetLogger(l *zap.Logger) func(*Env) error {
	return func(e *Env) error {
		e.Logger = l
		return nil
	}
}
