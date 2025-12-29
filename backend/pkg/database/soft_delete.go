package database

import (
	"time"

	"gorm.io/gorm"
)

// SoftDeleteHelper provides utilities for working with soft deletes
type SoftDeleteHelper struct {
	db *gorm.DB
}

// NewSoftDeleteHelper creates a new SoftDeleteHelper instance
func NewSoftDeleteHelper(db *gorm.DB) *SoftDeleteHelper {
	return &SoftDeleteHelper{db: db}
}

// Restore restores a soft-deleted record by setting deleted_at to NULL
func (h *SoftDeleteHelper) Restore(model interface{}, id string) error {
	return h.db.Unscoped().
		Model(model).
		Where("id = ?", id).
		Update("deleted_at", nil).
		Error
}

// PermanentDelete permanently deletes a record (hard delete)
func (h *SoftDeleteHelper) PermanentDelete(model interface{}, id string) error {
	return h.db.Unscoped().
		Where("id = ?", id).
		Delete(model).
		Error
}

// FindDeleted finds a soft-deleted record by ID
func (h *SoftDeleteHelper) FindDeleted(model interface{}, id string) error {
	return h.db.Unscoped().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		First(model).
		Error
}

// ListDeleted lists all soft-deleted records for a given model
func (h *SoftDeleteHelper) ListDeleted(model interface{}, userID string, limit, offset int) ([]interface{}, int64, error) {
	var results []interface{}
	var total int64

	query := h.db.Unscoped().
		Model(model).
		Where("deleted_at IS NOT NULL")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get results
	if err := query.
		Order("deleted_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// CleanupDeleted removes records that were soft-deleted more than the specified days ago
func (h *SoftDeleteHelper) CleanupDeleted(model interface{}, daysOld int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -daysOld)

	var result *gorm.DB
	result = h.db.Unscoped().
		Where("deleted_at IS NOT NULL AND deleted_at < ?", cutoffDate).
		Delete(model)

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// IsDeleted checks if a record is soft-deleted
func (h *SoftDeleteHelper) IsDeleted(model interface{}, id string) (bool, error) {
	var count int64
	err := h.db.Unscoped().
		Model(model).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
