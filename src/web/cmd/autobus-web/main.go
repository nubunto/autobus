package main

import (
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux"
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

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()

	dbURL := os.Getenv("AUTOBUS_WEB_MONGO_URL")

	sugar.Infow("connecting to db", "address", dbURL)
	env, err := NewEnv(
		DialDB(dbURL),
		SetLogger(logger),
	)
	if err != nil {
		panic(err)
	}

	env.ensureIndex("stops", "$2dsphere:location")

	sugar.Infow("mounting routes")
	mux := httptreemux.New()
	mux.GET("/gps", handleGetGPS(env))
	mux.POST("/stops", handleCreateStop(env))
	mux.GET("/stops", handleGetStops(env))
	sugar.Infow("done mounting routes")

	listenAddr := os.Getenv("AUTOBUS_WEB_LISTEN_ADDR")
	sugar.Infow("preparing to listen", "address", listenAddr)
	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		panic(err)
	}
}
