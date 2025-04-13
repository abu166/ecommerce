package application

import (
	"context"
	"ecommerce/internal/inventory/domain"
	"ecommerce/internal/inventory/infrastructure"
)

type Service struct {
	repo *infrastructure.Repository
}

func NewService(repo *infrastructure.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, p *domain.Product) error {
	return s.repo.Create(ctx, p)
}

func (s *Service) Get(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, p *domain.Product) error {
	return s.repo.Update(ctx, p)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, page, pageSize int, category string) ([]*domain.Product, int, error) {
	return s.repo.List(ctx, page, pageSize, category)
}
