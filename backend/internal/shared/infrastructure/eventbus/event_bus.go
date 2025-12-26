package eventbus

import (
	"fmt"
	"sync"
	"time"

	"gestao-financeira/backend/internal/shared/domain/events"

	"github.com/rs/zerolog/log"
)

// RetryConfig configures retry behavior for event handlers.
type RetryConfig struct {
	MaxRetries        int           // Maximum number of retry attempts
	InitialDelay      time.Duration // Initial delay before first retry
	MaxDelay          time.Duration // Maximum delay between retries
	BackoffMultiplier float64       // Multiplier for exponential backoff
}

// DefaultRetryConfig returns a default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:        3,
		InitialDelay:      100 * time.Millisecond,
		MaxDelay:          5 * time.Second,
		BackoffMultiplier: 2.0,
	}
}

// EventHandler represents a function that handles a domain event.
type EventHandler func(event events.DomainEvent) error

// HandlerWrapper wraps an event handler with retry and error handling logic.
type HandlerWrapper struct {
	handler     EventHandler
	retryConfig RetryConfig
	handlerName string
}

// NewHandlerWrapper creates a new handler wrapper with retry configuration.
func NewHandlerWrapper(handler EventHandler, config RetryConfig, name string) *HandlerWrapper {
	return &HandlerWrapper{
		handler:     handler,
		retryConfig: config,
		handlerName: name,
	}
}

// Execute executes the handler with retry logic.
func (hw *HandlerWrapper) Execute(event events.DomainEvent) error {
	var lastErr error
	delay := hw.retryConfig.InitialDelay

	for attempt := 0; attempt <= hw.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			time.Sleep(delay)

			// Calculate next delay with exponential backoff
			delay = time.Duration(float64(delay) * hw.retryConfig.BackoffMultiplier)
			if delay > hw.retryConfig.MaxDelay {
				delay = hw.retryConfig.MaxDelay
			}

			log.Warn().
				Str("event_type", event.EventType()).
				Str("handler", hw.handlerName).
				Int("attempt", attempt).
				Err(lastErr).
				Msg("Retrying event handler")
		}

		err := hw.handler(event)
		if err == nil {
			if attempt > 0 {
				log.Info().
					Str("event_type", event.EventType()).
					Str("handler", hw.handlerName).
					Int("attempt", attempt).
					Msg("Event handler succeeded after retry")
			}
			return nil
		}

		lastErr = err
	}

	// All retries failed
	log.Error().
		Str("event_type", event.EventType()).
		Str("handler", hw.handlerName).
		Int("max_retries", hw.retryConfig.MaxRetries).
		Err(lastErr).
		Msg("Event handler failed after all retries")

	return fmt.Errorf("handler %s failed after %d retries: %w", hw.handlerName, hw.retryConfig.MaxRetries, lastErr)
}

// EventBus is an in-memory event bus for publishing and subscribing to domain events.
// It supports retry logic and error handling for event handlers.
type EventBus struct {
	handlers      map[string][]*HandlerWrapper
	mu            sync.RWMutex
	defaultRetry  RetryConfig
	errorCallback func(eventType string, handlerName string, err error)
}

// NewEventBus creates a new EventBus instance with default retry configuration.
func NewEventBus() *EventBus {
	return &EventBus{
		handlers:     make(map[string][]*HandlerWrapper),
		defaultRetry: DefaultRetryConfig(),
	}
}

// NewEventBusWithConfig creates a new EventBus instance with custom retry configuration.
func NewEventBusWithConfig(retryConfig RetryConfig) *EventBus {
	return &EventBus{
		handlers:     make(map[string][]*HandlerWrapper),
		defaultRetry: retryConfig,
	}
}

// SetErrorCallback sets a callback function that will be called when a handler fails.
func (eb *EventBus) SetErrorCallback(callback func(eventType string, handlerName string, err error)) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.errorCallback = callback
}

// Subscribe registers an event handler for a specific event type.
// The handler will be called whenever an event of that type is published.
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.SubscribeWithRetry(eventType, handler, eb.defaultRetry, "")
}

// SubscribeWithRetry registers an event handler with custom retry configuration.
func (eb *EventBus) SubscribeWithRetry(eventType string, handler EventHandler, retryConfig RetryConfig, handlerName string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.handlers[eventType] == nil {
		eb.handlers[eventType] = make([]*HandlerWrapper, 0)
	}

	// Use handler name or generate one
	if handlerName == "" {
		handlerName = fmt.Sprintf("handler_%d", len(eb.handlers[eventType]))
	}

	wrapper := NewHandlerWrapper(handler, retryConfig, handlerName)
	eb.handlers[eventType] = append(eb.handlers[eventType], wrapper)
}

// Unsubscribe removes all event handlers for a specific event type.
func (eb *EventBus) Unsubscribe(eventType string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	delete(eb.handlers, eventType)
}

// Publish publishes a domain event to all registered handlers.
// Handlers are executed synchronously in the order they were registered.
// If a handler returns an error, retry logic is applied.
// Errors from handlers are collected and logged, but other handlers continue to execute.
func (eb *EventBus) Publish(event events.DomainEvent) error {
	eb.mu.RLock()
	handlers := eb.handlers[event.EventType()]
	eb.mu.RUnlock()

	if len(handlers) == 0 {
		// No handlers registered for this event type - this is not an error
		return nil
	}

	var errors []error
	for _, wrapper := range handlers {
		if err := wrapper.Execute(event); err != nil {
			errorMsg := fmt.Errorf("handler %s error for event %s: %w", wrapper.handlerName, event.EventType(), err)
			errors = append(errors, errorMsg)

			// Call error callback if set
			if eb.errorCallback != nil {
				eb.errorCallback(event.EventType(), wrapper.handlerName, err)
			}
		}
	}

	if len(errors) > 0 {
		errorSummary := fmt.Errorf("errors occurred while handling event %s: %v", event.EventType(), errors)
		log.Error().
			Str("event_type", event.EventType()).
			Int("error_count", len(errors)).
			Msg("Errors occurred while handling event")
		return errorSummary
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

		for _, wrapper := range handlers {
			wg.Add(1)
			go func(hw *HandlerWrapper) {
				defer wg.Done()
				if err := hw.Execute(event); err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("handler %s error for event %s: %w", hw.handlerName, event.EventType(), err))
					mu.Unlock()
				}
			}(wrapper)
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

	eb.handlers = make(map[string][]*HandlerWrapper)
}
