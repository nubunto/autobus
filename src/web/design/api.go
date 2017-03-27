package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("autobus-web", func() {
	Description("The web platform for the Autobus Tracker")
	Scheme("http")
	Host("localhost:8080")
})

var _ = Resource("GPS", func() {
	BasePath("/gps")
	Action("show", func() {
		Routing(GET("/"))
		Response(OK, ArrayOf(GPSMedia))
	})
})

var GPS = Type("gps", func() {
	Description("The GPS data received by the core Autobus application")
	Attribute("id", String)
	Attribute("time", String)
	Attribute("date", String)
	Attribute("longitude", String)
	Attribute("latitude", String)
	Attribute("status", String)
})

var GPSMedia = MediaType("autobus.web.platform/gps.media+json", func() {
	Description("The GPS data media type")
	Reference(GPS)
	Attributes(func() {
		Attribute("id")
		Attribute("time")
		Attribute("date")
		Attribute("longitude")
		Attribute("latitude")
		Attribute("status")
	})
	View("default", func() {
		Attribute("id")
		Attribute("time")
		Attribute("date")
		Attribute("longitude")
		Attribute("latitude")
		Attribute("status")
	})
})
