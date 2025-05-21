package user

import (
	"ecommerce/internal/config"
	"ecommerce/internal/user/application"
	"ecommerce/internal/user/infrastructure"
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

	lis, err := net.Listen("tcp", cfg.UserAddr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, server)
	log.Printf("User service running on %s", cfg.UserAddr)
	return s.Serve(lis)
}
