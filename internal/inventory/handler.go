package inventory

import (
	"context"
	"ecommerce/internal/inventory/application"
	"ecommerce/internal/inventory/domain"
	"ecommerce/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.UnimplementedInventoryServiceServer
	svc *application.Service
}

func NewServer(svc *application.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.ProductResponse, error) {
	p := &domain.Product{
		Name:     req.Name,
		Category: req.Category,
		Stock:    int(req.Stock),
		Price:    req.Price,
	}
	if err := s.svc.Create(ctx, p); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}
	return &proto.ProductResponse{
		Id:       p.ID,
		Name:     p.Name,
		Category: p.Category,
		Stock:    int32(p.Stock),
		Price:    p.Price,
	}, nil
}

func (s *Server) GetProduct(ctx context.Context, req *proto.GetProductRequest) (*proto.ProductResponse, error) {
	p, err := s.svc.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}
	return &proto.ProductResponse{
		Id:       p.ID,
		Name:     p.Name,
		Category: p.Category,
		Stock:    int32(p.Stock),
		Price:    p.Price,
	}, nil
}

func (s *Server) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.ProductResponse, error) {
	p := &domain.Product{
		ID:       req.Id,
		Name:     req.Name,
		Category: req.Category,
		Stock:    int(req.Stock),
		Price:    req.Price,
	}
	if err := s.svc.Update(ctx, p); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}
	return &proto.ProductResponse{
		Id:       p.ID,
		Name:     p.Name,
		Category: p.Category,
		Stock:    int32(p.Stock),
		Price:    p.Price,
	}, nil
}

func (s *Server) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.Empty, error) {
	if err := s.svc.Delete(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}
	return &proto.Empty{}, nil
}

func (s *Server) ListProducts(ctx context.Context, req *proto.ListProductsRequest) (*proto.ListProductsResponse, error) {
	products, total, err := s.svc.List(ctx, int(req.Page), int(req.PageSize), req.Category)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}
	var resp []*proto.ProductResponse
	for _, p := range products {
		resp = append(resp, &proto.ProductResponse{
			Id:       p.ID,
			Name:     p.Name,
			Category: p.Category,
			Stock:    int32(p.Stock),
			Price:    p.Price,
		})
	}
	return &proto.ListProductsResponse{
		Products: resp,
		Total:    int32(total),
	}, nil
}
