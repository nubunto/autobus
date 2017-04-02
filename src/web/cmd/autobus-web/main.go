package main

import (
	"os"
	"web/app"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Create service
	service := goa.New("autobus-web")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	env, err := NewEnv(
		DialDB(os.Getenv("AUTOBUS_WEB_MONGO_URL")),
	)
	if err != nil {
		service.LogError("startup", "err", err)
		os.Exit(1)
	}
	env.ensure2dsphereIndex("stops")

	// Mount "GPS" controller
	gpsController := NewGPSController(service, env)
	app.MountGPSController(service, gpsController)

	busStopController := NewBusStopController(service, env)
	app.MountBusStopController(service, busStopController)

	swaggerController := NewSwaggerController(service)
	app.MountSwaggerController(service, swaggerController)

	listenAddr := os.Getenv("AUTOBUS_WEB_LISTEN_ADDR")
	// Start service
	if err := service.ListenAndServe(listenAddr); err != nil {
		service.LogError("startup", "err", err)
	}
}
