package infrastructure

import (
	"context"
	"ecommerce/internal/inventory/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Repository defines the data access layer for the inventory service.
type Repository struct {
	db *gorm.DB
}

// NewRepository initializes a new repository with transaction support.
func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&domain.Product{}); err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// WithTransaction executes a function within a database transaction.
func (r *Repository) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	err := fn(ctx)
	if err != nil {
		return err
	}

	return tx.Commit().Error
}

// Create creates a new product.
func (r *Repository) Create(ctx context.Context, p *domain.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

// Get retrieves a product by ID.
func (r *Repository) Get(ctx context.Context, id string) (*domain.Product, error) {
	var p domain.Product
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// Update updates a product.
func (r *Repository) Update(ctx context.Context, p *domain.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

// Delete deletes a product by ID.
func (r *Repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, "id = ?", id).Error
}

// List lists products with pagination.
func (r *Repository) List(ctx context.Context, page, pageSize int) ([]*domain.Product, int, error) {
	var products []*domain.Product
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, int(total), nil
}
