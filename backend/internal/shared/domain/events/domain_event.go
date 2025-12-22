package events

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent represents a domain event that occurred in the system.
// Domain events are used to communicate state changes within and across bounded contexts.
type DomainEvent interface {
	// EventID returns a unique identifier for this event instance.
	EventID() string

	// EventType returns the type/name of the event (e.g., "UserRegistered", "AccountCreated").
	EventType() string

	// AggregateID returns the identifier of the aggregate that generated this event.
	AggregateID() string

	// AggregateType returns the type of the aggregate (e.g., "User", "Account").
	AggregateType() string

	// OccurredAt returns the timestamp when the event occurred.
	OccurredAt() time.Time

	// Version returns the version of the event schema (useful for event evolution).
	Version() int
}

// BaseDomainEvent provides a base implementation for domain events.
// It can be embedded in specific event types to reduce boilerplate.
type BaseDomainEvent struct {
	eventID       string
	eventType     string
	aggregateID   string
	aggregateType string
	occurredAt    time.Time
	version       int
}

// NewBaseDomainEvent creates a new base domain event.
func NewBaseDomainEvent(eventType, aggregateID, aggregateType string) BaseDomainEvent {
	return BaseDomainEvent{
		eventID:       generateEventID(),
		eventType:     eventType,
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
		occurredAt:    time.Now(),
		version:       1,
	}
}

// EventID returns the unique identifier for this event instance.
func (e BaseDomainEvent) EventID() string {
	return e.eventID
}

// EventType returns the type/name of the event.
func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

// AggregateID returns the identifier of the aggregate that generated this event.
func (e BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}

// AggregateType returns the type of the aggregate.
func (e BaseDomainEvent) AggregateType() string {
	return e.aggregateType
}

// OccurredAt returns the timestamp when the event occurred.
func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// Version returns the version of the event schema.
func (e BaseDomainEvent) Version() int {
	return e.version
}

// generateEventID generates a unique identifier for an event using UUID.
func generateEventID() string {
	return uuid.New().String()
}
