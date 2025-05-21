package application

import (
	"context"
	"ecommerce/internal/order/domain"
	"ecommerce/internal/order/infrastructure"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Service defines the application logic for the order service.
type Service struct {
	repo  *infrastructure.Repository
	cache infrastructure.Cache
}

// NewService creates a new order service.
func NewService(repo *infrastructure.Repository, cache infrastructure.Cache) *Service {
	return &Service{repo: repo, cache: cache}
}

// Create creates a new order with transaction support.
func (s *Service) Create(ctx context.Context, o *domain.Order) error {
	// Validate required fields
	if o.UserID == "" || len(o.Items) == 0 || o.Total <= 0 {
		return errors.New("user ID, items, and total are required")
	}

	// Validate user_id is a valid UUID
	if _, err := uuid.Parse(o.UserID); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":   o.UserID,
			"error":     err.Error(),
			"timestamp": "02:08 AM +05, Tuesday, May 20, 2025",
		}).Error("Invalid user ID format")
		return errors.New("invalid user ID format: must be a valid UUID")
	}

	// Validate each product_id is a valid UUID and quantity is positive
	for i, item := range o.Items {
		if item.ProductID == "" || item.Quantity <= 0 {
			return errors.New("invalid order item: product ID and quantity are required")
		}
		if _, err := uuid.Parse(item.ProductID); err != nil {
			logrus.WithFields(logrus.Fields{
				"product_id": item.ProductID,
				"item_index": i,
				"error":      err.Error(),
				"timestamp":  "02:08 AM +05, Tuesday, May 20, 2025",
			}).Error("Invalid product ID format")
			return errors.New("invalid product ID format: must be a valid UUID")
		}
	}

	// Create a new Order object, preserving the provided ID
	newOrder := &domain.Order{
		ID:     o.ID, // Use the ID from the input order
		UserID: o.UserID,
		Total:  o.Total,
		Status: o.Status,
	}
	// Validate the provided ID or generate a new one
	if newOrder.ID == "" {
		newOrder.ID = uuid.New().String()
		logrus.WithFields(logrus.Fields{
			"order_id":  newOrder.ID,
			"timestamp": "02:08 AM +05, Tuesday, May 20, 2025",
		}).Info("Generated new UUID for order")
	} else {
		// Validate the provided ID is a valid UUID
		if _, err := uuid.Parse(newOrder.ID); err != nil {
			logrus.WithFields(logrus.Fields{
				"order_id":  newOrder.ID,
				"error":     err.Error(),
				"timestamp": "02:08 AM +05, Tuesday, May 20, 2025",
			}).Error("Invalid order ID format")
			return errors.New("invalid order ID format: must be a valid UUID")
		}
	}

	// Deep copy OrderItems to avoid retaining old OrderID
	newOrder.Items = make([]domain.OrderItem, len(o.Items))
	for i, item := range o.Items {
		newOrder.Items[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	logrus.WithFields(logrus.Fields{
		"order_id_before_create": newOrder.ID,
		"timestamp":              "02:08 AM +05, Tuesday, May 20, 2025",
	}).Info("Order object before creation")

	err := s.repo.WithTransaction(ctx, func(txCtx context.Context) error {
		// Create order with the provided or generated ID
		if err := s.repo.Create(txCtx, newOrder); err != nil {
			logrus.WithFields(logrus.Fields{
				"order_id":           newOrder.ID,
				"error":              err.Error(),
				"transaction_status": "failed",
				"error_code":         "db_create_order",
				"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
			}).Error("Failed to create order in transaction")
			return err
		}
		logrus.WithFields(logrus.Fields{
			"order_id":           newOrder.ID,
			"transaction_status": "in_progress",
			"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
		}).Info("Order created, proceeding to items")

		// Set OrderID for items and handle duplicates
		for i := range newOrder.Items {
			newOrder.Items[i].OrderID = newOrder.ID
			// Check if the item already exists
			existingItem, err := s.repo.GetItem(txCtx, newOrder.ID, newOrder.Items[i].ProductID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.WithFields(logrus.Fields{
					"order_id":           newOrder.ID,
					"product_id":         newOrder.Items[i].ProductID,
					"item_index":         i,
					"error":              err.Error(),
					"transaction_status": "failed",
					"error_code":         "db_get_item",
					"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
				}).Error("Failed to check existing order item")
				return err
			}
			if existingItem != nil {
				// Item exists, update quantity
				existingItem.Quantity += newOrder.Items[i].Quantity
				if err := s.repo.UpdateItem(txCtx, existingItem); err != nil {
					logrus.WithFields(logrus.Fields{
						"order_id":           newOrder.ID,
						"product_id":         newOrder.Items[i].ProductID,
						"item_index":         i,
						"error":              err.Error(),
						"transaction_status": "failed",
						"error_code":         "db_update_item",
						"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
					}).Error("Failed to update order item in transaction")
					return err
				}
			} else {
				// Item does not exist, create it
				if err := s.repo.CreateItem(txCtx, &newOrder.Items[i]); err != nil {
					logrus.WithFields(logrus.Fields{
						"order_id":           newOrder.ID,
						"product_id":         newOrder.Items[i].ProductID,
						"item_index":         i,
						"error":              err.Error(),
						"transaction_status": "failed",
						"error_code":         "db_create_item",
						"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
					}).Error("Failed to create order item in transaction")
					return err
				}
			}
		}
		logrus.WithFields(logrus.Fields{
			"order_id":           newOrder.ID,
			"transaction_status": "pre_commit",
			"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
		}).Info("All items processed, preparing to commit")
		return nil
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"order_id":           newOrder.ID,
			"error":              err.Error(),
			"transaction_status": "rolled_back",
			"error_code":         "transaction_failed",
			"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
		}).Error("Transaction failed for order creation, rollback attempted")
		return err
	}

	uuidUserID, err := uuid.Parse(newOrder.UserID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":            newOrder.UserID,
			"error":              err.Error(),
			"transaction_status": "warning",
			"error_code":         "invalid_user_id",
			"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
		}).Warn("Invalid user ID for cache invalidation, proceeding")
	} else {
		if err := s.cache.DeleteOrders(ctx, uuidUserID); err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id":            newOrder.UserID,
				"error":              err.Error(),
				"transaction_status": "warning",
				"error_code":         "cache_invalidation_failed",
				"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
			}).Warn("Failed to invalidate orders cache, proceeding")
		}
	}
	logrus.WithFields(logrus.Fields{
		"order_id":           newOrder.ID,
		"transaction_status": "committed",
		"success":            true,
		"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
	}).Info("Order created successfully")
	// Update the input order with the final state
	o.ID = newOrder.ID
	o.Status = newOrder.Status
	o.Items = newOrder.Items
	return nil
}

// Get retrieves an order by ID.
func (s *Service) Get(ctx context.Context, id string) (*domain.Order, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}
	order, err := s.repo.Get(ctx, id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"order_id":   id,
			"error":      err.Error(),
			"error_code": "db_get_order",
			"timestamp":  "02:08 AM +05, Tuesday, May 20, 2025",
		}).Error("Failed to get order")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"order_id":  id,
		"success":   true,
		"timestamp": "02:08 AM +05, Tuesday, May 20, 2025",
	}).Info("Order retrieved successfully")
	return order, nil
}

// Update updates an existing order.
func (s *Service) Update(ctx context.Context, o *domain.Order) error {
	if o.ID == "" {
		return errors.New("order ID is required")
	}
	if err := s.repo.Update(ctx, o); err != nil {
		logrus.WithFields(logrus.Fields{
			"order_id":           o.ID,
			"error":              err.Error(),
			"transaction_status": "failed",
			"error_code":         "db_update_order",
			"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
		}).Error("Failed to update order")
		return err
	}

	uuidUserID, err := uuid.Parse(o.UserID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":            o.UserID,
			"error":              err.Error(),
			"transaction_status": "warning",
			"error_code":         "invalid_user_id",
			"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
		}).Warn("Invalid user ID for cache invalidation, proceeding")
	} else {
		if err := s.cache.DeleteOrders(ctx, uuidUserID); err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id":            o.UserID,
				"error":              err.Error(),
				"transaction_status": "warning",
				"error_code":         "cache_invalidation_failed",
				"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
			}).Warn("Failed to invalidate orders cache, proceeding")
		}
	}
	logrus.WithFields(logrus.Fields{
		"order_id":           o.ID,
		"transaction_status": "committed",
		"success":            true,
		"timestamp":          "02:08 AM +05, Tuesday, May 20, 2025",
	}).Info("Order updated successfully")
	return nil
}

// List lists orders for a user with pagination.
func (s *Service) List(ctx context.Context, userID string, page, pageSize int) ([]*domain.Order, int, error) {
	if userID == "" {
		return nil, 0, errors.New("user ID is required")
	}
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.New("page and pageSize must be positive")
	}
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, errors.New("invalid user ID format")
	}

	cachedOrders, err := s.cache.GetOrders(ctx, uuidUserID, page, pageSize)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":    userID,
			"page":       page,
			"page_size":  pageSize,
			"error":      err.Error(),
			"error_code": "cache_get_orders",
			"timestamp":  "02:08 AM +05, Tuesday, May 20, 2025",
		}).Error("Failed to get orders from cache")
		return nil, 0, err
	}
	if cachedOrders != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":   userID,
			"page":      page,
			"page_size": pageSize,
			"success":   true,
			"timestamp": "02:08 AM +05, Tuesday, May 20, 2025",
		}).Info("Cache hit for orders")
		return cachedOrders, len(cachedOrders), nil
	}

	orders, total, err := s.repo.List(ctx, userID, page, pageSize)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":    userID,
			"page":       page,
			"page_size":  pageSize,
			"error":      err.Error(),
			"error_code": "db_list_orders",
			"timestamp":  "02:08 AM +05, Tuesday, May 20, 2025",
		}).Error("Failed to list orders from repository")
		return nil, 0, err
	}

	if err := s.cache.SetOrders(ctx, uuidUserID, page, pageSize, orders); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":    userID,
			"page":       page,
			"page_size":  pageSize,
			"error":      err.Error(),
			"error_code": "cache_set_orders",
			"timestamp":  "02:08 AM +05, Tuesday, May 20, 2025",
		}).Warn("Failed to cache orders, proceeding")
	}
	logrus.WithFields(logrus.Fields{
		"user_id":   userID,
		"page":      page,
		"page_size": pageSize,
		"success":   true,
		"timestamp": "02:08 AM +05, Tuesday, May 20, 2025",
	}).Info("Orders retrieved and cached")
	return orders, total, nil
}
