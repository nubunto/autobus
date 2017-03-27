package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats"
)

type Message struct {
	ID        string
	Date      time.Time
	Latitude  string
	Longitude string
	Status    string
}

func main() {
	var horizontalConcurrency int
	concurrencyStr, exists := os.LookupEnv("AUTOBUS_PLATFORM_HORIZONTAL")
	if exists {
		horizontalConcurrency, _ = strconv.Atoi(concurrencyStr)
	} else {
		horizontalConcurrency = 1 << 10
	}

	fmt.Println("Connecting to nats...")
	nc, err := nats.Connect(os.Getenv("AUTOBUS_PLATFORM_NATS_URL"))
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", os.Getenv("AUTOBUS_PLATFORM_PGSQL"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Asynchronously waiting for messages...")
	for i := 0; i < horizontalConcurrency; i++ {
		go nc.QueueSubscribe("gps.update", "queue.pgsql", func(m *nats.Msg) {
			strData := string(m.Data)
			fmt.Println("Got message:", strData, len(strData))
			parsed, err := parseMessage(m.Data)
			if err != nil {
				// bail
				fmt.Println("error:", err)
				return
			}
			fmt.Println("Inserting in the database...")
			if _, err := db.Exec("INSERT INTO gps_data (id, time, date, latitude, longitude, status) VALUES ($1, $2, $3, $4, $5, $6)",
				parsed.ID,
				parsed.Date.Format("15:04:05"),
				parsed.Date.Format("02 Jan 2006"),
				parsed.Latitude,
				parsed.Longitude,
				parsed.Status,
			); err != nil {
				fmt.Println("got error:", err)
			}
		})
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}

func parseMessage(raw []byte) (msg Message, err error) {
	beginning := bytes.Index(raw, []byte("*"))
	end := bytes.LastIndex(raw, []byte("#"))
	if beginning == -1 {
		return Message{}, fmt.Errorf("malformed message: no beginning")
	}
	if end == -1 {
		return Message{}, fmt.Errorf("malformed message: no end")
	}
	fmt.Println("beginning:", beginning, "end:", end)
	raw = raw[beginning:end]
	if len(raw) < 60 {
		return Message{}, fmt.Errorf("the raw data has insufficient data")
	}
	rawStr := string(raw)
	parts := strings.Split(rawStr, ",")
	msg.ID = parts[1]
	//time := parts[3]
	msg.Latitude = parts[5]
	// TODO: South or North (parts[6] == N or S)
	msg.Longitude = parts[7]
	// TODO: East or West (parts[8] == E or W)
	//date := parts[11]
	msg.Status = parts[12]
	return msg, nil
}
