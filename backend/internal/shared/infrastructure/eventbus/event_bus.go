package eventbus

import (
	"fmt"
	"sync"

	"gestao-financeira/backend/internal/shared/domain/events"
)

// EventHandler represents a function that handles a domain event.
type EventHandler func(event events.DomainEvent) error

// EventBus is a simple in-memory event bus for publishing and subscribing to domain events.
type EventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewEventBus creates a new EventBus instance.
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe registers an event handler for a specific event type.
// The handler will be called whenever an event of that type is published.
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.handlers[eventType] == nil {
		eb.handlers[eventType] = make([]EventHandler, 0)
	}

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// Unsubscribe removes an event handler for a specific event type.
// Note: This removes all handlers of that type. For more granular control,
// consider using a handler ID or wrapper.
func (eb *EventBus) Unsubscribe(eventType string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	delete(eb.handlers, eventType)
}

// Publish publishes a domain event to all registered handlers.
// Handlers are executed synchronously in the order they were registered.
// If a handler returns an error, the error is logged but other handlers continue to execute.
func (eb *EventBus) Publish(event events.DomainEvent) error {
	eb.mu.RLock()
	handlers := eb.handlers[event.EventType()]
	eb.mu.RUnlock()

	if len(handlers) == 0 {
		// No handlers registered for this event type - this is not an error
		return nil
	}

	var errors []error
	for _, handler := range handlers {
		if err := handler(event); err != nil {
			errors = append(errors, fmt.Errorf("handler error for event %s: %w", event.EventType(), err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors occurred while handling event %s: %v", event.EventType(), errors)
	}

	return nil
}

// PublishAsync publishes a domain event to all registered handlers asynchronously.
// Handlers are executed in separate goroutines.
// Errors from handlers are collected and returned via a channel.
func (eb *EventBus) PublishAsync(event events.DomainEvent) <-chan error {
	errChan := make(chan error, 1)

	eb.mu.RLock()
	handlers := eb.handlers[event.EventType()]
	eb.mu.RUnlock()

	if len(handlers) == 0 {
		close(errChan)
		return errChan
	}

	go func() {
		var wg sync.WaitGroup
		var errors []error
		var mu sync.Mutex

		for _, handler := range handlers {
			wg.Add(1)
			go func(h EventHandler) {
				defer wg.Done()
				if err := h(event); err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("handler error for event %s: %w", event.EventType(), err))
					mu.Unlock()
				}
			}(handler)
		}

		wg.Wait()

		if len(errors) > 0 {
			errChan <- fmt.Errorf("errors occurred while handling event %s: %v", event.EventType(), errors)
		} else {
			close(errChan)
		}
	}()

	return errChan
}

// HasSubscribers checks if there are any handlers registered for a specific event type.
func (eb *EventBus) HasSubscribers(eventType string) bool {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	handlers, exists := eb.handlers[eventType]
	return exists && len(handlers) > 0
}

// GetSubscriberCount returns the number of handlers registered for a specific event type.
func (eb *EventBus) GetSubscriberCount(eventType string) int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	return len(eb.handlers[eventType])
}

// Clear removes all registered handlers.
func (eb *EventBus) Clear() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers = make(map[string][]EventHandler)
}
