package order

import (
	"context"
	"ecommerce/internal/order/application"
	"ecommerce/internal/order/domain"
	"ecommerce/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.UnimplementedOrderServiceServer
	svc *application.Service
}

func NewServer(svc *application.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.OrderResponse, error) {
	o := &domain.Order{
		UserID: req.UserId,
	}
	for _, item := range req.Items {
		o.Items = append(o.Items, domain.OrderItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
		})
	}
	if err := s.svc.Create(ctx, o); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create order: %v", err)
	}
	return toResponse(o), nil
}

func (s *Server) GetOrder(ctx context.Context, req *proto.GetOrderRequest) (*proto.OrderResponse, error) {
	o, err := s.svc.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}
	return toResponse(o), nil
}

func (s *Server) UpdateOrder(ctx context.Context, req *proto.UpdateOrderRequest) (*proto.OrderResponse, error) {
	o, err := s.svc.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}
	o.Status = req.Status
	if err := s.svc.Update(ctx, o); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order: %v", err)
	}
	return toResponse(o), nil
}

func (s *Server) ListOrders(ctx context.Context, req *proto.ListOrdersRequest) (*proto.ListOrdersResponse, error) {
	orders, total, err := s.svc.List(ctx, req.UserId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}
	var resp []*proto.OrderResponse
	for _, o := range orders {
		resp = append(resp, toResponse(o))
	}
	return &proto.ListOrdersResponse{
		Orders: resp,
		Total:  int32(total),
	}, nil
}

func toResponse(o *domain.Order) *proto.OrderResponse {
	var items []*proto.OrderItem
	for _, item := range o.Items {
		items = append(items, &proto.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
	}
	return &proto.OrderResponse{
		Id:     o.ID,
		UserId: o.UserID,
		Items:  items,
		Status: o.Status,
		Total:  o.Total,
	}
}
