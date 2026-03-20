package postgres

import (
	"context"
	"errors"
	"news-subscriber-core/src/models"

	"gorm.io/gorm"
)

type MemoryRepository struct {
	db *DB
}

func NewMemoryRepository(db *DB) *MemoryRepository {
	return &MemoryRepository{db: db}
}

func (r *MemoryRepository) Create(ctx context.Context, userID, title, content, category string) (*models.UserMemory, error) {
	memory := &models.UserMemory{
		UserID:   userID,
		Title:    title,
		Content:  content,
		Category: category,
	}

	if err := r.db.WithContext(ctx).Create(memory).Error; err != nil {
		return nil, err
	}

	return memory, nil
}

func (r *MemoryRepository) GetByUserID(ctx context.Context, userID string) ([]*models.UserMemory, error) {
	var memories []*models.UserMemory
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&memories).Error

	if err != nil {
		return nil, err
	}

	return memories, nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*models.UserMemory, error) {
	var memory models.UserMemory
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&memory).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &memory, nil
}

func (r *MemoryRepository) Update(ctx context.Context, id, title, content string) (*models.UserMemory, error) {
	err := r.db.WithContext(ctx).Model(&models.UserMemory{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"title":   title,
			"content": content,
		}).Error

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.UserMemory{}).Error
}

func (r *MemoryRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.UserMemory{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
