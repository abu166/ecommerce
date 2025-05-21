package consumer

import (
	"ecommerce/internal/config"
	"ecommerce/internal/consumer/application"
	"ecommerce/proto"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func Run(cfg *config.Config) error {
	nc, err := nats.Connect(cfg.NATSAddr)
	if err != nil {
		return err
	}
	defer nc.Close()

	// Connect to inventory-service using the configured address directly
	invConn, err := grpc.Dial(cfg.InventoryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer invConn.Close()

	invClient := proto.NewInventoryServiceClient(invConn)
	svc := application.NewService(nc, invClient)

	if err := svc.SubscribeToOrders(); err != nil {
		return err
	}

	log.Printf("Consumer service subscribed to order.created events")
	select {} // Keep running
}
