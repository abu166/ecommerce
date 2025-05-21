package domain

import (
	"time"
)

type Order struct {
	ID        string      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    string      `gorm:"type:uuid;not null"`
	Items     []OrderItem `gorm:"foreignKey:OrderID"`
	Total     float64     `gorm:"not null"`
	Status    string      `gorm:"default:'pending'"`
	CreatedAt time.Time   `gorm:"autoCreateTime"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime"`
}

type OrderItem struct {
	OrderID   string `gorm:"type:uuid;primaryKey"`
	ProductID string `gorm:"type:uuid;primaryKey"`
	Quantity  int    `gorm:"not null"`
}
