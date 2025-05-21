package application

import (
	"context"
	"ecommerce/internal/inventory/domain"
	"ecommerce/internal/inventory/infrastructure"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Service defines the application logic for the inventory service.
type Service struct {
	repo  *infrastructure.Repository
	cache infrastructure.Cache
}

// NewService creates a new inventory service.
func NewService(repo *infrastructure.Repository, cache infrastructure.Cache) *Service {
	return &Service{repo: repo, cache: cache}
}

// Create creates a new product.
func (s *Service) Create(ctx context.Context, p *domain.Product) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	if err := s.repo.Create(ctx, p); err != nil {
		logrus.WithError(err).Error("Failed to create product")
		return err
	}
	if err := s.cache.SetProduct(ctx, uuid.MustParse(p.ID), p); err != nil {
		logrus.WithError(err).Warn("Failed to cache product, proceeding")
	}
	logrus.WithField("product_id", p.ID).Info("Product created")
	return nil
}

// Get retrieves a product by ID with caching.
func (s *Service) Get(ctx context.Context, id string) (*domain.Product, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	// Check cache first
	if cachedProduct, err := s.cache.GetProduct(ctx, uuidID); err == nil && cachedProduct != nil {
		logrus.WithField("product_id", id).Info("Cache hit for product")
		return cachedProduct, nil
	}

	// Cache miss, query database
	product, err := s.repo.Get(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("Failed to get product")
		return nil, err
	}

	// Cache the product
	if err := s.cache.SetProduct(ctx, uuidID, product); err != nil {
		logrus.WithError(err).Warn("Failed to cache product, proceeding")
	}
	logrus.WithField("product_id", id).Info("Product retrieved and cached")
	return product, nil
}

// Update updates a product with transaction and cache invalidation.
func (s *Service) Update(ctx context.Context, p *domain.Product) error {
	uuidID, err := uuid.Parse(p.ID)
	if err != nil {
		return err
	}

	// Simulate a transaction with stock update and log
	err = s.repo.WithTransaction(ctx, func(txCtx context.Context) error {
		// Update the product
		if err := s.repo.Update(txCtx, p); err != nil {
			logrus.WithError(err).Error("Failed to update product in transaction")
			return err
		}

		// Simulate logging the update (in a real scenario, this would write to a Log table)
		logrus.WithFields(logrus.Fields{
			"product_id": p.ID,
			"stock":      p.Stock,
			"price":      p.Price,
		}).Info("Product update logged")

		return nil
	})
	if err != nil {
		logrus.WithError(err).Error("Transaction failed for product update")
		return err
	}

	// Invalidate cache
	if err := s.cache.DeleteProduct(ctx, uuidID); err != nil {
		logrus.WithError(err).Warn("Failed to invalidate product cache, proceeding")
	}
	logrus.WithField("product_id", p.ID).Info("Product updated and cache invalidated")
	return nil
}

// Delete deletes a product and invalidates cache.
func (s *Service) Delete(ctx context.Context, id string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		logrus.WithError(err).Error("Failed to delete product")
		return err
	}

	// Invalidate cache
	if err := s.cache.DeleteProduct(ctx, uuidID); err != nil {
		logrus.WithError(err).Warn("Failed to invalidate product cache, proceeding")
	}
	logrus.WithField("product_id", id).Info("Product deleted and cache invalidated")
	return nil
}

// List lists products with pagination.
func (s *Service) List(ctx context.Context, page, pageSize int) ([]*domain.Product, int, error) {
	products, total, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		logrus.WithError(err).Error("Failed to list products")
		return nil, 0, err
	}
	logrus.WithFields(logrus.Fields{"page": page, "page_size": pageSize}).Info("Products listed")
	return products, total, nil
}
