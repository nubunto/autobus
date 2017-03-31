// Code generated by goagen v1.1.0-dirty, command line:
// $ goagen
// --design=web/design
// --out=$(GOPATH)/src/web
// --version=v1.1.0-dirty
//
// API "autobus-web": Application Media Types
//
// The content of this file is auto-generated, DO NOT MODIFY

package app

import (
	"time"
)

// The GPS data media type (default view)
//
// Identifier: autobus.web.platform/gps.media+json; view=default
type GpsMedia struct {
	DateTime  *time.Time `form:"dateTime,omitempty" json:"dateTime,omitempty" xml:"dateTime,omitempty"`
	Direction *int       `form:"direction,omitempty" json:"direction,omitempty" xml:"direction,omitempty"`
	Head      *string    `form:"head,omitempty" json:"head,omitempty" xml:"head,omitempty"`
	ID        *string    `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	Latitude  *float64   `form:"latitude,omitempty" json:"latitude,omitempty" xml:"latitude,omitempty"`
	Longitude *float64   `form:"longitude,omitempty" json:"longitude,omitempty" xml:"longitude,omitempty"`
	Speed     *float64   `form:"speed,omitempty" json:"speed,omitempty" xml:"speed,omitempty"`
	Status    *string    `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
	Type      *string    `form:"type,omitempty" json:"type,omitempty" xml:"type,omitempty"`
	Valid     *bool      `form:"valid,omitempty" json:"valid,omitempty" xml:"valid,omitempty"`
}
