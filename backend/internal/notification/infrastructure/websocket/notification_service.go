package websocket

import (
	"encoding/json"

	"gestao-financeira/backend/internal/notification/application/dtos"

	"github.com/rs/zerolog/log"
)

// NotificationService handles sending notifications via WebSocket.
type NotificationService struct {
	hub *Hub
}

// NewNotificationService creates a new NotificationService instance.
func NewNotificationService(hub *Hub) *NotificationService {
	return &NotificationService{
		hub: hub,
	}
}

// SendNotification sends a notification to a user via WebSocket.
func (s *NotificationService) SendNotification(userID string, notification *dtos.CreateNotificationOutput) error {
	message := map[string]interface{}{
		"type": "notification",
		"data": notification,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("Failed to marshal notification message")
		return err
	}

	s.hub.SendToUser(userID, messageBytes)
	return nil
}

// SendNotificationUpdate sends a notification update (e.g., marked as read) via WebSocket.
func (s *NotificationService) SendNotificationUpdate(userID string, notificationID string, updateType string, data interface{}) error {
	message := map[string]interface{}{
		"type":            "notification_update",
		"notification_id": notificationID,
		"update_type":     updateType,
		"data":            data,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("Failed to marshal notification update message")
		return err
	}

	s.hub.SendToUser(userID, messageBytes)
	return nil
}
