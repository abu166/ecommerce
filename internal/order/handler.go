package order

import (
	"context"
	"ecommerce/internal/order/application"
	"ecommerce/internal/order/domain"
	"ecommerce/proto"
	"github.com/google/uuid"
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
	if req.UserId == "" || len(req.Items) == 0 || req.Total <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID, items, and total are required")
	}
	for _, item := range req.Items {
		if item.ProductId == "" || item.Quantity <= 0 {
			return nil, status.Error(codes.InvalidArgument, "invalid order item: product ID and quantity are required")
		}
	}
	items := make([]domain.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.OrderItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
		}
	}
	o := &domain.Order{
		ID:     uuid.New().String(),
		UserID: req.UserId,
		Items:  items,
		Status: "pending",
		Total:  req.Total,
	}
	if err := s.svc.Create(ctx, o); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}
	resp := &proto.OrderResponse{
		Id:     o.ID,
		UserId: o.UserID,
		Status: o.Status,
		Total:  o.Total,
	}
	for _, item := range o.Items {
		resp.Items = append(resp.Items, &proto.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
	}
	return resp, nil
}

func (s *Server) GetOrder(ctx context.Context, req *proto.GetOrderRequest) (*proto.OrderResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}
	o, err := s.svc.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get order: %v", err)
	}
	resp := &proto.OrderResponse{
		Id:     o.ID,
		UserId: o.UserID,
		Status: o.Status,
		Total:  o.Total,
	}
	for _, item := range o.Items {
		resp.Items = append(resp.Items, &proto.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
	}
	return resp, nil
}

func (s *Server) UpdateOrder(ctx context.Context, req *proto.UpdateOrderRequest) (*proto.OrderResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}
	if req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "status is required for update")
	}
	o, err := s.svc.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get order: %v", err)
	}
	o.Status = req.Status
	if err := s.svc.Update(ctx, o); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order: %v", err)
	}
	resp := &proto.OrderResponse{
		Id:     o.ID,
		UserId: o.UserID,
		Status: o.Status,
		Total:  o.Total,
	}
	for _, item := range o.Items {
		resp.Items = append(resp.Items, &proto.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
	}
	return resp, nil
}

func (s *Server) ListOrders(ctx context.Context, req *proto.ListOrdersRequest) (*proto.ListOrdersResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if req.Page <= 0 || req.PageSize <= 0 {
		return nil, status.Error(codes.InvalidArgument, "page and pageSize must be positive")
	}
	orders, total, err := s.svc.List(ctx, req.UserId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}
	resp := &proto.ListOrdersResponse{
		Total: int32(total),
	}
	for _, o := range orders {
		orderResp := &proto.OrderResponse{
			Id:     o.ID,
			UserId: o.UserID,
			Status: o.Status,
			Total:  o.Total,
		}
		for _, item := range o.Items {
			orderResp.Items = append(orderResp.Items, &proto.OrderItem{
				ProductId: item.ProductID,
				Quantity:  int32(item.Quantity),
			})
		}
		resp.Orders = append(resp.Orders, orderResp)
	}
	return resp, nil
}
