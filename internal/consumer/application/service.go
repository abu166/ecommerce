package application

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"ecommerce/proto"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type Service struct {
	nc        *nats.Conn
	invClient proto.InventoryServiceClient
}

func NewService(nc *nats.Conn, invClient proto.InventoryServiceClient) *Service {
	return &Service{nc: nc, invClient: invClient}
}

func (s *Service) SubscribeToOrders() error {
	// Retry subscription with exponential backoff
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		_, err := s.nc.Subscribe("order.created", func(msg *nats.Msg) {
			var order proto.OrderResponse
			if err := json.Unmarshal(msg.Data, &order); err != nil {
				logrus.Errorf("Failed to unmarshal order: %v", err)
				return
			}

			logrus.Infof("Received order.created event for order %s", order.Id)
			for _, item := range order.Items {
				var updateErr error
				// Retry stock update
				for retry := 0; retry < 3; retry++ {
					_, updateErr = s.invClient.UpdateProduct(context.Background(), &proto.UpdateProductRequest{
						Id:    item.ProductId,
						Stock: -item.Quantity, // Decrease stock
					})
					if updateErr == nil {
						logrus.Infof("Updated stock for product %s by -%d", item.ProductId, item.Quantity)
						break
					}
					logrus.Errorf("Retry %d: Failed to update stock for product %s: %v", retry+1, item.ProductId, updateErr)
					time.Sleep(time.Duration(retry*100) * time.Millisecond)
				}
				if updateErr != nil {
					logrus.Errorf("Failed to update stock for product %s after retries", item.ProductId)
				}
			}
		})
		if err == nil {
			logrus.Info("Successfully subscribed to order.created events")
			return nil
		}
		logrus.Errorf("Retry %d: Failed to subscribe to order.created: %v", i+1, err)
		time.Sleep(time.Duration(i*100) * time.Millisecond)
	}
	return ErrSubscriptionFailed
}

var ErrSubscriptionFailed = errors.New("failed to subscribe to NATS after retries")
