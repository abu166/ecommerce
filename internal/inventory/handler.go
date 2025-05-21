package inventory

import (
	"context"
	"ecommerce/internal/inventory/application"
	"ecommerce/internal/inventory/domain"
	"ecommerce/proto"
	"github.com/google/uuid"
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
		ID:       uuid.New().String(),
		Name:     req.Name,
		Category: req.Category,
		Stock:    int(req.Stock),
		Price:    req.Price,
	}
	if err := s.svc.Create(ctx, p); err != nil {
		return nil, status.Error(codes.Internal, "failed to create product")
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
		if status.Code(err) == codes.NotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "failed to get product")
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
	p, err := s.svc.Get(ctx, req.Id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "failed to get product")
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Category != "" {
		p.Category = req.Category
	}
	if req.Stock != 0 {
		p.Stock += int(req.Stock)
	}
	if req.Price != 0 {
		p.Price = req.Price
	}
	if err := s.svc.Update(ctx, p); err != nil {
		return nil, status.Error(codes.Internal, "failed to update product")
	}
	return &proto.ProductResponse{
		Id:       p.ID,
		Name:     p.Name,
		Category: p.Category,
		Stock:    int32(p.Stock),
		Price:    p.Price,
	}, nil
}

func (s *Server) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.InventoryEmpty, error) {
	if err := s.svc.Delete(ctx, req.Id); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete product")
	}
	return &proto.InventoryEmpty{}, nil
}

func (s *Server) ListProducts(ctx context.Context, req *proto.ListProductsRequest) (*proto.ListProductsResponse, error) {
	products, total, err := s.svc.List(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list products")
	}
	var protoProducts []*proto.ProductResponse
	for _, p := range products {
		protoProducts = append(protoProducts, &proto.ProductResponse{
			Id:       p.ID,
			Name:     p.Name,
			Category: p.Category,
			Stock:    int32(p.Stock),
			Price:    p.Price,
		})
	}
	return &proto.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(total),
	}, nil
}
