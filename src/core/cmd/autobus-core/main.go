package main

import (
	"log"
	"os"
)

var Version string

func main() {
	hubLogger := log.New(os.Stdout, "hub ", log.LstdFlags)
	np, err := NewNatsProtocol(os.Getenv("AUTOBUS_CORE_NATS_URL"))
	if err != nil {
		panic(err)
	}

	hubLogger.Println("Version:", Version)

	h, err := NewHub(hubLogger,
		DebugFromEnv("AUTOBUS_CORE_DEBUG"),
		ListenFromEnv("AUTOBUS_CORE_TCP_HOST"),
		AcceptGoroutinesFromEnv("AUTOBUS_CORE_ACCEPT"),
		HandlerGoroutinesFromEnv("AUTOBUS_CORE_HANDLERS"),
		WithProtocol(np),
	)
	if err != nil {
		panic(err)
	}

	if err := h.Start(); err != nil {
		panic(err)
	}
	h.Wait()
}
