package dtos

// DeleteNotificationInput represents the input data for deleting a notification.
type DeleteNotificationInput struct {
	NotificationID string `json:"notification_id" validate:"required,uuid"`
	UserID         string `json:"user_id" validate:"required,uuid"`
}

// DeleteNotificationOutput represents the output data after deleting a notification.
type DeleteNotificationOutput struct {
	NotificationID string `json:"notification_id"`
	Deleted        bool   `json:"deleted"`
}
