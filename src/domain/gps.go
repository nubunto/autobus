package domain

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
)

const stupidDateTimeLayout = "020106150405"

type GPSMessage struct {
	MessageHead string
	ID          string
	Type        string
	Valid       bool
	Loc         *Location
	DateTime    time.Time
	Speed       float64
	Direction   int64
	Status      string
}

func (msg *GPSMessage) UnmarshalText(raw []byte) (err error) {
	beginning := bytes.Index(raw, []byte("*"))
	end := bytes.LastIndex(raw, []byte("#"))
	if end > len(raw) {
		return errors.Errorf("end of message is impossible to reach. beginning=%d end=%d", beginning, end)
	}
	if beginning == -1 {
		return errors.New("malformed message: no beginning")
	}
	if end == -1 {
		return errors.New("malformed message: no end")
	}

	raw = raw[beginning:end]
	if len(raw) < 60 {
		return errors.New("the raw data has insufficient data")
	}

	rawStr := string(raw)
	parts := strings.Split(rawStr, ",")
	msg.MessageHead = parts[0]
	msg.ID = parts[1]
	msg.Type = parts[2]

	// TODO time
	timePart := parts[3]

	valid := parts[4]
	if "A" == valid {
		msg.Valid = true
	} else if "S" == valid {
		msg.Valid = false
	} else {
		return errors.Errorf("error decoding valid (raw: %s)", valid)
	}

	msg.Loc = new(Location)
	if err := msg.Loc.UnmarshalText(raw[27:51]); err != nil {
		return errors.Wrapf(err, "error while decoding latitude/longitude information")
	}

	speed := parts[9]
	if speed != "" {
		msg.Speed, err = strconv.ParseFloat(speed, 64)
		if err != nil {
			return errors.Wrapf(err, "error decoding speed (raw: %s)", speed)
		}
	}

	direction := parts[10]
	if direction != "" {
		msg.Direction, err = strconv.ParseInt(direction, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error decoding direction (raw: %s)", direction)
		}
	}

	// TODO date
	datePart := parts[11]

	fullDateAndTime := datePart + timePart
	msg.DateTime, err = time.Parse(stupidDateTimeLayout, fullDateAndTime)
	if err != nil {
		return errors.Wrapf(err, "error decoding time (raw: time: %s - date: %s)", timePart, datePart)
	}
	msg.Status = parts[12]
	return nil
}

func (m *GPSMessage) Insert(session *mgo.Session) error {
	transient := session.DB("autobus").C("gps_data_transient")
	persisted := session.DB("autobus").C("gps_data")
	if err := transient.Insert(m); err != nil {
		return errors.Wrap(err, "error while inserting to a transient collection")
	}
	if err := persisted.Insert(m); err != nil {
		return errors.Wrap(err, "error while inserting to a persisted collection")
	}
	return nil
}

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
