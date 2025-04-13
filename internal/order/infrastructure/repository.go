package infrastructure

import (
	"context"
	"ecommerce/internal/order/domain"
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
	if err := db.AutoMigrate(&domain.Order{}, &domain.OrderItem{}); err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) Create(ctx context.Context, o *domain.Order) error {
	o.ID = uuid.New().String()
	for i := range o.Items {
		o.Items[i].OrderID = o.ID
	}
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *Repository) Get(ctx context.Context, id string) (*domain.Order, error) {
	var o domain.Order
	err := r.db.WithContext(ctx).Preload("Items").First(&o, "id = ?", id).Error
	return &o, err
}

func (r *Repository) Update(ctx context.Context, o *domain.Order) error {
	return r.db.WithContext(ctx).Save(o).Error
}

func (r *Repository) List(ctx context.Context, userID string, page, pageSize int) ([]*domain.Order, int, error) {
	var orders []*domain.Order
	query := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID)
	var total int64
	query.Model(&domain.Order{}).Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error
	return orders, int(total), err
}
