package main

import (
	"log"
	"time"
)

// Protocol is a interface for messages of clients.
type Protocol interface {
	HandleMessage(msg []byte) ([]byte, error)
}

// ProtocolFunc is a Protocol implementation as a closure
type ProtocolFunc func(msg []byte) ([]byte, error)

// HandleMessage implements Protocol for ProtocolFunc
func (pf ProtocolFunc) HandleMessage(msg []byte) ([]byte, error) {
	return pf(msg)
}

type Decorator func(Protocol) Protocol

func Logging(logger *log.Logger) Decorator {
	return func(p Protocol) Protocol {
		return ProtocolFunc(func(msg []byte) ([]byte, error) {
			start := time.Now()
			ret, err := p.HandleMessage(msg)
			logger.Println("Message:", msg)
			logger.Println("Response:", ret)
			logger.Println("Took:", time.Since(start))
			return ret, err
		})
	}
}

func Decorate(root Protocol, decorators ...Decorator) Protocol {
	decorated := root
	for _, d := range decorators {
		decorated = d(decorated)
	}
	return decorated
}
