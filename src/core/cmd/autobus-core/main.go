package main

import (
	"log"
	"os"
)

func main() {
	hubLogger := log.New(os.Stdout, "hub ", log.LstdFlags)
	np, err := NewNatsProtocol(os.Getenv("AUTOBUS_NATS_URL"))
	if err != nil {
		panic(err)
	}

	h, err := NewHub(hubLogger,
		DebugFromEnv("AUTOBUS_DEBUG"),
		ListenFromEnv("AUTOBUS_TCP_HOST"),
		AcceptGoroutinesFromEnv("AUTOBUS_ACCEPT"),
		HandlerGoroutinesFromEnv("AUTOBUS_HANDLERS"),
		WithProtocol(np),
	)
	if err != nil {
		panic(err)
	}

	if err := h.Start(); err != nil {
		panic(err)
	}
	h.wg.Wait()
}
