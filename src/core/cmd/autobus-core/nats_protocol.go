package main

import "github.com/nats-io/nats"

type NatsProtocol struct {
	client *nats.Conn
}

const (
	SubjectMessageReceived string = "gps.update"
)

func (np *NatsProtocol) HandleMessage(msg []byte) ([]byte, error) {
	if err := np.client.Publish(SubjectMessageReceived, msg); err != nil {
		return nil, err
	}
	return nil, nil
}

func NewNatsProtocol(urls string) (Protocol, error) {
	nc, err := nats.Connect(urls)
	if err != nil {
		return nil, err
	}
	return &NatsProtocol{
		client: nc,
	}, nil
}
