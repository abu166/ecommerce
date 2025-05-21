package infrastructure

import (
	"context"
	"ecommerce/internal/user/domain"
	"errors"
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
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	result := r.db.WithContext(ctx).Create(u)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("failed to create user")
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	result := r.db.WithContext(ctx).First(&u, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &u, nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var u domain.User
	result := r.db.WithContext(ctx).First(&u, "username = ?", username)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &u, nil
}
