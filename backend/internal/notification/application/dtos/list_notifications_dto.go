package dtos

// ListNotificationsInput represents the input data for listing notifications.
type ListNotificationsInput struct {
	UserID   string `json:"user_id" validate:"required,uuid"`
	Status   string `json:"status,omitempty" validate:"omitempty,oneof=UNREAD READ ARCHIVED"`
	Type     string `json:"type,omitempty" validate:"omitempty,oneof=INFO WARNING SUCCESS ERROR"`
	Page     int    `json:"page,omitempty" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size,omitempty" validate:"omitempty,min=1,max=100"`
}

// ListNotificationsOutput represents the output data for listing notifications.
type ListNotificationsOutput struct {
	Notifications []NotificationItem `json:"notifications"`
	Total         int64              `json:"total"`
	Page          int                `json:"page"`
	PageSize      int                `json:"page_size"`
	TotalPages    int                `json:"total_pages"`
}

// NotificationItem represents a single notification in the list.
type NotificationItem struct {
	NotificationID string                 `json:"notification_id"`
	Title          string                 `json:"title"`
	Message        string                 `json:"message"`
	Type           string                 `json:"type"`
	Status         string                 `json:"status"`
	ReadAt         *string                `json:"read_at,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt      string                 `json:"created_at"`
}
