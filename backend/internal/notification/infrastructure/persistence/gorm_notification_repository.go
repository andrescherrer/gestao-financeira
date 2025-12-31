package persistence

import (
	"errors"
	"fmt"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/notification/domain/entities"
	"gestao-financeira/backend/internal/notification/domain/repositories"
	notificationvalueobjects "gestao-financeira/backend/internal/notification/domain/valueobjects"

	"gorm.io/gorm"
)

// GormNotificationRepository implements NotificationRepository using GORM.
type GormNotificationRepository struct {
	db *gorm.DB
}

// NewGormNotificationRepository creates a new GORM notification repository.
func NewGormNotificationRepository(db *gorm.DB) repositories.NotificationRepository {
	return &GormNotificationRepository{db: db}
}

// FindByID finds a notification by its ID.
func (r *GormNotificationRepository) FindByID(id notificationvalueobjects.NotificationID) (*entities.Notification, error) {
	var model NotificationModel
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id.Value()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find notification by ID: %w", err)
	}

	return r.toDomain(&model)
}

// FindByUserID finds all notifications for a given user.
func (r *GormNotificationRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Notification, error) {
	var models []NotificationModel
	if err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID.Value()).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find notifications by user ID: %w", err)
	}

	notifications := make([]*entities.Notification, 0, len(models))
	for _, model := range models {
		notification, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert notification model to domain: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// FindByUserIDWithFiltersWithPagination finds notifications for a given user with filters and pagination.
func (r *GormNotificationRepository) FindByUserIDWithFiltersWithPagination(
	userID identityvalueobjects.UserID,
	status string,
	notifType string,
	offset, limit int,
) ([]*entities.Notification, int64, error) {
	var models []NotificationModel
	var total int64

	query := r.db.Model(&NotificationModel{}).Where("user_id = ? AND deleted_at IS NULL", userID.Value())

	// Apply status filter if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Apply type filter if provided
	if notifType != "" {
		query = query.Where("type = ?", notifType)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find notifications: %w", err)
	}

	notifications := make([]*entities.Notification, 0, len(models))
	for _, model := range models {
		notification, err := r.toDomain(&model)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert notification model to domain: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, total, nil
}

// FindUnreadByUserID finds all unread notifications for a given user.
func (r *GormNotificationRepository) FindUnreadByUserID(userID identityvalueobjects.UserID) ([]*entities.Notification, error) {
	var models []NotificationModel
	if err := r.db.Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID.Value(), notificationvalueobjects.StatusUnread).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find unread notifications by user ID: %w", err)
	}

	notifications := make([]*entities.Notification, 0, len(models))
	for _, model := range models {
		notification, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert notification model to domain: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// Save saves or updates a notification.
func (r *GormNotificationRepository) Save(notification *entities.Notification) error {
	model := r.toModel(notification)

	// Check if notification exists
	var existing NotificationModel
	err := r.db.Where("id = ?", model.ID).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new notification
		if err := r.db.Create(model).Error; err != nil {
			return fmt.Errorf("failed to create notification: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check notification existence: %w", err)
	} else {
		// Update existing notification
		if err := r.db.Save(model).Error; err != nil {
			return fmt.Errorf("failed to update notification: %w", err)
		}
	}

	return nil
}

// Delete deletes a notification by its ID (soft delete).
func (r *GormNotificationRepository) Delete(id notificationvalueobjects.NotificationID) error {
	if err := r.db.Where("id = ?", id.Value()).Delete(&NotificationModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// Exists checks if a notification with the given ID already exists.
func (r *GormNotificationRepository) Exists(id notificationvalueobjects.NotificationID) (bool, error) {
	var count int64
	if err := r.db.Model(&NotificationModel{}).Where("id = ? AND deleted_at IS NULL", id.Value()).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check notification existence: %w", err)
	}
	return count > 0, nil
}

// Count returns the total number of notifications for a given user.
func (r *GormNotificationRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	var count int64
	if err := r.db.Model(&NotificationModel{}).Where("user_id = ? AND deleted_at IS NULL", userID.Value()).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count notifications: %w", err)
	}
	return count, nil
}

// CountUnread returns the number of unread notifications for a given user.
func (r *GormNotificationRepository) CountUnread(userID identityvalueobjects.UserID) (int64, error) {
	var count int64
	if err := r.db.Model(&NotificationModel{}).Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID.Value(), notificationvalueobjects.StatusUnread).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count unread notifications: %w", err)
	}
	return count, nil
}

// toDomain converts a persistence model to a domain entity.
func (r *GormNotificationRepository) toDomain(model *NotificationModel) (*entities.Notification, error) {
	// Create value objects
	notificationID, err := notificationvalueobjects.NewNotificationID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid notification ID: %w", err)
	}

	userID, err := identityvalueobjects.NewUserID(model.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	title, err := notificationvalueobjects.NewNotificationTitle(model.Title)
	if err != nil {
		return nil, fmt.Errorf("invalid notification title: %w", err)
	}

	message, err := notificationvalueobjects.NewNotificationMessage(model.Message)
	if err != nil {
		return nil, fmt.Errorf("invalid notification message: %w", err)
	}

	notifType, err := notificationvalueobjects.NewNotificationType(model.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid notification type: %w", err)
	}

	status, err := notificationvalueobjects.NewNotificationStatus(model.Status)
	if err != nil {
		return nil, fmt.Errorf("invalid notification status: %w", err)
	}

	// Convert JSONB metadata to map[string]interface{}
	var metadata map[string]interface{}
	if model.Metadata != nil {
		metadata = map[string]interface{}(model.Metadata)
	} else {
		metadata = make(map[string]interface{})
	}

	return entities.NotificationFromPersistence(
		notificationID,
		userID,
		title,
		message,
		notifType,
		status,
		model.ReadAt,
		metadata,
		model.CreatedAt,
		model.UpdatedAt,
	), nil
}

// toModel converts a domain entity to a persistence model.
func (r *GormNotificationRepository) toModel(notification *entities.Notification) *NotificationModel {
	// Convert metadata to JSONB
	var metadata JSONB
	if notification.Metadata() != nil {
		metadata = JSONB(notification.Metadata())
	} else {
		metadata = make(JSONB)
	}

	return &NotificationModel{
		ID:        notification.ID().Value(),
		UserID:    notification.UserID().Value(),
		Title:     notification.Title().Value(),
		Message:   notification.Message().Value(),
		Type:      notification.Type().Value(),
		Status:    notification.Status().Value(),
		ReadAt:    notification.ReadAt(),
		Metadata:  metadata,
		CreatedAt: notification.CreatedAt(),
		UpdatedAt: notification.UpdatedAt(),
	}
}
