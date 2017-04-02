package design

import (
	"os"

	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("autobus-web", func() {
	Description("The web platform for the Autobus Tracker")
	Scheme("http")
	Host(os.Getenv("AUTOBUS_WEB_HOST"))
})

var _ = Resource("GPS", func() {
	BasePath("/gps")
	Action("show", func() {
		Routing(GET("/"))
		Response(OK, ArrayOf(GPSMedia))
	})
})

var _ = Resource("BusStop", func() {
	BasePath("/stops")
	Action("create", func() {
		Routing(POST("/"))
		Payload(BusStopPayload)
		Response(Created, BusStopMedia)
	})
	Action("nearest", func() {
		Routing(GET("/"))
		Params(func() {
			Param("longitude", Number)
			Param("latitude", Number)
			Param("radius", Number, "distance in meters")
			Required("longitude", "latitude", "radius")
		})
		Response(OK, ArrayOf(BusStopMedia))
	})
})

var _ = Resource("swagger", func() {
	Origin("*", func() {
		Methods("GET")
	})
	Files("/swagger.json", "swagger/swagger.json")
})

var BusStop = Type("BusStop", func() {
	Description("The stops where the bus passes through")
	Attribute("id", Integer)
	Attribute("name", String)
	Attribute("loc", GeoJSON)
})

var BusStopPayload = Type("BusStopPayload", func() {
	Attribute("latitude", Number)
	Attribute("longitude", Number)
	Attribute("name", String)

	Required("name", "latitude", "longitude")
})

var GPS = Type("gps", func() {
	Description("The GPS data received by the core Autobus application")
	Attribute("head", String)
	Attribute("id", String)
	Attribute("type", String)
	Attribute("valid", Boolean)
	Attribute("dateTime", DateTime)
	Attribute("latitude", Number)
	Attribute("longitude", Number)
	Attribute("speed", Number)
	Attribute("direction", Integer)
	Attribute("status", String)
})

var GeoJSON = Type("GeoJSON", func() {
	Attribute("type", String)
	Attribute("coordinates", ArrayOf(Number))
})

var GPSMedia = MediaType("autobus.web.platform/gps.media+json", func() {
	Description("The GPS data media type")
	Reference(GPS)
	Attributes(func() {
		Attribute("head")
		Attribute("id")
		Attribute("type")
		Attribute("valid")
		Attribute("dateTime")
		Attribute("latitude")
		Attribute("longitude")
		Attribute("speed")
		Attribute("direction")
		Attribute("status")
	})
	View("default", func() {
		Attribute("head")
		Attribute("id")
		Attribute("type")
		Attribute("valid")
		Attribute("dateTime")
		Attribute("latitude")
		Attribute("longitude")
		Attribute("speed")
		Attribute("direction")
		Attribute("status")
	})
})
var BusStopMedia = MediaType("autobus.web.platform/bus-stop.media+json", func() {
	Description("The Bus Stop media type")
	Reference(BusStop)
	Attributes(func() {
		Attribute("id")
		Attribute("name")
		Attribute("loc")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("loc")
	})
})
