package handlers

import (
	"gestao-financeira/backend/internal/shared/domain/events"

	"github.com/rs/zerolog/log"
)

// EventLoggerHandler logs all domain events for observability.
type EventLoggerHandler struct{}

// NewEventLoggerHandler creates a new EventLoggerHandler instance.
func NewEventLoggerHandler() *EventLoggerHandler {
	return &EventLoggerHandler{}
}

// Handle logs a domain event.
func (h *EventLoggerHandler) Handle(event events.DomainEvent) error {
	log.Info().
		Str("event_id", event.EventID()).
		Str("event_type", event.EventType()).
		Str("aggregate_id", event.AggregateID()).
		Str("aggregate_type", event.AggregateType()).
		Time("occurred_at", event.OccurredAt()).
		Int("version", event.Version()).
		Msg("Domain event published")

	return nil
}
