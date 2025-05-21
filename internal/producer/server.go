package producer

import (
	"context"
	"ecommerce/internal/config"
	"ecommerce/internal/producer/application"
	"ecommerce/proto"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	proto.UnimplementedProducerServiceServer
	svc *application.Service
}

func NewServer(svc *application.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) NotifyOrderCreated(ctx context.Context, req *proto.OrderResponse) (*proto.ProducerEmpty, error) {
	err := s.svc.NotifyOrderCreated(ctx, req)
	if err != nil {
		return nil, err
	}
	return &proto.ProducerEmpty{}, nil
}

func Run(cfg *config.Config) error {
	nc, err := nats.Connect(cfg.NATSAddr)
	if err != nil {
		return err
	}
	defer nc.Close()

	svc := application.NewService(nc)
	server := NewServer(svc)

	lis, err := net.Listen("tcp", cfg.ProducerAddr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	proto.RegisterProducerServiceServer(s, server)
	log.Printf("Producer service running on %s", cfg.ProducerAddr)
	return s.Serve(lis)
}
