package application

import (
	"context"
	"ecommerce/internal/user/domain"
	"ecommerce/internal/user/infrastructure"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *infrastructure.Repository
}

func NewService(repo *infrastructure.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, u *domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return s.repo.Create(ctx, u)
}

func (s *Service) Authenticate(ctx context.Context, username, password string) (string, error) {
	u, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	return u.ID, nil
}

func (s *Service) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.Get(ctx, id)
}
