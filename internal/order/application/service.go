package application

import (
	"context"
	"ecommerce/internal/inventory/application"
	"ecommerce/internal/order/domain"
	"ecommerce/internal/order/infrastructure"
	"errors"
)

type Service struct {
	repo *infrastructure.Repository
	inv  *application.Service
}

func NewService(repo *infrastructure.Repository, inv *application.Service) *Service {
	return &Service{repo: repo, inv: inv}
}

func (s *Service) Create(ctx context.Context, o *domain.Order) error {
	var total float64
	for i, item := range o.Items {
		p, err := s.inv.Get(ctx, item.ProductID)
		if err != nil {
			return err
		}
		if p.Stock < item.Quantity {
			return errors.New("insufficient stock")
		}
		total += p.Price * float64(item.Quantity)
		p.Stock -= item.Quantity
		if err := s.inv.Update(ctx, p); err != nil {
			return err
		}
		o.Items[i].OrderID = o.ID
	}
	o.Total = total
	o.Status = "pending"
	return s.repo.Create(ctx, o)
}

func (s *Service) Get(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, o *domain.Order) error {
	return s.repo.Update(ctx, o)
}

func (s *Service) List(ctx context.Context, userID string, page, pageSize int) ([]*domain.Order, int, error) {
	return s.repo.List(ctx, userID, page, pageSize)
}
