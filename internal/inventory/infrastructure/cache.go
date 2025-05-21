package infrastructure

import (
	"context"
	"encoding/json"
	"time"

	"ecommerce/internal/inventory/domain"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Cache defines the interface for caching operations.
type Cache interface {
	GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	SetProduct(ctx context.Context, id uuid.UUID, product *domain.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

// RedisCache implements the Cache interface using Redis.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache initializes a new Redis client.
func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{client: client}
}

// GetProduct retrieves a product from Redis by ID.
func (c *RedisCache) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	key := "product:" + id.String()
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		logrus.WithField("product_id", id).Info("Cache miss for product")
		return nil, nil
	}
	if err != nil {
		logrus.WithError(err).Error("Failed to get product from cache")
		return nil, err
	}

	var product domain.Product
	if err := json.Unmarshal(data, &product); err != nil {
		logrus.WithError(err).Error("Failed to unmarshal product from cache")
		return nil, err
	}
	logrus.WithField("product_id", id).Info("Cache hit for product")
	return &product, nil
}

// SetProduct stores a product in Redis with a TTL of 1 hour.
func (c *RedisCache) SetProduct(ctx context.Context, id uuid.UUID, product *domain.Product) error {
	key := "product:" + id.String()
	data, err := json.Marshal(product)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal product for cache")
		return err
	}
	if err := c.client.Set(ctx, key, data, time.Hour).Err(); err != nil {
		logrus.WithError(err).Error("Failed to set product in cache")
		return err
	}
	logrus.WithField("product_id", id).Info("Product cached successfully")
	return nil
}

// DeleteProduct removes a product from Redis by ID.
func (c *RedisCache) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	key := "product:" + id.String()
	if err := c.client.Del(ctx, key).Err(); err != nil {
		logrus.WithError(err).Error("Failed to delete product from cache")
		return err
	}
	logrus.WithField("product_id", id).Info("Product cache invalidated")
	return nil
}
