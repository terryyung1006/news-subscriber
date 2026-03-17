package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserMemory stores long-term memory/preferences for a user
type UserMemory struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `gorm:"type:uuid;not null;index"`
	Title     string    `gorm:"type:varchar(100);not null"`
	Content   string    `gorm:"type:text;not null"`
	Category  string    `gorm:"type:varchar(50)"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

// BeforeCreate hook to generate UUID if not set
func (m *UserMemory) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
