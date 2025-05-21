package application

import (
	"context"
	"ecommerce/internal/user/domain"
	"ecommerce/internal/user/infrastructure"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo  *infrastructure.Repository
	cache infrastructure.Cache
}

func NewService(repo *infrastructure.Repository, cache infrastructure.Cache) *Service {
	return &Service{repo: repo, cache: cache}
}

func (s *Service) Register(ctx context.Context, u *domain.User) error {
	if u.Username == "" || u.Password == "" || u.Email == "" {
		return errors.New("username, password, and email are required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithError(err).Error("Failed to hash password")
		return err
	}
	u.Password = string(hash)

	if err := s.repo.Create(ctx, u); err != nil {
		logrus.WithError(err).Error("Failed to create user in repository")
		return err
	}

	if err := s.cache.SetUser(ctx, u); err != nil {
		logrus.WithError(err).Warn("Failed to cache user, proceeding")
	}
	logrus.WithField("user_id", u.ID).Info("User registered successfully")
	return nil
}

func (s *Service) Authenticate(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username and password are required")
	}
	u, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user by username")
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		logrus.WithField("username", username).Error("Invalid password")
		return "", errors.New("invalid credentials")
	}
	logrus.WithField("user_id", u.ID).Info("User authenticated successfully")
	return u.ID, nil
}

func (s *Service) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	cachedUser, err := s.cache.GetUser(ctx, uuidID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user from cache")
		return nil, err
	}
	if cachedUser != nil {
		logrus.WithField("user_id", id).Info("Returning cached user profile")
		return cachedUser, nil
	}

	user, err := s.repo.Get(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user from repository")
		return nil, err
	}

	if err := s.cache.SetUser(ctx, user); err != nil {
		logrus.WithError(err).Warn("Failed to cache user, proceeding")
	}
	logrus.WithField("user_id", id).Info("User profile retrieved and cached")
	return user, nil
}
