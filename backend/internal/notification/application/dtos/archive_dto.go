package dtos

// ArchiveNotificationInput represents the input data for archiving a notification.
type ArchiveNotificationInput struct {
	NotificationID string `json:"notification_id" validate:"required,uuid"`
	UserID         string `json:"user_id" validate:"required,uuid"`
}

// ArchiveNotificationOutput represents the output data after archiving a notification.
type ArchiveNotificationOutput struct {
	NotificationID string `json:"notification_id"`
	Status         string `json:"status"`
}
