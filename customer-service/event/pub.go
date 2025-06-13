package event

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/oliknight1/retail-isa-investment/customer-service/model"
)

type EventPublisher interface {
	PublishCustomer(customer model.Customer) error
}

type NatsPublisher struct {
	nc *nats.Conn
}

func NewNatsPublisher(url string) *NatsPublisher {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	return &NatsPublisher{nc}
}

func (p *NatsPublisher) PublishCustomer(customer model.Customer) error {
	msg, _ := json.Marshal(customer)
	return p.nc.Publish("customer.created", msg)
}
