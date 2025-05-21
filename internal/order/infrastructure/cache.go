package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"ecommerce/internal/order/domain"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Cache interface {
	GetOrders(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*domain.Order, error)
	SetOrders(ctx context.Context, userID uuid.UUID, page, pageSize int, orders []*domain.Order) error
	DeleteOrders(ctx context.Context, userID uuid.UUID) error
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return &RedisCache{client: client}
}

func (c *RedisCache) GetOrders(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*domain.Order, error) {
	key := fmt.Sprintf("orders:%s:%d:%d", userID.String(), page, pageSize)
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		logrus.WithFields(logrus.Fields{
			"user_id":   userID,
			"page":      page,
			"page_size": pageSize,
		}).Info("Cache miss for orders")
		return nil, nil
	}
	if err != nil {
		logrus.WithError(err).Error("Failed to get orders from cache")
		return nil, err
	}

	var orders []*domain.Order
	if err := json.Unmarshal(data, &orders); err != nil {
		logrus.WithError(err).Error("Failed to unmarshal orders from cache")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"user_id":   userID,
		"page":      page,
		"page_size": pageSize,
	}).Info("Cache hit for orders")
	return orders, nil
}

func (c *RedisCache) SetOrders(ctx context.Context, userID uuid.UUID, page, pageSize int, orders []*domain.Order) error {
	key := fmt.Sprintf("orders:%s:%d:%d", userID.String(), page, pageSize)
	data, err := json.Marshal(orders)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal orders for cache")
		return err
	}
	if err := c.client.Set(ctx, key, data, time.Hour).Err(); err != nil {
		logrus.WithError(err).Error("Failed to set orders in cache")
		return err
	}
	logrus.WithFields(logrus.Fields{
		"user_id":   userID,
		"page":      page,
		"page_size": pageSize,
	}).Info("Orders cached successfully")
	return nil
}

func (c *RedisCache) DeleteOrders(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("invalid user ID")
	}
	pattern := "orders:" + userID.String() + ":*"
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			logrus.WithError(err).Error("Failed to delete orders from cache")
			return err
		}
	}
	if err := iter.Err(); err != nil {
		logrus.WithError(err).Error("Failed to scan orders cache")
		return err
	}
	logrus.WithField("user_id", userID).Info("Orders cache invalidated")
	return nil
}
