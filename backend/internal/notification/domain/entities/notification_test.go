package entities

import (
	"testing"
	"time"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	notificationvalueobjects "gestao-financeira/backend/internal/notification/domain/valueobjects"
)

func TestNewNotification(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	title := notificationvalueobjects.MustNotificationTitle("Test Notification")
	message := notificationvalueobjects.MustNotificationMessage("This is a test notification")
	notifType := notificationvalueobjects.MustNotificationType(notificationvalueobjects.TypeInfo)
	metadata := map[string]interface{}{
		"entity_id":   "123",
		"entity_type": "Transaction",
	}

	tests := []struct {
		name      string
		userID    identityvalueobjects.UserID
		title     notificationvalueobjects.NotificationTitle
		message   notificationvalueobjects.NotificationMessage
		notifType notificationvalueobjects.NotificationType
		metadata  map[string]interface{}
		wantErr   bool
	}{
		{
			name:      "valid notification",
			userID:    userID,
			title:     title,
			message:   message,
			notifType: notifType,
			metadata:  metadata,
			wantErr:   false,
		},
		{
			name:      "empty user ID",
			userID:    identityvalueobjects.UserID{},
			title:     title,
			message:   message,
			notifType: notifType,
			metadata:  metadata,
			wantErr:   true,
		},
		{
			name:      "empty title",
			userID:    userID,
			title:     notificationvalueobjects.NotificationTitle{},
			message:   message,
			notifType: notifType,
			metadata:  metadata,
			wantErr:   true,
		},
		{
			name:      "empty message",
			userID:    userID,
			title:     title,
			message:   notificationvalueobjects.NotificationMessage{},
			notifType: notifType,
			metadata:  metadata,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNotification(
				tt.userID,
				tt.title,
				tt.message,
				tt.notifType,
				tt.metadata,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("NewNotification() returned nil for valid input")
			}
			if !tt.wantErr && got != nil {
				if got.ID().IsEmpty() {
					t.Error("NewNotification() should generate a non-empty ID")
				}
				if !got.Status().IsUnread() {
					t.Error("NewNotification() should create notification with UNREAD status")
				}
				if got.ReadAt() != nil {
					t.Error("NewNotification() should create notification with nil ReadAt")
				}
			}
		})
	}
}

func TestNotification_MarkAsRead(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	title := notificationvalueobjects.MustNotificationTitle("Test Notification")
	message := notificationvalueobjects.MustNotificationMessage("This is a test notification")
	notifType := notificationvalueobjects.MustNotificationType(notificationvalueobjects.TypeInfo)

	notification, _ := NewNotification(userID, title, message, notifType, nil)

	// Mark as read
	err := notification.MarkAsRead()
	if err != nil {
		t.Errorf("Notification.MarkAsRead() error = %v, want nil", err)
	}

	if !notification.Status().IsRead() {
		t.Error("Notification.Status() should be READ after MarkAsRead()")
	}

	if notification.ReadAt() == nil {
		t.Error("Notification.ReadAt() should not be nil after MarkAsRead()")
	}

	// Try to mark as read again (should fail)
	err = notification.MarkAsRead()
	if err == nil {
		t.Error("Notification.MarkAsRead() should return error when already read")
	}
}

func TestNotification_MarkAsUnread(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	title := notificationvalueobjects.MustNotificationTitle("Test Notification")
	message := notificationvalueobjects.MustNotificationMessage("This is a test notification")
	notifType := notificationvalueobjects.MustNotificationType(notificationvalueobjects.TypeInfo)

	notification, _ := NewNotification(userID, title, message, notifType, nil)
	notification.MarkAsRead()

	// Mark as unread
	err := notification.MarkAsUnread()
	if err != nil {
		t.Errorf("Notification.MarkAsUnread() error = %v, want nil", err)
	}

	if !notification.Status().IsUnread() {
		t.Error("Notification.Status() should be UNREAD after MarkAsUnread()")
	}

	if notification.ReadAt() != nil {
		t.Error("Notification.ReadAt() should be nil after MarkAsUnread()")
	}

	// Try to mark as unread again (should fail)
	err = notification.MarkAsUnread()
	if err == nil {
		t.Error("Notification.MarkAsUnread() should return error when already unread")
	}
}

func TestNotification_Archive(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	title := notificationvalueobjects.MustNotificationTitle("Test Notification")
	message := notificationvalueobjects.MustNotificationMessage("This is a test notification")
	notifType := notificationvalueobjects.MustNotificationType(notificationvalueobjects.TypeInfo)

	notification, _ := NewNotification(userID, title, message, notifType, nil)

	// Archive
	err := notification.Archive()
	if err != nil {
		t.Errorf("Notification.Archive() error = %v, want nil", err)
	}

	if !notification.Status().IsArchived() {
		t.Error("Notification.Status() should be ARCHIVED after Archive()")
	}

	// Try to archive again (should fail)
	err = notification.Archive()
	if err == nil {
		t.Error("Notification.Archive() should return error when already archived")
	}

	// Try to mark as read when archived (should fail)
	err = notification.MarkAsRead()
	if err == nil {
		t.Error("Notification.MarkAsRead() should return error when archived")
	}
}

func TestNotification_Unarchive(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	title := notificationvalueobjects.MustNotificationTitle("Test Notification")
	message := notificationvalueobjects.MustNotificationMessage("This is a test notification")
	notifType := notificationvalueobjects.MustNotificationType(notificationvalueobjects.TypeInfo)

	notification, _ := NewNotification(userID, title, message, notifType, nil)
	notification.MarkAsRead()
	notification.Archive()

	// Unarchive
	err := notification.Unarchive()
	if err != nil {
		t.Errorf("Notification.Unarchive() error = %v, want nil", err)
	}

	if !notification.Status().IsRead() {
		t.Error("Notification.Status() should be READ after Unarchive() (was read before)")
	}

	// Test unarchive from unread
	notification2, _ := NewNotification(userID, title, message, notifType, nil)
	notification2.Archive()
	notification2.Unarchive()

	if !notification2.Status().IsUnread() {
		t.Error("Notification.Status() should be UNREAD after Unarchive() (was unread before)")
	}

	// Try to unarchive when not archived (should fail)
	err = notification.Unarchive()
	if err == nil {
		t.Error("Notification.Unarchive() should return error when not archived")
	}
}

func TestNotification_UpdateMetadata(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	title := notificationvalueobjects.MustNotificationTitle("Test Notification")
	message := notificationvalueobjects.MustNotificationMessage("This is a test notification")
	notifType := notificationvalueobjects.MustNotificationType(notificationvalueobjects.TypeInfo)

	notification, _ := NewNotification(userID, title, message, notifType, nil)

	originalUpdatedAt := notification.UpdatedAt()
	time.Sleep(10 * time.Millisecond) // Ensure time difference

	newMetadata := map[string]interface{}{
		"new_key": "new_value",
	}
	notification.UpdateMetadata(newMetadata)

	if notification.Metadata()["new_key"] != "new_value" {
		t.Error("Notification.UpdateMetadata() should update metadata")
	}

	if !notification.UpdatedAt().After(originalUpdatedAt) {
		t.Error("Notification.UpdateMetadata() should update UpdatedAt timestamp")
	}
}

func TestNotificationFromPersistence(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	title := notificationvalueobjects.MustNotificationTitle("Test Notification")
	message := notificationvalueobjects.MustNotificationMessage("This is a test notification")
	notifType := notificationvalueobjects.MustNotificationType(notificationvalueobjects.TypeInfo)
	status := notificationvalueobjects.MustNotificationStatus(notificationvalueobjects.StatusRead)
	readAt := time.Now()
	metadata := map[string]interface{}{"key": "value"}
	createdAt := time.Now().Add(-1 * time.Hour)
	updatedAt := time.Now()

	notification := NotificationFromPersistence(
		notificationvalueobjects.GenerateNotificationID(),
		userID,
		title,
		message,
		notifType,
		status,
		&readAt,
		metadata,
		createdAt,
		updatedAt,
	)

	if notification == nil {
		t.Error("NotificationFromPersistence() should not return nil")
	}

	if !notification.Status().IsRead() {
		t.Error("NotificationFromPersistence() should preserve status")
	}

	if notification.ReadAt() == nil || !notification.ReadAt().Equal(readAt) {
		t.Error("NotificationFromPersistence() should preserve ReadAt")
	}

	if len(notification.GetEvents()) != 0 {
		t.Error("NotificationFromPersistence() should not have domain events")
	}
}
