package main

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

const (
	acceptGoroutinesDefault  = 1024
	handlerGoroutinesDefault = 2048
)

type hub struct {
	*log.Logger
	listener net.Listener
	conns    chan net.Conn
	err      chan error
	debug    bool
	*sync.WaitGroup

	addr                                string
	acceptGoroutines, handlerGoroutines int
	Protocol
}

type hubOption func(*hub) error

func NewHub(logger *log.Logger, options ...hubOption) (*hub, error) {
	h := &hub{
		Logger: logger,
		err:    make(chan error),
		conns:  make(chan net.Conn),
	}
	for _, opt := range options {
		if err := opt(h); err != nil {
			return nil, err
		}
	}
	return h, nil
}

func ListenOn(addr string) hubOption {
	return func(h *hub) (err error) {
		if addr == "" {
			h.Println("[WARNING] the connection hub will start at the default port (9009), which may not be what you expect.")
			addr = "0.0.0.0:9009"
		}
		h.addr = addr
		return nil
	}
}

func ListenFromEnv(envVar string) hubOption {
	return ListenOn(os.Getenv(envVar))
}

func Debug(debug bool) hubOption {
	return func(h *hub) error {
		if debug {
			h.Println("[WARNING] Be aware that debug can slow things down.")
		}
		h.debug = debug
		return nil
	}
}

func DebugFromEnv(envVar string) hubOption {
	b, err := strconv.ParseBool(os.Getenv(envVar))
	if err != nil {
		b = false
	}
	return Debug(b)
}

func parseIntFromEnv(env string, defaultValue int) (r int, err error) {
	n, exists := os.LookupEnv(env)
	if !exists {
		return defaultValue, nil
	}
	r, err = strconv.Atoi(n)
	return
}

func AcceptGoroutines(count int) hubOption {
	return func(h *hub) error {
		h.acceptGoroutines = count
		return nil
	}
}

func HandlerGoroutines(count int) hubOption {
	return func(h *hub) error {
		h.handlerGoroutines = count
		return nil
	}
}

func AcceptGoroutinesFromEnv(env string) hubOption {
	count, err := parseIntFromEnv(env, acceptGoroutinesDefault)
	if err != nil {
		panic(err)
	}
	return AcceptGoroutines(count)
}

func HandlerGoroutinesFromEnv(env string) hubOption {
	count, err := parseIntFromEnv(env, handlerGoroutinesDefault)
	if err != nil {
		panic(err)
	}
	return HandlerGoroutines(count)
}

func WithProtocol(p Protocol) hubOption {
	return func(h *hub) error {
		h.Protocol = p
		return nil
	}
}

func (h *hub) Start() error {
	h.Println("Starting connection hub @", h.addr)
	ln, err := net.Listen("tcp", h.addr)
	if err != nil {
		return err
	}
	h.listener = ln

	h.WaitGroup = new(sync.WaitGroup)
	// accept + handlers + the intercept goroutine
	h.WaitGroup.Add(h.acceptGoroutines + h.handlerGoroutines + 1)

	for i := 0; i < h.acceptGoroutines; i++ {
		go h.accept()
	}

	for i := 0; i < h.handlerGoroutines; i++ {
		go h.openHandlers()
	}

	go h.interceptErrors()
	return nil
}

func (h *hub) accept() {
	for {
		conn, err := h.listener.Accept()
		if err != nil {
			h.err <- err
			break
		}
		h.conns <- conn
	}
	h.WaitGroup.Done()
}

func (h *hub) openHandlers() {
	for conn := range h.conns {
		for {
			msg := make([]byte, 256)
			var n int
			var err error
			if n, err = conn.Read(msg); err != nil {
				h.err <- err
				conn.Close()
				break
			}
			ret, err := h.Protocol.HandleMessage(msg[:n])
			if err != nil {
				h.logDebug("Dropping this message. Reason:", err)
				continue
			}
			if ret == nil {
				// if we return a nil buffer,
				// don't even bother.
				continue
			}
			if _, err := conn.Write(ret); err != nil {
				h.err <- err
				conn.Close()
				break
			}
		}
	}
	h.WaitGroup.Done()
}

func (h *hub) interceptErrors() {
	for err := range h.err {
		if err != io.EOF {
			h.Println("Got error:", err)
		}
	}
	h.WaitGroup.Done()
}

func (h *hub) logDebug(args ...interface{}) {
	if h.debug {
		h.Println(args...)
	}
}
