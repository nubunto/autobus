package main

import (
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

var Version string

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()

	sugar.Infow("autobus-web", "version", Version)
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

	mux := httprouter.New()
	mux.GET("/gps", handleGetGPS(env))

	mux.POST("/stops", handleCreateStop(env))
	mux.GET("/stops", handleGetStops(env))

	mux.GET("/lines", handleGetLines(env))
	mux.POST("/lines", handleCreateLine(env))

	mux.GET("/version", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte(Version))
	})

	listenAddr := os.Getenv("AUTOBUS_WEB_LISTEN_ADDR")
	sugar.Infow("preparing to listen", "address", listenAddr)
	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		panic(err)
	}
}
