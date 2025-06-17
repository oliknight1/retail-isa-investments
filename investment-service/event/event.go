package event

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type EventHandler interface {
	Publish(subject string, data any) error
	Close()
}

type NatsPublisher struct {
	conn *nats.Conn
}

func NewNatsPublisher(url string) (EventHandler, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsPublisher{conn: conn}, nil
}

func (p *NatsPublisher) Publish(subject string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal event for subject %s: %v", subject, err)
		return err
	}
	return p.conn.Publish(subject, data)
}

func (p *NatsPublisher) Close() {
	p.conn.Close()
}
