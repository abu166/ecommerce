package infrastructure

import (
	"context"
	"ecommerce/internal/order/domain"
	"errors"
	"github.com/sirupsen/logrus"
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
	// Ensure the schema is up-to-date with the domain structs
	if err := db.AutoMigrate(&domain.Order{}, &domain.OrderItem{}); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err.Error(),
			"timestamp": "01:38 AM +05, Tuesday, May 20, 2025",
		}).Error("Failed to auto-migrate database schema")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"timestamp": "01:38 AM +05, Tuesday, May 20, 2025",
	}).Info("Database schema migrated successfully")
	return &Repository{db: db}, nil
}

func (r *Repository) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logrus.WithFields(logrus.Fields{
				"transaction_status": "panic_rolled_back",
				"timestamp":          "01:38 AM +05, Tuesday, May 20, 2025",
			}).Error("Transaction panicked and rolled back")
		}
	}()

	err := fn(ctx)
	if err != nil {
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			logrus.WithFields(logrus.Fields{
				"error":              rollbackErr.Error(),
				"transaction_status": "rollback_failed",
				"timestamp":          "01:38 AM +05, Tuesday, May 20, 2025",
			}).Error("Failed to rollback transaction")
		} else {
			logrus.WithFields(logrus.Fields{
				"transaction_status": "rolled_back",
				"timestamp":          "01:38 AM +05, Tuesday, May 20, 2025",
			}).Info("Transaction rolled back successfully")
		}
		return err
	}

	return tx.Commit().Error
}

func (r *Repository) Create(ctx context.Context, o *domain.Order) error {
	logrus.WithFields(logrus.Fields{
		"order_id_before_create": o.ID,
		"timestamp":              "01:38 AM +05, Tuesday, May 20, 2025",
	}).Info("Creating order with ID")
	result := r.db.WithContext(ctx).Create(o)
	if result.Error != nil {
		logrus.WithFields(logrus.Fields{
			"order_id":  o.ID,
			"error":     result.Error.Error(),
			"timestamp": "01:38 AM +05, Tuesday, May 20, 2025",
		}).Error("Failed to create order")
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("failed to create order")
	}
	logrus.WithFields(logrus.Fields{
		"order_id":  o.ID,
		"timestamp": "01:38 AM +05, Tuesday, May 20, 2025",
	}).Info("Order created successfully in database")
	return nil
}

func (r *Repository) CreateItem(ctx context.Context, item *domain.OrderItem) error {
	result := r.db.WithContext(ctx).Create(item)
	if result.Error != nil {
		logrus.WithFields(logrus.Fields{
			"order_id":   item.OrderID,
			"product_id": item.ProductID,
			"error":      result.Error.Error(),
			"timestamp":  "01:38 AM +05, Tuesday, May 20, 2025",
		}).Error("Failed to create order item")
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("failed to create order item")
	}
	return nil
}

func (r *Repository) GetItem(ctx context.Context, orderID, productID string) (*domain.OrderItem, error) {
	var item domain.OrderItem
	result := r.db.WithContext(ctx).Where("order_id = ? AND product_id = ?", orderID, productID).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *Repository) UpdateItem(ctx context.Context, item *domain.OrderItem) error {
	result := r.db.WithContext(ctx).Model(&domain.OrderItem{}).Where("order_id = ? AND product_id = ?", item.OrderID, item.ProductID).Update("quantity", item.Quantity)
	if result.Error != nil {
		logrus.WithFields(logrus.Fields{
			"order_id":   item.OrderID,
			"product_id": item.ProductID,
			"error":      result.Error.Error(),
			"timestamp":  "01:38 AM +05, Tuesday, May 20, 2025",
		}).Error("Failed to update order item")
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("failed to update order item")
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (*domain.Order, error) {
	var o domain.Order
	result := r.db.WithContext(ctx).Preload("Items").First(&o, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, result.Error
	}
	return &o, nil
}

func (r *Repository) Update(ctx context.Context, o *domain.Order) error {
	result := r.db.WithContext(ctx).Save(o)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("failed to update order")
	}
	return nil
}

func (r *Repository) List(ctx context.Context, userID string, page, pageSize int) ([]*domain.Order, int, error) {
	var orders []*domain.Order
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}
