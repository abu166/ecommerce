package order

import (
	"ecommerce/internal/config"
	"ecommerce/internal/order/application"
	"ecommerce/internal/order/infrastructure"
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
	cache := infrastructure.NewRedisCache(cfg.RedisAddr)
	svc := application.NewService(repo, cache)
	server := NewServer(svc)

	lis, err := net.Listen("tcp", cfg.OrderAddr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	proto.RegisterOrderServiceServer(s, server)
	log.Printf("Order service running on %s", cfg.OrderAddr)
	return s.Serve(lis)
}
