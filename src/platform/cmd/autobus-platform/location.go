package main

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Location struct {
	Type        string
	Coordinates []float64
}

func (l *Location) UnmarshalText(raw []byte) error {
	if len(raw) < 20 {
		return errors.Errorf("insufficient data (raw: %s)", raw)
	}
	rawStr := string(raw)
	parts := strings.Split(rawStr, ",")

	// Latitude

	// Degree
	latitudePart := parts[0]
	latitudeDegree, err := strconv.ParseInt(latitudePart[:2], 10, 64)
	if err != nil {
		return errors.Wrapf(err, "error decoding the latitude degree (raw: %s)", latitudeDegree)
	}

	// Minute
	latitudeMinute, err := strconv.ParseFloat(latitudePart[2:], 64)
	if err != nil {
		return errors.Wrapf(err, "error decoding the latitude minute (raw: %s)", latitudeMinute)
	}

	// Direction
	latitudeDir := parts[1]
	if latitudeDir != "S" && latitudeDir != "N" {
		return errors.Errorf("latitude direction should be either S or N (south or north), it is %s", latitudeDir)
	}
	if "S" == latitudeDir {
		latitudeDegree = -latitudeDegree
	}
	// End Latitude

	// Longitude
	longitudePart := parts[2]

	// Degree
	longitudeDegree, err := strconv.ParseFloat(longitudePart[:3], 64)
	if err != nil {
		return errors.Wrapf(err, "error decoding the longitude (raw: %s)", longitudePart)
	}

	// Minute
	longitudeMinute, err := strconv.ParseFloat(longitudePart[3:], 64)
	if err != nil {
		return errors.Wrapf(err, "error decoding the latitude minute (raw: %s)", latitudePart)
	}

	// Direction
	longitudeDir := parts[3]
	if longitudeDir != "E" && longitudeDir != "W" {
		return errors.Errorf("longitude direction should be either E or W (east or west), it is %s", longitudeDir)
	}
	if "W" == longitudeDir {
		longitudeDegree = -longitudeDegree
	}
	// End Longitude

	// convert the thing
	latitude := float64(latitudeDegree) + (latitudeMinute / 60)
	longitude := float64(longitudeDegree) + (longitudeMinute / 60)

	// TODO: use the minutes as part of the full latitude/longitude information
	// Location is a GeoJSON value
	l.Type = "Point"
	l.Coordinates = []float64{float64(longitude), float64(latitude)}
	return nil
}
