package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email     string    `gorm:"type:varchar(255);not null;unique;index"`
	GoogleID  string    `gorm:"type:varchar(255);not null;unique;index"`
	Name      string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

type InviteCode struct {
	ID        string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code      string     `gorm:"type:varchar(50);not null;unique;index"`
	UsedBy    *string    `gorm:"type:uuid;index"` // User ID who used it, nullable
	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UsedAt    *time.Time `gorm:""` // When it was used, nullable
}

// BeforeCreate hook to generate UUID if not set
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

func (ic *InviteCode) BeforeCreate(tx *gorm.DB) error {
	if ic.ID == "" {
		ic.ID = uuid.New().String()
	}
	return nil
}
