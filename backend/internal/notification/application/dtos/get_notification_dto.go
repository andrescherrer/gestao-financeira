package dtos

// GetNotificationInput represents the input data for getting a notification.
type GetNotificationInput struct {
	NotificationID string `json:"notification_id" validate:"required,uuid"`
	UserID         string `json:"user_id" validate:"required,uuid"`
}

// GetNotificationOutput represents the output data for getting a notification.
type GetNotificationOutput struct {
	NotificationID string                 `json:"notification_id"`
	UserID         string                 `json:"user_id"`
	Title          string                 `json:"title"`
	Message        string                 `json:"message"`
	Type           string                 `json:"type"`
	Status         string                 `json:"status"`
	ReadAt         *string                `json:"read_at,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt      string                 `json:"created_at"`
	UpdatedAt      string                 `json:"updated_at"`
}
