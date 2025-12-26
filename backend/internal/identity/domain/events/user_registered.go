package events

import (
	"gestao-financeira/backend/internal/shared/domain/events"
)

// UserRegistered represents a domain event that occurs when a user is registered.
type UserRegistered struct {
	events.BaseDomainEvent
	email string
	name  string
}

// NewUserRegistered creates a new UserRegistered event.
func NewUserRegistered(userID string, email string, name string) *UserRegistered {
	baseEvent := events.NewBaseDomainEvent(
		"UserRegistered",
		userID,
		"User",
	)

	return &UserRegistered{
		BaseDomainEvent: baseEvent,
		email:           email,
		name:            name,
	}
}

// Email returns the user's email.
func (e *UserRegistered) Email() string {
	return e.email
}

// Name returns the user's name.
func (e *UserRegistered) Name() string {
	return e.name
}
