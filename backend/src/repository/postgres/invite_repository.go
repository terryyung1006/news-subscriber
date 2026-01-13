package postgres

import (
	"context"
	"errors"
	"news-subscriber-core/src/models"
	"time"

	"gorm.io/gorm"
)

type InviteCodeRepository struct {
	db *DB
}

func NewInviteCodeRepository(db *DB) *InviteCodeRepository {
	return &InviteCodeRepository{db: db}
}

func (r *InviteCodeRepository) IsValid(ctx context.Context, code string) (bool, error) {
	var inviteCode models.InviteCode
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&inviteCode).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil // Code doesn't exist
	}
	if err != nil {
		return false, err
	}

	// Valid if not yet used
	return inviteCode.UsedBy == nil, nil
}

func (r *InviteCodeRepository) MarkAsUsed(ctx context.Context, code, userID string) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&models.InviteCode{}).
		Where("code = ? AND used_by IS NULL", code).
		Updates(map[string]interface{}{
			"used_by": userID,
			"used_at": now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("invite code already used or doesn't exist")
	}

	return nil
}
