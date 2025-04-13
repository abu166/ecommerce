package domain

import "time"

type Order struct {
	ID        string `gorm:"primaryKey;type:uuid"` // CHANGED: Added type:uuid
	UserID    string
	Items     []OrderItem
	Status    string
	Total     float64
	CreatedAt time.Time // ADDED: To support timestamps
	UpdatedAt time.Time // ADDED: To support timestamps
}

type OrderItem struct {
	OrderID   string `gorm:"primaryKey;type:uuid"` // CHANGED: Added type:uuid
	ProductID string
	Quantity  int
}
