package infrastructure

import (
	"context"
	"ecommerce/internal/user/domain"
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
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) Create(ctx context.Context, u *domain.User) error {
	u.ID = uuid.New().String()
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *Repository) Get(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error
	return &u, err
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).First(&u, "username = ?", username).Error
	return &u, err
}
