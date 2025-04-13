package domain

import "time"

type Product struct {
	ID        string `gorm:"primaryKey;type:uuid"` // CHANGED: Added type:uuid
	Name      string
	Category  string
	Stock     int
	Price     float64
	CreatedAt time.Time // ADDED: To support timestamps
	UpdatedAt time.Time // ADDED: To support timestamps
}
