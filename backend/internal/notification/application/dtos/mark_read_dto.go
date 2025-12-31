package dtos

// MarkNotificationReadInput represents the input data for marking a notification as read.
type MarkNotificationReadInput struct {
	NotificationID string `json:"notification_id" validate:"required,uuid"`
	UserID         string `json:"user_id" validate:"required,uuid"`
}

// MarkNotificationReadOutput represents the output data after marking a notification as read.
type MarkNotificationReadOutput struct {
	NotificationID string `json:"notification_id"`
	Status           string `json:"status"`
	ReadAt           string `json:"read_at"`
}

