package apigateway

import (
	"ecommerce/internal/config"
	"ecommerce/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Server struct {
	invClient proto.InventoryServiceClient
	ordClient proto.OrderServiceClient
	usrClient proto.UserServiceClient
}

func NewServer(cfg *config.Config) (*Server, error) {
	invConn, err := grpc.Dial("localhost"+cfg.InventoryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	ordConn, err := grpc.Dial("localhost"+cfg.OrderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	usrConn, err := grpc.Dial("localhost"+cfg.UserAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Server{
		invClient: proto.NewInventoryServiceClient(invConn),
		ordClient: proto.NewOrderServiceClient(ordConn),
		usrClient: proto.NewUserServiceClient(usrConn),
	}, nil
}

func Run(cfg *config.Config) error {
	srv, err := NewServer(cfg)
	if err != nil {
		return err
	}

	r := gin.Default()
	srv.SetupRoutes(r)

	log.Printf("API Gateway running on %s", cfg.APIGatewayAddr)
	return r.Run(cfg.APIGatewayAddr)
}
