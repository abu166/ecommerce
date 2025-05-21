package application

import (
	"context"
	"ecommerce/proto"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
)

type Service struct {
	nc *nats.Conn
}

func NewService(nc *nats.Conn) *Service {
	return &Service{nc: nc}
}

func (s *Service) NotifyOrderCreated(ctx context.Context, order *proto.OrderResponse) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = s.nc.Publish("order.created", data)
	if err != nil {
		return err
	}

	log.Printf("Published order.created event for order ID: %s", order.Id)
	return nil
}
