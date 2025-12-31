package dtos

// CreateNotificationInput represents the input data for notification creation.
type CreateNotificationInput struct {
	UserID   string                 `json:"user_id" validate:"required,uuid"`
	Title    string                 `json:"title" validate:"required,min=1,max=200,no_sql_injection,no_xss,utf8"`
	Message  string                 `json:"message" validate:"required,min=1,max=1000,no_sql_injection,no_xss,utf8"`
	Type     string                 `json:"type" validate:"required,oneof=INFO WARNING SUCCESS ERROR"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CreateNotificationOutput represents the output data after notification creation.
type CreateNotificationOutput struct {
	NotificationID string                 `json:"notification_id"`
	UserID         string                 `json:"user_id"`
	Title          string                 `json:"title"`
	Message        string                 `json:"message"`
	Type           string                 `json:"type"`
	Status         string                 `json:"status"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt      string                 `json:"created_at"`
}
