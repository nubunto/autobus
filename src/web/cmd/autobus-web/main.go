package main

import (
	"os"
	"web/app"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	_ "github.com/lib/pq"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	// Create service
	service := goa.New("autobus-web")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "GPS" controller
	session, err := mgo.Dial(os.Getenv("AUTOBUS_WEB_MONGO_URL"))
	if err != nil {
		panic(err)
	}

	c := NewGPSController(service, session)
	app.MountGPSController(service, c)

	autobusHost := os.Getenv("AUTOBUS_WEB_HOST")
	// Start service
	if err := service.ListenAndServe(autobusHost); err != nil {
		service.LogError("startup", "err", err)
	}
}
