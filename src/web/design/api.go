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
