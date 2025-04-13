package inventory

import (
	"ecommerce/internal/config"
	"ecommerce/internal/inventory/application"
	"ecommerce/internal/inventory/infrastructure"
	"ecommerce/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Run(cfg *config.Config) error {
	repo, err := infrastructure.NewRepository(cfg.DSN())
	if err != nil {
		return err
	}
	svc := application.NewService(repo)
	server := NewServer(svc)

	lis, err := net.Listen("tcp", cfg.InventoryAddr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	proto.RegisterInventoryServiceServer(s, server)
	log.Printf("Inventory service running on %s", cfg.InventoryAddr)
	return s.Serve(lis)
}
