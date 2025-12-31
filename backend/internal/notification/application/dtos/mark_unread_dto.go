package dtos

// MarkNotificationUnreadInput represents the input data for marking a notification as unread.
type MarkNotificationUnreadInput struct {
	NotificationID string `json:"notification_id" validate:"required,uuid"`
	UserID         string `json:"user_id" validate:"required,uuid"`
}

// MarkNotificationUnreadOutput represents the output data after marking a notification as unread.
type MarkNotificationUnreadOutput struct {
	NotificationID string `json:"notification_id"`
	Status         string `json:"status"`
}
