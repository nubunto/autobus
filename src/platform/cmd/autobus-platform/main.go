package main

import (
	"bytes"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"time"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats"
	"github.com/pkg/errors"
)

const stupidDateTimeLayout = "020106150405"

type Message struct {
	MessageHead string
	ID          string
	Type        string
	Valid       bool
	Latitude    float64
	Longitude   float64
	DateTime    time.Time
	Speed       float64
	Direction   int64
	Status      string
}

func main() {
	var horizontalConcurrency int
	concurrencyStr, exists := os.LookupEnv("AUTOBUS_PLATFORM_HORIZONTAL")
	if exists {
		horizontalConcurrency, _ = strconv.Atoi(concurrencyStr)
	} else {
		horizontalConcurrency = 1 << 10
	}

	logger := log.New(os.Stdout, "autobus-platform: ", log.LstdFlags)
	natsURL := os.Getenv("AUTOBUS_PLATFORM_NATS_URL")
	logger.Println("Connecting to nats @", natsURL)
	nc, err := nats.Connect(natsURL)
	if err != nil {
		logger.Fatal("error while connecting to nats:", err)
	}

	dbURL := os.Getenv("AUTOBUS_PLATFORM_PGSQL")
	logger.Println("Connecting to postgres @", dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Fatal("error while connecting to postgres:", err)
	}

	logger.Println("Asynchronously waiting for messages...")
	for i := 0; i < horizontalConcurrency; i++ {
		go nc.QueueSubscribe("gps.update", "queue.pgsql", func(m *nats.Msg) {
			log.Println("Got message:", m.Data, "length:", len(m.Data))
			var parsed Message
			if err := parsed.UnmarshalText(m.Data); err != nil {
				logger.Println("[ERROR] error while parsing the gps message: ", err)
				return
			}
			logger.Println("Inserting in the database... parsed:", parsed)
			if err := parsed.Insert(db); err != nil {
				logger.Println("[ERROR] error while inserting gps data to the database: ", err)
				return
			}
		})
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	logger.Println("shutting autobus-platform down...")
}

func (msg *Message) UnmarshalText(raw []byte) (err error) {
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
	log.Println(parts)
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

	longitude := parts[5]
	msg.Latitude, err = strconv.ParseFloat(longitude, 64)
	if err != nil {
		return errors.Wrapf(err, "error decoding the latitude (raw: %s)", longitude)
	}

	latitudeDir := parts[6]
	if latitudeDir != "S" && latitudeDir != "N" {
		return errors.Errorf("latitude direction should be either S or N (south or north), it is %s", latitudeDir)
	}
	if "S" == latitudeDir {
		msg.Latitude = -msg.Latitude
	}

	latitude := parts[7]
	msg.Longitude, err = strconv.ParseFloat(latitude, 64)
	if err != nil {
		return errors.Wrapf(err, "error decoding the longitude (raw: %s)", latitude)
	}

	longitudeDir := parts[8]
	if longitudeDir != "E" && longitudeDir != "W" {
		return errors.Errorf("longitude direction should be either E or W (east or west), it is %s", longitudeDir)
	}
	if "W" == longitudeDir {
		msg.Longitude = -msg.Longitude
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

func (m *Message) Insert(db *sql.DB) (err error) {
	_, err = db.Exec("INSERT INTO gps_data (id, time, date, latitude, longitude, status) VALUES ($1, $2, $3, $4, $5, $6)",
		m.ID,
		m.DateTime.Format("15:04:05"),
		m.DateTime.Format("02 Jan 2006"),
		m.Latitude,
		m.Longitude,
		m.Status,
	)
	return
}
