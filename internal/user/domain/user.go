package domain

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;type:uuid"` // CHANGED: Added type:uuid
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"unique"` // CHANGED: Added unique to match schema
	CreatedAt time.Time // ADDED: To match existing created_at column
	UpdatedAt time.Time // ADDED: To match existing updated_at column
}
