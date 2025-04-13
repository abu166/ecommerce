package order

import (
	"ecommerce/internal/config"
	invApp "ecommerce/internal/inventory/application"
	invInfra "ecommerce/internal/inventory/infrastructure"
	ordApp "ecommerce/internal/order/application"
	ordInfra "ecommerce/internal/order/infrastructure"
	"ecommerce/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Run(cfg *config.Config) error {
	// Initialize inventory repository and service
	invRepo, err := invInfra.NewRepository(cfg.DSN())
	if err != nil {
		return err
	}
	invSvc := invApp.NewService(invRepo)

	// Initialize order repository and service
	ordRepo, err := ordInfra.NewRepository(cfg.DSN())
	if err != nil {
		return err
	}
	ordSvc := ordApp.NewService(ordRepo, invSvc)

	// Start gRPC server
	server := NewServer(ordSvc)
	lis, err := net.Listen("tcp", cfg.OrderAddr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	proto.RegisterOrderServiceServer(s, server)
	log.Printf("Order service running on %s", cfg.OrderAddr)
	return s.Serve(lis)
}
