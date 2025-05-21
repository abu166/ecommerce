package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"ecommerce/internal/user/domain"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Cache defines the interface for caching operations.
type Cache interface {
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	SetUser(ctx context.Context, user *domain.User) error
}

// RedisCache implements the Cache interface using Redis.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache initializes a new Redis client.
func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // Set password if required
		DB:       0,  // Use default DB
	})
	return &RedisCache{client: client}
}

// GetUser retrieves a user from Redis by ID.
func (c *RedisCache) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	key := "user:" + id.String()
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		logrus.WithField("user_id", id).Info("Cache miss for user")
		return nil, nil
	}
	if err != nil {
		logrus.WithError(err).Error("Failed to get user from cache")
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal(data, &user); err != nil {
		logrus.WithError(err).Error("Failed to unmarshal user from cache")
		return nil, err
	}
	logrus.WithField("user_id", id).Info("Cache hit for user")
	return &user, nil
}

// SetUser stores a user in Redis with a TTL of 1 hour.
func (c *RedisCache) SetUser(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		return errors.New("user ID is required")
	}
	key := "user:" + user.ID
	data, err := json.Marshal(user)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal user for cache")
		return err
	}
	if err := c.client.Set(ctx, key, data, time.Hour).Err(); err != nil {
		logrus.WithError(err).Error("Failed to set user in cache")
		return err
	}
	logrus.WithField("user_id", user.ID).Info("User cached successfully")
	return nil
}
