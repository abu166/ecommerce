package infrastructure

import (
	"context"
	"ecommerce/internal/inventory/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

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

func (r *Repository) Create(ctx context.Context, p *domain.Product) error {
	p.ID = uuid.New().String()
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repository) Get(ctx context.Context, id string) (*domain.Product, error) {
	var p domain.Product
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	return &p, err
}

func (r *Repository) Update(ctx context.Context, p *domain.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, "id = ?", id).Error
}

func (r *Repository) List(ctx context.Context, page, pageSize int, category string) ([]*domain.Product, int, error) {
	var products []*domain.Product
	query := r.db.WithContext(ctx)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	var total int64
	query.Model(&domain.Product{}).Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error
	return products, int(total), err
}
