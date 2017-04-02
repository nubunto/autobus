package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
)

const stupidDateTimeLayout = "020106150405"

func createCappedCollection(session *mgo.Session) error {
	transient := session.DB("autobus").C("gps_data_transient")
	if err := transient.Create(&mgo.CollectionInfo{
		Capped:   true,
		MaxBytes: 1 << 10, // 1 KB
		MaxDocs:  500,
	}); err != nil {
		return errors.Wrap(err, "error creating transient collection")
	}
	return nil
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

	dbURL := os.Getenv("AUTOBUS_PLATFORM_MONGO_URL")
	logger.Println("Connecting to db @", dbURL)

	session, err := mgo.Dial(dbURL)
	if err != nil {
		logger.Fatal("error while connecting to db:", err)
	}
	defer session.Close()

	// create the transient capped collection
	if err := createCappedCollection(session); err != nil {
		logger.Println("[WARN] error while creating collections:", err)
	}

	logger.Println("Asynchronously waiting for messages...")
	for i := 0; i < horizontalConcurrency; i++ {
		go nc.QueueSubscribe("gps.update", "queue.web.database", func(m *nats.Msg) {
			log.Println("Got message:", m.Data, "length:", len(m.Data))

			var parsed Message
			if err := parsed.UnmarshalText(m.Data); err != nil {
				logger.Println("[ERROR] error while parsing the gps message: ", err)
				return
			}

			logger.Println("Inserting in the database... parsed:", parsed)
			if err := parsed.Insert(session); err != nil {
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
