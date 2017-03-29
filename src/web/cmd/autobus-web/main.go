package main

import (
	"database/sql"
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

	// Mount "GPS" controller
	db, err := sql.Open("postgres", os.Getenv("AUTOBUS_WEB_PGSQL"))
	if err != nil {
		panic(err)
	}

	c := NewGPSController(service, db)
	app.MountGPSController(service, c)

	autobusHost := os.Getenv("AUTOBUS_WEB_HOST")
	// Start service
	if err := service.ListenAndServe(autobusHost); err != nil {
		service.LogError("startup", "err", err)
	}
}
